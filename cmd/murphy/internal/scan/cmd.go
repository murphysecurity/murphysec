package scan

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/maven"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var jsonOutput bool
var isDeep bool
var noBuild bool
var projectNameCli string
var mavenSettingsPath string
var onlyTaskId bool
var privateSourceId string
var privateSourceName string
var projectTagNames []string
var concurrentNumber int
var sbomOutputConfig string
var sbomOutputType common.SBOMFormatFlag

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan <DIR>"
	c.Short = "Detects open source vulnerabilities by scanning various file types within the project"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	c.Flags().BoolVar(&jsonOutput, "json", false, "output in json format")
	c.Flags().BoolVar(&isDeep, "deep", false, "enable enhanced deep insight, code features identification, vulnerability accessibility analysis")
	c.Flags().BoolVar(&noBuild, "no-build", false, "skip project building")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringVar(&privateSourceId, "maven-setting-id", "", "specify the id of the Maven settings.xml file used during the scan")
	c.Flags().StringVar(&privateSourceName, "maven-setting-name", "", "specify the name of the Maven settings.xml file used during the scan")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().IntVarP(&concurrentNumber, "max-concurrent-uploads", "j", 1, "Set the maximum number of parallel uploads.")
	return &c
}

func DfCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "dfscan <DIR>"
	c.Args = cobra.ExactArgs(1)
	c.Run = dfScanRun
	c.Short = "Detects open source vulnerabilities by scanning package management files"
	c.Flags().BoolVar(&jsonOutput, "json", false, "output in json format")
	c.Flags().BoolVar(&isDeep, "deep", false, "enable enhanced deep insight, code features identification, vulnerability accessibility analysis")
	c.Flags().BoolVar(&noBuild, "no-build", false, "skip project building")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().StringVar(&mavenSettingsPath, "maven-settings", "", "specify the path of maven settings")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().StringVar(&sbomOutputConfig, "sbom-output", "-", "Specify the SBOM output file path, use \"-\" to output to stdout")
	c.Flags().Var(&sbomOutputType, "sbom-format", "(Required) Specify the SBOM format, currently supports: murphysec1.1+json")
	return &c
}

func EnvCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "envscan"
	c.Run = envScanRun
	c.Short = "Detects open source vulnerabilities environment"
	c.Flags().BoolVar(&jsonOutput, "json", false, "output in json format")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().StringVar(&sbomOutputConfig, "sbom-output", "-", "Specify the SBOM output file path, use \"-\" to output to stdout")
	c.Flags().Var(&sbomOutputType, "sbom-format", "(Required) Specify the SBOM format, currently supports: murphysec1.1+json")
	return &c
}

func commonInitNoAPI(ctx context.Context) (context.Context, error) {
	// init logging
	ctx, e := common.InitLogger(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		reportIdeError(ctx, model.IDEStatusLogFileCreationError, e)
		exitcode.Set(1)
		return nil, e
	}
	return ctx, nil
}

func commonInit(ctx context.Context) (context.Context, error) {
	// init logging
	ctx, e := common.InitLogger(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		reportIdeError(ctx, model.IDEStatusLogFileCreationError, e)
		exitcode.Set(1)
		return nil, e
	}
	var logger = logctx.Use(ctx).Sugar()
	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		logger.Error(e)
		reportIdeError(ctx, model.IDEStatusAPIFail, e)
		exitcode.Set(1)
		return nil, e
	}
	return ctx, nil
}

func commonScanPreCheck(ctx context.Context, scanDir string) (string, error) {
	// get absolute path and check if a directory
	scanDir, e := filepath.Abs(scanDir)
	if e != nil {
		cv.DisplayScanInvalidPath(ctx, e)
		return "", e
	}
	if !utils.IsDir(scanDir) {
		cv.DisplayScanInvalidPathMustDir(ctx, nil)
		exitcode.Set(1)
		return "", fmt.Errorf("dir invalid")
	}
	return scanDir, nil
}

func scanRun(cmd *cobra.Command, args []string) {
	var ctx = context.TODO()
	if jsonOutput {
		ctx = ui.With(ctx, ui.IDEA)
	} else {
		ctx = ui.With(ctx, ui.CLI)
	}
	scanDir := args[0]
	scanDir, e := commonScanPreCheck(ctx, scanDir)
	if e != nil {
		return
	}
	ctx, e = commonInit(ctx)
	if e != nil {
		return
	}
	logger := logctx.Use(ctx).Sugar()
	r, e := scan(ctx, scanDir, model.AccessTypeCli, model.ScanModeStandard)
	if errors.Is(e, inspector.ErrNoWait) {
		return
	}
	if e != nil {
		logger.Error(e)
		autoReportIde(ctx, e)
		exitcode.Set(1)
		return
	}
	if onlyTaskId {
		return
	}
	if jsonOutput {
		fmt.Println(string(must.A(json.MarshalIndent(model.GetIDEAOutput(r), "", "  "))))
	}
}

func envScanRun(cmd *cobra.Command, args []string) {
	var ctx = context.TODO()
	if sbomOutputType.Valid {
		ctx = ui.With(ctx, ui.None)
	} else if jsonOutput {
		ctx = ui.With(ctx, ui.IDEA)
	} else if onlyTaskId {
		ctx = ui.With(ctx, ui.None)
	} else {
		ctx = ui.With(ctx, ui.CLI)
	}
	var e error
	if sbomOutputType.Valid {
		ctx, e = commonInitNoAPI(ctx)
	} else {
		ctx, e = commonInit(ctx)
	}
	if e != nil {
		return
	}
	logger := logctx.Use(ctx).Sugar()
	var r *model.ScanTask
	if sbomOutputType.Valid {
		r, e = envScanSbomOnly(ctx)
		if e != nil {
			exitcode.Set(1)
		}
		doSBOMOnlyPrint(ctx, r)
		return
	} else {
		r, e = envScan(ctx)
	}
	if errors.Is(e, inspector.ErrNoWait) {
		return
	}
	if e != nil {
		logger.Error(e)
		autoReportIde(ctx, e)
		exitcode.Set(1)
		return
	}
	if onlyTaskId {
		return
	}
	if jsonOutput {
		fmt.Println(string(must.A(json.MarshalIndent(model.GetIDEAOutput(r), "", "  "))))
	}
}

func dfScanRun(cmd *cobra.Command, args []string) {
	var ctx = context.TODO()
	ctx = scanerr.WithCtx(ctx)
	if sbomOutputType.Valid {
		ctx = ui.With(ctx, ui.None)
	} else if jsonOutput {
		ctx = ui.With(ctx, ui.IDEA)
	} else if onlyTaskId {
		ctx = ui.With(ctx, ui.None)
	} else {
		ctx = ui.With(ctx, ui.CLI)
	}

	if mavenSettingsPath != "" {
		//nolint:all
		ctx = context.WithValue(ctx, maven.M2SettingsFilePathCtxKey, mavenSettingsPath)
	}
	scanDir := args[0]
	scanDir, e := commonScanPreCheck(ctx, scanDir)
	if e != nil {
		return
	}
	if sbomOutputType.Valid {
		ctx, e = commonInitNoAPI(ctx)
	} else {
		ctx, e = commonInit(ctx)
	}
	if e != nil {
		return
	}
	if sbomOutputType.Valid {
		scanSbomOnly(ctx, scanDir)
		return
	}
	logger := logctx.Use(ctx).Sugar()
	r, e := scan(ctx, scanDir, model.AccessTypeCli, model.ScanModeSource)
	if errors.Is(e, inspector.ErrNoWait) {
		return
	}
	if e != nil {
		logger.Error(e)
		autoReportIde(ctx, e)
		exitcode.Set(1)
		return
	}
	if onlyTaskId {
		return
	}
	if jsonOutput {
		fmt.Println(string(must.A(json.MarshalIndent(model.GetIDEAOutput(r), "", "  "))))
	}
}

func IdeaScan() *cobra.Command {
	var c cobra.Command
	c.Use = "ideascan <DIR>"
	c.Args = cobra.ExactArgs(1)
	c.Run = ideascanRun
	c.Hidden = true
	c.Flags().String("ide", "", "unused")
	must.M(c.Flags().MarkHidden("ide"))
	return &c
}

func ideascanRun(cmd *cobra.Command, args []string) {
	ctx := ui.With(context.TODO(), ui.IDEA)
	accessType := model.AccessTypeIdea
	scanDir := args[0]
	// get absolute path and check if a directory
	scanDir, e := filepath.Abs(scanDir)
	if e != nil {
		reportIdeError(ctx, model.IDEStatusScanDirInvalid, e)
		exitcode.Set(1)
		return
	}
	if !utils.IsDir(scanDir) {
		reportIdeError(ctx, model.IDEStatusScanDirInvalid, fmt.Errorf("not a dir"))
		exitcode.Set(1)
		return
	}

	// init logging
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		reportIdeError(ctx, model.IDEStatusLogFileCreationError, e)
		exitcode.Set(1)
		return
	}
	var logger = logctx.Use(ctx).Sugar()

	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		reportIdeError(ctx, model.IDEStatusAPIFail, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}

	task, e := scan(ctx, scanDir, accessType, model.ScanModeSource)
	if errors.Is(e, inspector.ErrNoWait) {
		return
	}
	if e != nil {
		autoReportIde(ctx, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}
	fmt.Println(string(must.A(json.MarshalIndent(model.GetIDEAOutput(task), "", "  "))))
}

type ideErrorResp struct {
	ErrCode model.IDEStatus `json:"err_code"`
	ErrMsg  string          `json:"err_msg"`
}

func autoReportIde(ctx context.Context, e error) {
	if errors.Is(e, api.ErrTokenInvalid) {
		reportIdeError(ctx, model.IDEStatusTokenInvalid, e)
		return
	}
	if errors.Is(e, api.ErrServerFail) {
		reportIdeError(ctx, model.IDEStatusServerFail, e)
		return
	}
	if errors.Is(e, api.ErrGeneralError) {
		reportIdeError(ctx, model.IDEStatusGeneralAPIError, e)
		return
	}
	if errors.Is(e, api.ErrRequest) {
		reportIdeError(ctx, model.IDEStatusAPIFail, e)
		return
	}
	reportIdeError(ctx, model.IDEStatusUnknownError, e)
}

func reportIdeError(ctx context.Context, status model.IDEStatus, e error) {
	if ui.Use(ctx) != ui.IDEA {
		return
	}
	resp := ideErrorResp{
		ErrCode: status,
		ErrMsg:  status.String(),
	}
	if e != nil {
		resp.ErrMsg = e.Error()
	}
	fmt.Println(string(must.A(json.MarshalIndent(resp, "", "  "))))
}

func doSBOMOnlyPrint(ctx context.Context, task *model.ScanTask) {
	var logger = logctx.Use(ctx)
	_ = logger.Sync()
	if sbomOutputConfig == "" {
		panic("sbomOutputConfig == \"\"")
	}
	var writer io.Writer
	if sbomOutputConfig == "-" {
		writer = os.Stdout
	} else {
		f, e := os.OpenFile(sbomOutputConfig, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if e != nil {
			panic(e)
		}
		writer = f
		defer func() {
			var e = f.Close()
			if e != nil {
				panic(e)
			}
		}()
	}
	var bufioWriter = bufio.NewWriter(writer)
	var enc = json.NewEncoder(bufioWriter)
	must.M(bufioWriter.Flush())
	enc.SetIndent("", "    ")
	if task.Modules == nil {
		task.Modules = make([]model.Module, 0)
	}
	must.M(enc.Encode(map[string]any{"modules": task.Modules}))
}

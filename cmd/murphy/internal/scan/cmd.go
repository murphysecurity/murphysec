package scan

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/gradle"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/toolver"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var jsonOutput bool
var isDeep bool
var noBuild bool
var projectNameCli string
var projectsNameCli string
var onlyTaskId bool
var privateSourceId string
var privateSourceName string
var projectTagNames []string
var concurrentNumber int
var sbomOutputConfig string
var sbomOutputType common.SBOMFormatFlag
var webhookAddr string
var webhookMode common.WebhookModeFlag
var extraData string
var scanCodeHash bool
var gradleProjectFilter gradle.ProjectFilter
var branch string
var mavenModuleName []string
var binaryOnly bool
var scanProcess bool
var distribution common.DistributionFlag
var disableWindowsPatchScan bool
var windowsPatchScanTimeout int
var webhookToken []string

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan <DIR>"
	c.Short = "Detects open source vulnerabilities by scanning various file types within the project"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	c.Flags().BoolVar(&jsonOutput, "json", false, "output in json format")
	c.Flags().BoolVar(&isDeep, "deep", false, "enable enhanced deep insight, code features identification, vulnerability accessibility analysis")
	c.Flags().BoolVar(&noBuild, "no-build", false, "skip project building")
	c.Flags().StringVar(&branch, "branch", "", "")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().StringVar(&projectsNameCli, "projects-name", "", "specify projects name(group)")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringVar(&privateSourceId, "maven-setting-id", "", "specify the id of the Maven settings.xml file used during the scan")
	c.Flags().StringVar(&privateSourceName, "maven-setting-name", "", "specify the name of the Maven settings.xml file used during the scan")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().IntVarP(&concurrentNumber, "max-concurrent-uploads", "j", 1, "Set the maximum number of parallel uploads.")
	c.Flags().StringVar(&webhookAddr, "webhook-addr", "", "specify the webhook address")
	c.Flags().Var(&webhookMode, "webhook-mode", "specify the webhook mode, currently supports: simple, full")
	c.Flags().Var(&distribution, "distribution", "specify the distribution, currently supports: external, internal, saas, open_source")
	c.Flags().StringVar(&extraData, "extra-data", "", "specify the extra data")
	c.Flags().BoolVar(&scanCodeHash, "scan-snippets", false, "Enable scanning of code snippets to detect SBOM and  vulnerabilities. Disabled by default")
	c.Flags().BoolVar(&binaryOnly, "binary-only", false, "only scan binary files, skip source code scanning")
	c.Flags().StringArrayVar(&toolver.Default.Maven.AdditionalPrependArgs, "maven-prepend-arg", []string{}, "Prepend an argument to the Maven command. Can be specified multiple times.")
	c.Flags().StringArrayVar(&toolver.Default.Maven.AdditionalArgs, "maven-arg", []string{}, "Append an argument to the Maven command. Can be specified multiple times.")
	c.Flags().StringVar(&toolver.Default.Maven.JdkVersion, "maven-jdk", "", "specify JDK version for Maven build")
	c.Flags().StringVar(&toolver.Default.Maven.MavenVersion, "maven-version", "", "specify Maven version for Maven build")
	// 多个入参： murphy scan ./project --webhook-token Authorization=Bearer123 --webhook-token X-Custom-Header=value
	c.Flags().StringArrayVar(&webhookToken, "webhook-token", make([]string, 0), "specify the webhook token in key=value format. Can be specified multiple times.")
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
	c.Flags().StringVar(&branch, "branch", "", "")
	c.Flags().StringArrayVar(&mavenModuleName, "maven-module-name", make([]string, 0), "retains module")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().StringVar(&projectsNameCli, "projects-name", "", "specify projects name(group)")
	c.Flags().StringVar(&toolver.Default.Maven.MavenSettingPath, "maven-settings", "", "specify the path of maven settings")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().StringVar(&sbomOutputConfig, "sbom-output", "-", "Specify the SBOM output file path, use \"-\" to output to stdout")
	c.Flags().Var(&sbomOutputType, "sbom-format", "(Required) specify the SBOM format, currently supports: msdx1.1+json")
	c.Flags().StringVar(&webhookAddr, "webhook-addr", "", "specify the webhook address")
	c.Flags().Var(&webhookMode, "webhook-mode", "specify the webhook mode, currently supports: simple, full(default)")
	c.Flags().Var(&distribution, "distribution", "specify the distribution, currently supports: external, internal, saas, open_source")
	c.Flags().StringVar(&extraData, "extra-data", "", "specify the extra data")
	c.Flags().BoolVar(&scanCodeHash, "scan-snippets", false, "Enable scanning of code snippets to detect SBOM and  vulnerabilities. Disabled by default")
	c.Flags().StringArrayVar(&gradleProjectFilter.ProjectNames, "gradle-project-name", make([]string, 0), "specify the name of the Gradle project")
	c.Flags().StringArrayVar(&toolver.Default.Maven.AdditionalPrependArgs, "maven-prepend-arg", []string{}, "Prepend an argument to the Maven command. Can be specified multiple times.")
	c.Flags().StringArrayVar(&toolver.Default.Maven.AdditionalArgs, "maven-arg", []string{}, "Append an argument to the Maven command. Can be specified multiple times.")
	c.Flags().StringArrayVar(&webhookToken, "webhook-token", make([]string, 0), "specify the webhook token in key=value format. Can be specified multiple times.")
	return &c
}

func EnvCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "envscan"
	c.Run = envScanRun
	c.Short = "Detects open source vulnerabilities environment"
	c.Flags().BoolVar(&jsonOutput, "json", false, "output in json format")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	c.Flags().StringVar(&projectsNameCli, "projects-name", "", "specify projects name(group)")
	c.Flags().BoolVar(&onlyTaskId, "only-task-id", false, "print task id after task created, the scan result will not be printed")
	c.Flags().StringArrayVar(&projectTagNames, "project-tag", make([]string, 0), "specify the tag of the project")
	c.Flags().StringVar(&sbomOutputConfig, "sbom-output", "-", "Specify the SBOM output file path, use \"-\" to output to stdout")
	c.Flags().Var(&sbomOutputType, "sbom-format", "(Required) specify the SBOM format, currently supports: msdx1.1+json")
	c.Flags().StringVar(&webhookAddr, "webhook-addr", "", "specify the webhook address")
	c.Flags().Var(&webhookMode, "webhook-mode", "specify the webhook mode, currently supports: simple, full(default)")
	c.Flags().StringVar(&extraData, "extra-data", "", "specify the extra data")
	c.Flags().BoolVar(&scanProcess, "scan-process", false, "Enable scanning of process to detect SBOM. Disabled by default")
	c.Flags().BoolVar(&disableWindowsPatchScan, "disable-windows-patch-scan", false, "Disable scanning of Windows patches. Enabled by default")
	c.Flags().IntVar(&windowsPatchScanTimeout, "windows-patch-scan-timeout", 60, "Timeout for Windows patch scan in seconds. Default is 60 seconds")
	c.Flags().StringArrayVar(&webhookToken, "webhook-token", make([]string, 0), "specify the webhook token in key=value format. Can be specified multiple times.")
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
	if windowsPatchScanTimeout < 1 {
		windowsPatchScanTimeout = 0
	}
	var windowsPatchScanTimeoutDuration = time.Duration(windowsPatchScanTimeout) * time.Second
	if disableWindowsPatchScan {
		windowsPatchScanTimeoutDuration = 0
	}
	if sbomOutputType.Valid {
		r, e = envScanSbomOnly(ctx, windowsPatchScanTimeoutDuration)
		if e != nil {
			exitcode.Set(1)
		}
		doSBOMOnlyPrint(ctx, r)
		return
	} else {
		r, e = envScan(ctx, windowsPatchScanTimeoutDuration)
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
	ctx = context.WithValue(ctx, gradle.ProjectFilterCtxKey, gradleProjectFilter)
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
	for i := range task.Modules {
		for j := range task.Modules[i].Dependencies {
			task.Modules[i].Dependencies[j].Postprocess()
		}
	}
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
	enc.SetIndent("", "    ")
	if task.Modules == nil {
		task.Modules = make([]model.Module, 0)
	}
	must.M(enc.Encode(map[string]any{
		"modules":             task.Modules,
		"scan_warnings_codes": lo.Uniq(fp.Map(func(it scanerr.Param) string { return it.Kind })(scanerr.GetAll(ctx))),
	}))
	must.M(bufioWriter.Flush())
}

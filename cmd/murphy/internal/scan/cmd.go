package scan

import (
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
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

var jsonOutput bool
var isDeep bool
var noBuild bool

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan <DIR>"
	c.Short = "scan project directory"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	c.Flags().BoolVar(&jsonOutput, "json", false, "")
	c.Flags().BoolVar(&isDeep, "deep", false, "")
	c.Flags().BoolVar(&noBuild, "no-build", false, "")
	return &c
}

func DfCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "dfscan <DIR>"
	c.Args = cobra.ExactArgs(1)
	c.Run = dfScanRun
	c.Flags().BoolVar(&jsonOutput, "json", false, "")
	c.Flags().BoolVar(&isDeep, "deep", false, "")
	c.Flags().BoolVar(&noBuild, "no-build", false, "")
	return &c
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
		return "", e
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
	_, e = scan(ctx, scanDir, model.AccessTypeCli, model.ScanModeStandard)
	if e != nil {
		logger.Error(e)
		autoReportIde(ctx, e)
		exitcode.Set(1)
		return
	}
}

func dfScanRun(cmd *cobra.Command, args []string) {
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
	_, e = scan(ctx, scanDir, model.AccessTypeCli, model.ScanModeSource)
	if e != nil {
		logger.Error(e)
		autoReportIde(ctx, e)
		exitcode.Set(1)
		return
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

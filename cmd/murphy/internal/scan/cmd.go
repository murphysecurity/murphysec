package scan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

var cliTaskIdOverride string
var jsonOutput bool

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan <DIR>"
	c.Short = "scan project directory"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	c.Flags().StringVar(&cliTaskIdOverride, "task-id", "", "specify task id, and write it to config")
	c.Flags().BoolVar(&jsonOutput, "json", false, "")
	return &c
}

func scanRun(cmd *cobra.Command, args []string) {
	if jsonOutput {
		ideascanRun(cmd, args)
		return
	}
	var (
		// workaround
		ctx     = ui.With(context.TODO(), ui.CLI{})
		scanDir = args[0]
		e       error
	)
	// get absolute path and check if a directory
	scanDir, e = filepath.Abs(scanDir)
	if e != nil {
		cv.DisplayScanInvalidPath(ctx, e)
	}
	if !utils.IsDir(scanDir) {
		cv.DisplayScanInvalidPathMustDir(ctx, nil)
		exitcode.Set(1)
		return
	}

	// init logging
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		exitcode.Set(1)
		return
	}
	var logger = logctx.Use(ctx).Sugar()

	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}

	if cliTaskIdOverride != "" {
		logger.Infof("CLI task id override: %s", cliTaskIdOverride)
		cf := config.RepoConfig{TaskId: cliTaskIdOverride}
		if e := cf.Validate(); e != nil {
			cv.DisplayBadTaskId(ctx)
			logger.Error(e)
			exitcode.Set(1)
			return
		}
		e = config.WriteRepoConfig(ctx, scanDir, model.AccessTypeCli, config.RepoConfig{TaskId: cliTaskIdOverride})
		if e != nil {
			cv.DisplayInitializeFailed(ctx, e)
			logger.Error(e)
			exitcode.Set(1)
			return
		}
	}
	_, e = scan(ctx, scanDir, model.AccessTypeCli)
	if e != nil {
		logger.Error(e)
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
	c.Flags().StringVar(&cliTaskIdOverride, "task-id", "", "specify task id, and write it to config")
	return &c
}

func ideascanRun(cmd *cobra.Command, args []string) {
	var (
		// workaround
		ctx        = ui.With(context.TODO(), ui.None{})
		scanDir    = args[0]
		e          error
		accessType = model.AccessTypeIdea
	)
	if cmd.Use == "scan" {
		accessType = model.AccessTypeCli
	}
	// get absolute path and check if a directory
	scanDir, e = filepath.Abs(scanDir)
	if e != nil {
		reportIdeError(model.IDEStatusScanDirInvalid, e)
		exitcode.Set(1)
		return
	}
	if !utils.IsDir(scanDir) {
		reportIdeError(model.IDEStatusScanDirInvalid, fmt.Errorf("not a dir"))
		exitcode.Set(1)
		return
	}

	// init logging
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		reportIdeError(model.IDEStatusLogFileCreationError, e)
		exitcode.Set(1)
		return
	}
	var logger = logctx.Use(ctx).Sugar()

	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		reportIdeError(model.IDEStatusAPIFail, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}

	if cliTaskIdOverride != "" {
		logger.Infof("CLI task id override: %s", cliTaskIdOverride)
		cf := config.RepoConfig{TaskId: cliTaskIdOverride}
		if e := cf.Validate(); e != nil {
			cv.DisplayBadTaskId(ctx)
			logger.Error(e)
			exitcode.Set(1)
			return
		}
		e = config.WriteRepoConfig(ctx, scanDir, accessType, config.RepoConfig{TaskId: cliTaskIdOverride})
		if e != nil {
			cv.DisplayInitializeFailed(ctx, e)
			logger.Error(e)
			exitcode.Set(1)
			return
		}
	}

	task, e := scan(ctx, scanDir, accessType)
	if e != nil {
		autoReportIde(e)
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

func autoReportIde(e error) {
	if errors.Is(e, api.ErrTaskNotFound) {
		reportIdeError(model.IDEStatusTaskNotExists, e)
		return
	}
	if errors.Is(e, api.ErrTokenInvalid) {
		reportIdeError(model.IDEStatusTokenInvalid, e)
		return
	}
	if errors.Is(e, api.ErrServerFail) {
		reportIdeError(model.IDEStatusServerFail, e)
		return
	}
	if errors.Is(e, api.ErrGeneralError) {
		reportIdeError(model.IDEStatusGeneralAPIError, e)
		return
	}
	if errors.Is(e, api.ErrRequest) {
		reportIdeError(model.IDEStatusAPIFail, e)
		return
	}
	reportIdeError(model.IDEStatusUnknownError, e)
}

func reportIdeError(status model.IDEStatus, e error) {
	resp := ideErrorResp{
		ErrCode: status,
		ErrMsg:  status.String(),
	}
	if e != nil {
		resp.ErrMsg = e.Error()
	}
	fmt.Println(string(must.A(json.MarshalIndent(resp, "", "  "))))
}

package binscan

import (
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/chunkupload"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

var cliIOTScan bool
var projectNameCli string

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "binscan <DIR>"
	c.Args = cobra.ExactArgs(1)
	c.Run = binScanRun
	c.Short = "Detects open source vulnerabilities by scanning binary files"
	c.Flags().BoolVar(&cliIOTScan, "iot", false, "IOT scan mode")
	c.Flags().StringVar(&projectNameCli, "project-name", "", "specify project name")
	return &c
}

func binScanRun(cmd *cobra.Command, args []string) {
	var (
		// workaround
		ctx      = ui.With(context.TODO(), ui.CLI)
		scanPath = args[0]
		e        error
	)
	// get absolute path and check if a directory
	scanPath, e = filepath.Abs(scanPath)
	if e != nil {
		cv.DisplayScanInvalidPath(ctx, e)
	}
	if !utils.IsFile(scanPath) {
		cv.DisplayScanInvalidPathMustFile(ctx, nil)
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

	e = binScan(ctx, scanPath)
	if errors.Is(e, inspector.ErrNoWait) {
		return
	}
	if e != nil {
		logger.Error(e)
		exitcode.Set(1)
		return
	}

}

func binScan(ctx context.Context, scanPath string) error {
	var mode = model.ScanModeBinary
	if cliIOTScan {
		mode = model.ScanModeIot
	}
	var subtaskName = filepath.Base(scanPath)
	taskResp, e := api.CreateSubTask(api.DefaultClient(), &api.CreateSubTaskRequest{
		AccessType:  model.AccessTypeCli,
		ScanMode:    mode,
		SubtaskName: subtaskName,
		Dir:         filepath.Dir(scanPath),
		ProjectName: projectNameCli,
	})
	if e != nil {
		cv.DisplayCreateSubtaskErr(ctx, e)
		return e
	}
	cv.DisplayAlertMessage(ctx, taskResp.AlertMessage)
	cv.DisplaySubtaskCreated(ctx, taskResp.ProjectsName, taskResp.SubtaskID)
	cv.DisplayUploading(ctx)
	defer cv.DisplayStatusClear(ctx)
	e = chunkupload.UploadFile(ctx, scanPath, chunkupload.Params{
		TaskId:    taskResp.TaskID,
		SubtaskId: taskResp.SubtaskID,
	})
	if e != nil {
		cv.DisplayUploadErr(ctx, e)
		return e
	}

	e = api.StartCheck(api.DefaultClient(), taskResp.SubtaskID)
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return e
	}
	if env.NoWait {
		return inspector.ErrNoWait
	}
	cv.DisplayWaitingResponse(ctx)
	// query result
	var result *model.ScanResultResponse
	result, e = api.QueryResult(ctx, api.DefaultClient(), taskResp.SubtaskID)
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
	}
	cv.DisplayStatusClear(ctx)
	cv.DisplayScanResultSummary(ctx, result.RelyNum, result.LeakNum, len(result.VulnInfoMap))

	return nil
}

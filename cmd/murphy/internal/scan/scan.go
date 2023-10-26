package scan

import (
	"context"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/chunkupload"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/collect"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/gitinfo"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ref"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"path/filepath"
)

func scan(ctx context.Context, dir string, accessType model.AccessType, mode model.ScanMode) (*model.ScanTask, error) {
	must.NotNil(ctx)
	must.True(filepath.IsAbs(dir))
	must.True(accessType.Valid())

	var (
		e      error
		logger = logctx.Use(ctx).Sugar()
	)
	logger.Infof("Scan dir: %s", dir)
	cv.DisplayScanning(ctx)
	defer cv.DisplayStatusClear(ctx)

	var createSubtask api.CreateSubTaskRequest
	createSubtask.SubtaskName = filepath.Base(dir)
	createSubtask.AccessType = accessType
	createSubtask.ScanMode = mode
	createSubtask.Dir = dir
	createSubtask.IsBuild = !noBuild
	createSubtask.IsDeep = isDeep
	createSubtask.ProjectName = projectNameCli

	// get git info
	var gitSummary *gitinfo.Summary
	gitSummary, e = gitinfo.GetSummary(ctx, dir)
	if e != nil {
		logger.Warnf("get git info failed: %v", e)
	} else {
		assignGitInfoToCreateSubtaskReq(&createSubtask, gitSummary)
	}

	// call API
	createTaskResp, e := api.CreateSubTask(api.DefaultClient(), &createSubtask)
	if errors.Is(e, api.ErrTLSError) {
		cv.DisplayTLSNotice(ctx)
		return nil, e
	}
	if e != nil {
		cv.DisplayCreateSubtaskErr(ctx, e)
		return nil, e
	}
	logger.Infof("subtask created, %s / %s", createTaskResp.TaskID, createTaskResp.SubtaskID)
	if onlyTaskId {
		fmt.Println("subtask_id=", createTaskResp.SubtaskID)
	}
	cv.DisplayAlertMessage(ctx, createTaskResp.AlertMessage)
	cv.DisplaySubtaskCreated(ctx, createTaskResp.ProjectsName, createTaskResp.SubtaskID)
	cv.DisplayReportUrl(ctx, api.DefaultClient().BaseURLText(), createTaskResp.TaskID, createTaskResp.SubtaskID)

	// create task object
	task := &model.ScanTask{
		Ctx:         ctx,
		Mode:        mode,
		AccessType:  accessType,
		ProjectPath: dir,
		TaskId:      createTaskResp.TaskID,
		SubtaskId:   createTaskResp.SubtaskID,
		SubtaskName: createSubtask.SubtaskName,
	}

	ctx = model.WithScanTask(ctx, task)
	if task.Mode == model.ScanModeSource && !isDeep {
		// do scan
		e = inspector.ManagedInspect(ctx)
		if e != nil {
			cv.DisplayScanFailed(ctx, e)
			return nil, e
		}

		// submit SBOM
		e = api.SubmitSBOM(api.DefaultClient(), task.SubtaskId, task.Modules, task.CodeFragments)
		if e != nil {
			cv.DisplaySubmitSBOMErr(ctx, e)
			return nil, e
		}
	} else {
		cv.DisplayUploading(ctx)
		e = chunkupload.UploadDirectory(ctx, task.ProjectPath, chunkupload.DiscardDot, chunkupload.Params{
			SubtaskId: task.SubtaskId,
		})
		if e != nil {
			cv.DisplayUploadErr(ctx, e)
			return nil, e
		}
	}

	// start check
	e = api.StartCheck(api.DefaultClient(), task.SubtaskId)
	if e != nil {
		cv.DisplaySubmitSBOMErr(ctx, e)
		return nil, e
	}

	if onlyTaskId {
		return task, nil
	}
	// 收集贡献者信息
	cu, e := collect.CollectDir(ctx, task.ProjectPath)
	if e != nil {
		logger.Warn("收集贡献者信息失败", zap.Error(e))
	} else {
		cu.RepoInfo.SubtaskId = createTaskResp.SubtaskID
		api.ReportCollectedContributors(ctx, api.DefaultClient(), cu)
		logger.Info("报送贡献者信息成功")
	}
	if env.NoWait {
		return nil, inspector.ErrNoWait
	}
	// query result
	cv.DisplayWaitingResponse(ctx)
	defer cv.DisplayStatusClear(ctx)
	var result *model.ScanResultResponse
	result, e = api.QueryResult(ctx, api.DefaultClient(), task.SubtaskId)
	task.Result = result
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return nil, e
	}
	cv.DisplayStatusClear(ctx)
	cv.DisplayScanResultSummary(ctx, result.RelyNum, result.LeakNum, len(result.VulnInfoMap))

	return task, nil
}

func assignGitInfoToCreateSubtaskReq(createSubtask *api.CreateSubTaskRequest, gitSummary *gitinfo.Summary) {
	createSubtask.Addr = ref.OmitZero(gitSummary.RemoteAddr)
	createSubtask.Author = ref.OmitZero(gitSummary.AuthorEmail)
	createSubtask.Branch = ref.OmitZero(gitSummary.BranchName)
	createSubtask.PushTime = ref.OmitZero(gitSummary.CommitTime)
	createSubtask.Commit = ref.OmitZero(gitSummary.CommitHash)
	createSubtask.Message = ref.OmitZero(gitSummary.CommitMessage)
}

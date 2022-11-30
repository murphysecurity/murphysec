package scan

import (
	"context"
	"errors"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/gitinfo"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ref"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"path/filepath"
)

func scan(ctx context.Context, dir string, accessType model.AccessType) (*model.ScanTask, error) {
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
	createSubtask.ScanMode = model.ScanModeSource

	// get repo config
	var shouldWriteConfig = false
	var repoConfig *config.RepoConfig
	repoConfig, e = config.ReadRepoConfig(ctx, dir, accessType)
	if e == nil {
		createSubtask.TaskID = ref.OmitZero(repoConfig.TaskId)
		if repoConfig.TaskId == "" {
			logger.Infof("task id not set, will be written soom")
			shouldWriteConfig = true
		}
	}
	if errors.Is(e, config.ErrRepoConfigNotFound) {
		logger.Infof("config not found, will be written soon")
		shouldWriteConfig = true
	}

	// get git info
	var gitSummary *gitinfo.Summary
	gitSummary, e = gitinfo.GetSummary(ctx, dir)
	if e != nil {
		logger.Warnf("get git info failed: %v", e)
	} else {
		createSubtask.Addr = ref.OmitZero(gitSummary.RemoteAddr)
		createSubtask.Author = ref.OmitZero(gitSummary.AuthorEmail)
		createSubtask.Branch = ref.OmitZero(gitSummary.BranchName)
		createSubtask.PushTime = ref.OmitZero(gitSummary.CommitTime)
		createSubtask.Commit = ref.OmitZero(gitSummary.CommitHash)
		createSubtask.Message = ref.OmitZero(gitSummary.CommitMessage)
	}

	// call API
	var createTaskResp *api.CreateSubTaskResponse
	createTaskResp, e = api.CreateSubTask(api.DefaultClient(), &createSubtask)
	if errors.Is(e, api.ErrTLSError) {
		cv.DisplayTLSNotice(ctx)
		return nil, e
	}
	if e != nil {
		cv.DisplayCreateSubtaskErr(ctx, e)
		return nil, e
	}
	logger.Infof("subtask created, %s / %s", createTaskResp.TaskID, createTaskResp.SubtaskID)
	cv.DisplaySubtaskCreated(ctx, createTaskResp.ProjectsName, createTaskResp.TaskName, createTaskResp.TaskID, createSubtask.SubtaskName, createTaskResp.SubtaskID)
	if shouldWriteConfig {
		logger.Infof("creating repo config...")
		e = config.WriteRepoConfig(ctx, dir, accessType, config.RepoConfig{TaskId: createTaskResp.TaskID})
		if e != nil {
			logger.Warnf("repo config: %v", e)
		}
	}

	// create task object
	task := &model.ScanTask{
		Ctx:         ctx,
		Mode:        model.ScanModeSource,
		AccessType:  accessType,
		ProjectPath: dir,
		TaskId:      createTaskResp.TaskID,
		SubtaskId:   createTaskResp.SubtaskID,
		SubtaskName: createSubtask.SubtaskName,
	}

	// do scan
	ctx = model.WithScanTask(ctx, task)
	e = inspector.ManagedInspect(ctx)
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return nil, e
	}

	// submit SBOM
	e = api.SubmitSBOM(api.DefaultClient(), task.SubtaskId, task.Modules)
	if e != nil {
		cv.DisplaySubmitSBOMErr(ctx, e)
		return nil, e
	}

	// start check
	e = api.StartCheck(api.DefaultClient(), task.SubtaskId)
	if e != nil {
		cv.DisplaySubmitSBOMErr(ctx, e)
		return nil, e
	}

	cv.DisplayWaitingResponse(ctx)
	defer cv.DisplayStatusClear(ctx)
	// query result
	var result *model.ScanResultResponse
	result, e = api.QueryResult(ctx, api.DefaultClient(), task.SubtaskId)
	task.Result = result
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return nil, e
	}
	cv.DisplayStatusClear(ctx)
	cv.DisplayScanResultSummary(ctx, result.RelyNum, result.LeakNum)

	return task, nil
}

package scan

import (
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/gitinfo"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ref"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"path/filepath"
)

func scan(ctx context.Context, dir string, accessType model.AccessType) error {
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
		return e
	}
	if e != nil {
		cv.DisplayCreateSubtaskErr(ctx, e)
		return e
	}
	cv.DisplaySubtaskCreated(ctx, createTaskResp.ProjectsName, createTaskResp.TaskName, createTaskResp.TaskID)
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
	}

	// do scan
	ctx = model.WithScanTask(ctx, task)
	e = inspector.ManagedInspect(ctx)
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return e
	}

	// submit SBOM
	e = api.SubmitSBOM(api.DefaultClient(), task.SubtaskId, task.Modules)
	if e != nil {
		cv.DisplaySubmitSBOMErr(ctx, e)
		return e
	}

	// start check
	e = api.StartCheck(api.DefaultClient(), task.SubtaskId)
	if e != nil {
		cv.DisplaySubmitSBOMErr(ctx, e)
		return e
	}

	cv.DisplayWaitingResponse(ctx)
	defer cv.DisplayStatusClear(ctx)
	// query result
	var result *model.ScanResultResponse
	result, e = api.QueryResult(ctx, api.DefaultClient(), task.SubtaskId)
	task.Result = result
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return e
	}
	cv.DisplayStatusClear(ctx)
	cv.DisplayScanResultSummary(ctx, result.RelyNum, result.LeakNum)

	return nil
}

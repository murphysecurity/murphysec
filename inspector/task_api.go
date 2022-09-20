package inspector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/conf"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"os"
	"strings"
)

var _CPPModuleUUID = uuid.Must(uuid.Parse("794a5c39-ce6b-458e-8f26-ff26298bab09"))

func submitModuleInfoApi(ctx context.Context) error {
	task := model.UseScanTask(ctx)
	req := new(api.SendDetectRequest)
	req.TaskInfo = task.TaskId
	for _, it := range task.Modules {
		req.Modules = append(req.Modules, api.VoModule{
			Dependencies:   it.Dependencies,
			FileHashList:   nil,
			Language:       it.Language,
			Name:           it.Name,
			PackageManager: it.PackageManager,
			RelativePath:   it.FilePath,
			RuntimeInfo:    it.RuntimeInfo,
			Version:        it.Version,
			ModuleUUID:     it.UUID,
			ModuleType:     api.ModuleTypeVersion,
			ScanStrategy:   string(it.ScanStrategy),
		})
	}
	if len(task.FileHashes) != 0 && env.AllowFileHash {
		list := make([]api.VoFileHash, 0)
		for _, it := range task.FileHashes {
			for _, hash := range it.Hash {
				list = append(list, api.VoFileHash{
					Path: it.Path,
					Hash: hash,
				})
			}
		}
		req.Modules = append(req.Modules, api.VoModule{
			FileHashList: list,
			Language:     model.Cxx,
			ModuleType:   api.ModuleTypeFileHash,
			ModuleUUID:   _CPPModuleUUID,
		})
	}
	if e := api.SendDetect(req); e != nil {
		Logger.Error("Module data commit failed", zap.Error(e))
		return e
	}
	Logger.Info("Module data committed", zap.Int("total_module", len(task.Modules)))
	return nil
}

func createTaskApi(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	req := &api.CreateTaskRequest{
		CliVersion:      version.Version(),
		TaskType:        scanTask.TaskType,
		UserAgent:       version.UserAgent(),
		CmdLine:         strings.Join(os.Args, " "),
		ApiToken:        conf.APIToken(),
		ProjectName:     scanTask.ProjectName,
		TargetAbsPath:   scanTask.ProjectDir,
		ProjectType:     scanTask.ProjectType,
		ContributorList: scanTask.ContributorList,
		ProjectId:       scanTask.ProjectId,
	}
	if g := scanTask.GitInfo; g != nil {
		v := &api.VoGitInfo{
			Commit:        g.HeadCommitHash,
			GitRef:        g.HeadRefName,
			GitRemoteUrl:  g.RemoteURL,
			CommitMessage: g.CommitMsg,
			CommitEmail:   g.CommitterEmail,
			CommitTime:    g.CommitTime,
		}
		req.GitInfo = v
	}
	if env.SpecificProjectName != "" {
		// force set project dir, in order to create new project
		req.TargetAbsPath = fmt.Sprintf(`/%s`, env.SpecificProjectName)
	}
	var res *api.CreateTaskResponse
	res, e = api.CreateTask(req)
	if e != nil {
		Logger.Error("Task create failed", zap.Error(e))
		return e
	}
	scanTask.TaskId = res.TaskInfo
	scanTask.TotalContributors = res.TotalContributors
	scanTask.ProjectId = res.ProjectId
	scanTask.Username = res.Username
	Logger.Info("Task created", zap.Any("task_id", scanTask.TaskId), zap.Any("project_id", scanTask.ProjectId))
	if res.AlertMessage != "" {
		ui.Display(res.AlertLevel, res.AlertMessage)
	}
	return
}

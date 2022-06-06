package inspector

import (
	"context"
	"github.com/google/uuid"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/env"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/version"
	"os"
	"strings"
)

func createTask(ctx context.Context) error {
	c := model.UseScanTask(ctx)
	req := &api.CreateTaskRequest{
		CliVersion:      version.Version(),
		TaskType:        c.TaskType,
		UserAgent:       version.UserAgent(),
		CmdLine:         strings.Join(os.Args, " "),
		ApiToken:        conf.APIToken(),
		ProjectName:     c.ProjectName,
		TargetAbsPath:   c.ProjectDir,
		ProjectType:     c.ProjectType,
		ContributorList: c.ContributorList,
		ProjectId:       c.ProjectId,
	}

	if g := c.GitInfo; g != nil {
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
	if res, e := api.CreateTask(req); e == nil {
		c.TaskId = res.TaskInfo
		c.TotalContributors = res.TotalContributors
		c.ProjectId = res.ProjectId
		logger.Info.Println("task created, id:", res.TaskInfo)
		return nil
	} else {
		logger.Warn.Println("task create failed", e.Error())
		return e
	}
}

var CPPModuleUUID = uuid.Must(uuid.Parse("794a5c39-ce6b-458e-8f26-ff26298bab09"))

func submitModuleInfo(ctx context.Context) error {
	task := model.UseScanTask(ctx)
	req := new(api.SendDetectRequest)
	req.TaskInfo = task.TaskId
	for _, it := range task.Modules {
		req.Modules = append(req.Modules, api.VoModule{
			Dependencies:   it.Dependencies,
			FileHashList:   nil,
			Language:       it.Language,
			Name:           it.Name,
			PackageFile:    it.PackageFile,
			PackageManager: it.PackageManager,
			RelativePath:   it.FilePath,
			RuntimeInfo:    it.RuntimeInfo,
			Version:        it.Version,
			ModuleUUID:     it.UUID,
			ModuleType:     api.ModuleTypeVersion,
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
			ModuleUUID:   CPPModuleUUID,
		})
	}
	if e := api.SendDetect(req); e != nil {
		logger.Err.Println("send module info failed.", e.Error())
		return e
	}
	return nil
}

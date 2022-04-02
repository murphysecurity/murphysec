package inspector

import (
	"github.com/google/uuid"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/env"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/version"
	"os"
	"strings"
)

func createTask(ctx *ScanContext) error {
	req := &api.CreateTaskRequest{
		CliVersion:      version.Version(),
		TaskType:        ctx.TaskType,
		UserAgent:       version.UserAgent(),
		CmdLine:         strings.Join(os.Args, " "),
		ApiToken:        conf.APIToken(),
		ProjectName:     ctx.ProjectName,
		TargetAbsPath:   ctx.ProjectDir,
		ProjectType:     ctx.ProjectType,
		ContributorList: ctx.ContributorList,
		ProjectId:       ctx.ProjectId,
	}
	req.GitInfo = ctx.GitInfo.ApiVo()
	logger.Info.Printf("create task: %#v", ctx)
	if res, e := api.CreateTask(req); e == nil {
		ctx.TaskId = res.TaskInfo
		ctx.TotalContributors = res.TotalContributors
		ctx.ProjectId = res.ProjectId
		logger.Info.Println("task created, id:", res.TaskInfo)
		return nil
	} else {
		logger.Warn.Println("task create failed", e.Error())
		return e
	}
}

var CPPModuleUUID = uuid.Must(uuid.Parse("794a5c39-ce6b-458e-8f26-ff26298bab09"))

func submitModuleInfo(ctx *ScanContext) error {
	req := new(api.SendDetectRequest)
	req.TaskInfo = ctx.TaskId
	for _, it := range ctx.ManagedModules {
		req.Modules = append(req.Modules, *it.ApiVo())
	}
	if len(ctx.FileHashes) != 0 && env.AllowFileHash {
		list := make([]api.VoFileHash, 0)
		for _, it := range ctx.FileHashes {
			list = append(list, api.VoFileHash{Hash: it})
		}
		req.Modules = append(req.Modules, api.VoModule{
			FileHashList: list,
			Language:     "C/C++",
			ModuleType:   "file_hash",
			ModuleUUID:   CPPModuleUUID,
		})
	}
	if e := api.SendDetect(req); e != nil {
		logger.Err.Println("send module info failed.", e.Error())
		return e
	}
	return nil
}

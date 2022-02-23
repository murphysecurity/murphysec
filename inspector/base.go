package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"strings"
	"time"
)

var ErrNoEngineMatched = errors.New("ErrNoEngineMatched")
var ErrAPITokenInvalid = errors.New("ErrAPITokenInvalid")
var ErrNoModule = errors.New("ErrNoModule")

type ScanContext struct {
	GitInfo        *GitInfo
	ProjectName    string
	TaskType       api.InspectTaskType
	ProjectDir     string
	ManagedModules []base.Module
	StartTime      time.Time
	TaskId         string
	FileHashes     []string
}

func createTaskContext(baseDir string, taskType api.InspectTaskType) *ScanContext {
	ctx := readProjectInfo(baseDir)
	ctx.TaskType = taskType
	ctx.StartTime = time.Now()
	return ctx
}

func createTask(ctx *ScanContext) error {
	req := &api.CreateTaskRequest{
		CliVersion:    version.Version(),
		TaskType:      ctx.TaskType,
		UserAgent:     version.UserAgent(),
		CmdLine:       strings.Join(os.Args, " "),
		ApiToken:      conf.APIToken(),
		ProjectName:   ctx.ProjectName,
		TargetAbsPath: ctx.ProjectDir,
	}
	req.GitInfo = ctx.GitInfo.ApiVo()
	logger.Info.Printf("create task: %#v", ctx)
	if taskId, e := api.CreateTask(req); e == nil {
		ctx.TaskId = *taskId
		logger.Info.Println("task created, id:", *taskId)
		return nil
	} else {
		logger.Warn.Println("task create failed", e.Error())
		return e
	}
}

func Scan(dir string, source api.InspectTaskType) (interface{}, error) {
	ctx := createTaskContext(dir, source)
	if e := createTask(ctx); e != nil {
		logger.Err.Println("Create task failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		return nil, e
	}

	if source == api.TaskTypeCli {
		fmt.Printf("项目创建成功，项目名称：%s\n", ctx.ProjectName)
	}

	if e := managedInspectScan(ctx); e != nil {
		if source == api.TaskTypeCli {
			fmt.Println("受管理扫描失败，执行文件哈希扫描")
		}
		logger.Err.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		// if managed inspect failed, start file hash scan
		FileHashScan(ctx)
		if e == ErrNoEngineMatched && len(ctx.FileHashes) == 0 {
			if source == api.TaskTypeCli {
				fmt.Println("扫描器无法支持当前项目")
			}
			return nil, ErrNoEngineMatched
		}
	}
	if source == api.TaskTypeCli {
		fmt.Printf("共扫描到模块：%d个\n，正在传输数据", len(ctx.ManagedModules))
	}
	// merge module info

	req := new(api.SendDetectRequest)
	req.TaskInfo = ctx.TaskId
	for _, it := range ctx.ManagedModules {
		req.Modules = append(req.Modules, *it.ApiVo())
	}
	if len(ctx.FileHashes) != 0 {
		req.Modules = append(req.Modules, api.VoModule{
			Hash:       ctx.FileHashes,
			Language:   "C/C++",
			ModuleType: "file_hash",
		})
	}

	if e := api.SendDetect(req); e != nil {
		logger.Err.Println("send module info failed.", e.Error())
		return nil, e
	}

	if source == api.TaskTypeCli {
		fmt.Println("检测中...")
	}

	if e := api.StartCheck(ctx.TaskId); e != nil {
		logger.Err.Println("send start check command failed.", e.Error())
		return nil, e
	}

	resp, e := api.QueryResult(ctx.TaskId)
	if e != nil {
		logger.Err.Println("query result failed.", e.Error())
		return nil, e
	}
	// todo: resp

	if source == api.TaskTypeCli {
		fmt.Printf("项目扫描成功，依赖数：%d，漏洞数：%d\n", resp.DependenciesCount, resp.IssuesCompsCount)
	} else if source == api.TaskTypeJenkins || source == api.TaskTypeIdea {
		fmt.Println(must.Byte(json.Marshal(mapForIdea(resp))))
	}
	return nil, nil
}

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
	"path/filepath"
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

func commitModuleInfo(ctx *ScanContext) error {
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
		return e
	}
	return nil
}

func displayTaskCreating(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("正在创建扫描任务，请稍候，项目名称：", ctx.ProjectName)
	}
}

func displayTaskCreated(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("扫描任务已创建")
	}
}

func displayManagedScanning(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("正在执行扫描")
	}
}

func shouldUploadFile(ctx *ScanContext) bool {
	if len(ctx.ManagedModules) == 0 {
		return true
	}
	for _, it := range ctx.ManagedModules {
		if it.PackageManager == "Maven" {
			return true
		}
	}
	return false
}

func Scan(dir string, source api.InspectTaskType, deepScan bool) (interface{}, error) {
	{
		info, e := os.Stat(dir)
		if e != nil || !info.IsDir() {
			if source == api.TaskTypeCli {
				fmt.Println("项目目录不存在或无效")
			}
			return nil, errors.New("Invalid project dir")
		}
	}
	ctx := createTaskContext(dir, source)

	displayTaskCreating(ctx)
	if e := createTask(ctx); e != nil {
		logger.Err.Println("Create task failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		return nil, e
	}
	displayTaskCreated(ctx)

	displayManagedScanning(ctx)
	if e := managedInspectScan(ctx); e != nil {
		logger.Debug.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%v", e)
		// if managed inspect failed, start file hash scan
		FileHashScan(ctx)
		if e == ErrNoEngineMatched && len(ctx.FileHashes) == 0 {
			if source == api.TaskTypeCli {
				fmt.Println("扫描器无法支持当前项目")
			}
			return nil, ErrNoEngineMatched
		}
	}
	if e := commitModuleInfo(ctx); e != nil {
		if source == api.TaskTypeCli {
			fmt.Println("提交模块信息失败", e.Error())
		}
		logger.Debug.Printf("%+v", e)
		logger.Err.Println(e.Error())
	}

	if deepScan && shouldUploadFile(ctx) {
		logger.Info.Printf("deep scan enabled, upload source code")
		fmt.Println("正在上传文件到服务端以进行深度检测")
		if e := UploadCodeFile(ctx); e != nil {
			if source == api.TaskTypeCli {
				fmt.Println("深度检测上传文件失败！")
			}
		} else {
			if source == api.TaskTypeCli {
				fmt.Println("深度检测上传文件成功！")
			}
		}
	}

	if source == api.TaskTypeCli {
		fmt.Println("检测中，等待返回结果...")
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

	if source == api.TaskTypeCli {
		fmt.Printf("项目扫描成功，依赖数：%d，漏洞数：%d\n", resp.DependenciesCount, resp.IssuesCompsCount)
	} else if source == api.TaskTypeJenkins || source == api.TaskTypeIdea {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(resp)))))
	}
	return nil, nil
}

func ScannerScan(dir string) {
	ctx := createTaskContext(must.String(filepath.Abs(dir)), api.TaskTypeIdea)
	if e := managedInspectScan(ctx); e != nil {
		logger.Err.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%+v", e)
	}
	if ctx.ManagedModules == nil {
		ctx.ManagedModules = []base.Module{}
	}
	fmt.Println(string(must.Byte(json.Marshal(ctx.ManagedModules))))
}

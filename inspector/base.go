package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	"io/fs"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"strconv"
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
	ProjectType    string
}

func (s *ScanContext) AddManagedModule(module base.Module) {
	module.UUID = uuid.New()
	s.ManagedModules = append(s.ManagedModules, module)
}

func createTaskContext(baseDir string, taskType api.InspectTaskType) *ScanContext {
	ctx := readProjectInfo(baseDir)
	ctx.TaskType = taskType
	ctx.StartTime = time.Now()
	if ctx.GitInfo != nil {
		ctx.ProjectType = "Local"
	} else {
		ctx.ProjectType = "Git"
	}
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
		ProjectType:   ctx.ProjectType,
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

var CPPModuleUUID = uuid.Must(uuid.Parse("794a5c39-ce6b-458e-8f26-ff26298bab09"))

func commitModuleInfo(ctx *ScanContext) error {
	req := new(api.SendDetectRequest)
	req.TaskInfo = ctx.TaskId
	for _, it := range ctx.ManagedModules {
		req.Modules = append(req.Modules, *it.ApiVo())
	}
	if len(ctx.FileHashes) != 0 {
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

func checkProjectDirAvail(dir string) bool {
	info, e := os.Stat(dir)
	return e == nil && info.IsDir()
}

func Scan(dir string, source api.InspectTaskType, deepScan bool) (interface{}, error) {
	ui := display.NONE
	if source == api.TaskTypeCli {
		ui = display.CLI
		EnableANSI()
	}
	if !checkProjectDirAvail(dir) {
		ui.Display(display.MsgError, "项目目录不存在或无效")
		return nil, errors.New("Invalid project dir")
	}
	ctx := createTaskContext(dir, source)
	ui.Display(display.MsgInfo, fmt.Sprint("项目名称：", ctx.ProjectName))
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")
	if e := createTask(ctx); e != nil {
		ui.Display(display.MsgError, fmt.Sprint("项目创建失败：", e.Error()))
		logger.Err.Println("Create task failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		return nil, e
	}
	ui.Display(display.MsgInfo, fmt.Sprint("项目创建成功，项目唯一标识：", ctx.TaskId))

	ui.UpdateStatus(display.StatusRunning, "正在进行扫描...")
	if e := managedInspectScan(ctx); e != nil {
		logger.Debug.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%v", e)
	}

	{
		enableCxx := false
		filepath.Walk(ctx.ProjectDir, func(path string, info fs.FileInfo, err error) error {
			if enableCxx {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			enableCxx = CxxExtSet[filepath.Ext(info.Name())]
			return nil
		})
		if enableCxx {
			FileHashScan(ctx)
		}
	}

	ui.UpdateStatus(display.StatusRunning, "项目扫描结束，正在提交信息...")
	if e := commitModuleInfo(ctx); e != nil {
		ui.Display(display.MsgError, fmt.Sprint("信息提交失败：", e.Error()))
		logger.Debug.Printf("%+v", e)
		logger.Err.Println(e.Error())
	}
	if deepScan && shouldUploadFile(ctx) {
		logger.Info.Printf("deep scan enabled, upload source code")
		ui.UpdateStatus(display.StatusRunning, "正在上传文件到服务端以进行深度检测")
		if e := UploadCodeFile(ctx); e != nil {
			ui.Display(display.MsgError, "上传文件失败："+e.Error())
		} else {
			ui.Display(display.MsgInfo, "上传文件成功")
		}
		ui.ClearStatus()
	}

	ui.UpdateStatus(display.StatusRunning, "检测中，等待返回结果...")

	if e := api.StartCheck(ctx.TaskId); e != nil {
		ui.Display(display.MsgError, "启动检测失败："+e.Error())
		logger.Err.Println("send start check command failed.", e.Error())
		return nil, e
	}
	ui.ClearStatus()
	resp, e := api.QueryResult(ctx.TaskId)
	ui.ClearStatus()
	if e != nil {
		ui.Display(display.MsgError, "获取检测结果失败："+e.Error())
		logger.Err.Println("query result failed.", e.Error())
		return nil, e
	}
	ui.Display(display.MsgNotice, fmt.Sprint("项目扫描成功，依赖数：", termenv.String(strconv.Itoa(resp.DependenciesCount)).Foreground(termenv.ANSIBrightCyan), "，漏洞数：", termenv.String(strconv.Itoa(resp.IssuesCompsCount)).Foreground(termenv.ANSIBrightRed)))
	if source == api.TaskTypeJenkins || source == api.TaskTypeIdea {
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

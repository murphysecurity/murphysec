package inspector

import (
	"context"
	"errors"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"strconv"
)

func Scan(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.Display(display.MsgInfo, fmt.Sprint("项目名称：", scanTask.ProjectName))
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")

	if e := createTask(ctx); e != nil {
		logger.Err.Println("Create task failed.", e.Error())
		logger.Debug.Printf("%+v", e)
		ui.Display(display.MsgError, fmt.Sprint("项目创建失败"))
		if errors.Is(api.ErrTokenInvalid, e) {
			ui.Display(display.MsgError, "当前 Token 无效")
		} else {
			ui.Display(display.MsgError, e.Error())
		}
		return e
	}

	ui.Display(display.MsgInfo, fmt.Sprint("项目创建成功"))
	ui.UpdateStatus(display.StatusRunning, "正在进行扫描...")

	if e := managedInspect(ctx); e != nil {
		logger.Debug.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%v", e)
		return e
	}

	if env.AllowFileHash && len(scanTask.Modules) == 0 {
		logger.Info.Println("File hash scanning...")
		if e := FileHashScan(ctx); e != nil {
			logger.Err.Println("FileHash calc failed.", e.Error())
			ui.Display(display.MsgInfo, "文件哈希计算失败："+e.Error())
		}
	}

	if env.AllowDeepScan && scanTask.EnableDeepScan {
		logger.Info.Println("DeepScan......")
		ui.Display(display.MsgInfo, "正在上传代码进行深度检测")
		ui.UpdateStatus(display.StatusRunning, "代码上传中")
		if e := UploadCodeFile(ctx); e != nil {
			logger.Err.Println("Code upload failed.", e.Error())
			ui.Display(display.MsgError, "代码上传失败："+e.Error())
		}
	}

	ui.UpdateStatus(display.StatusRunning, "项目扫描结束，正在提交信息...")
	if e := submitModuleInfo(ctx); e != nil {
		ui.Display(display.MsgError, fmt.Sprint("信息提交失败：", e.Error()))
		logger.Debug.Printf("%+v", e)
		logger.Err.Println(e.Error())
		return e
	}

	if e := api.StartCheck(scanTask.TaskId); e != nil {
		ui.Display(display.MsgError, "启动检测失败："+e.Error())
		logger.Err.Println("send start check command failed.", e.Error())
		return e
	}
	ui.ClearStatus()
	resp, e := api.QueryResult(scanTask.TaskId)
	ui.ClearStatus()
	if e != nil {
		ui.Display(display.MsgError, "获取检测结果失败："+e.Error())
		logger.Err.Println("query result failed.", e.Error())
		return e
	}
	scanTask.ScanResult = resp
	totalDep := strconv.Itoa(scanTask.ScanResult.DependenciesCount)
	totalVuln := strconv.Itoa(scanTask.ScanResult.IssuesCompsCount)
	t := fmt.Sprint(
		"项目扫描完成，依赖数：",
		termenv.String(totalDep).Foreground(termenv.ANSIBrightCyan),
		"，漏洞数：",
		termenv.String(totalVuln).Foreground(termenv.ANSIBrightRed),
	)
	if scanTask.ScanResult.InspectReportUrl != "" {
		ui.Display(display.MsgNotice, fmt.Sprintf("检测报告详见：%s", scanTask.ScanResult.InspectReportUrl))
	}
	ui.Display(display.MsgNotice, t)

	return nil
}

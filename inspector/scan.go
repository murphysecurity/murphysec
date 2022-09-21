package inspector

import (
	"context"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
	"strconv"
)

func Scan(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.Display(display.MsgInfo, fmt.Sprint("项目名称：", scanTask.ProjectName))
	if e := createTaskC(ctx); e != nil {
		return e
	}

	ui.UpdateStatus(display.StatusRunning, "正在进行扫描...")

	if e := managedInspect(ctx); e != nil {
		Logger.Error("Managed inspect failed", zap.Error(e))
		return e
	}

	if env.AllowFileHash && !scanTask.EnableDeepScan {
		Logger.Info("File hash scanning")
		if e := FileHashScan(ctx); e != nil {
			Logger.Error("FileHash calc failed", zap.Error(e))
			ui.Display(display.MsgInfo, "文件哈希计算失败："+e.Error())
		}
	}

	if env.AllowDeepScan && scanTask.EnableDeepScan {
		Logger.Info("DeepScan")
		ui.Display(display.MsgInfo, "正在上传代码进行深度检测")
		ui.UpdateStatus(display.StatusRunning, "代码上传中")
		if e := UploadCodeFile(ctx); e != nil {
			Logger.Error("Code upload failed", zap.Error(e))
			ui.Display(display.MsgError, "代码上传失败："+e.Error())
		}
	}

	ui.Display(display.MsgInfo, "项目扫描完成")

	if err := submitModuleInfoC(ctx); err != nil {
		return err
	}

	if err := startCheckC(ctx); err != nil {
		return err
	}
	if e := queryResultC(ctx); e != nil {
		return e
	}
	totalDep := strconv.Itoa(scanTask.ScanResult.DependenciesCount)
	totalVuln := strconv.Itoa(scanTask.ScanResult.IssuesCompsCount)
	t := fmt.Sprint(
		"项目扫描完成，依赖数：",
		termenv.String(totalDep).Foreground(termenv.ANSIBrightCyan),
		"，漏洞数：",
		termenv.String(totalVuln).Foreground(termenv.ANSIBrightRed),
	)

	if scanTask.ScanResult.ReportURL() != "" {
		ui.Display(display.MsgNotice, fmt.Sprintf("检测报告详见：%s", scanTask.ScanResult.ReportURL()))
	}
	ui.Display(display.MsgNotice, t)

	return nil
}

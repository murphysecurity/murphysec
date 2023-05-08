package cv

import (
	"context"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/infra/ui"
	"strconv"
)

func DisplayInitializeFailed(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "初始化失败："+e.Error())
}

func DisplayScanInvalidPath(ctx context.Context, e error) {
	msg := "您指定了一个无效目录"
	if e != nil {
		msg += "：" + e.Error()
	}
	ui.Use(ctx).Display(ui.MsgError, msg)
}

func DisplayScanInvalidPathMustDir(ctx context.Context, e error) {
	msg := "您指定了一个无效路径，请指定一个目录"
	if e != nil {
		msg += "：" + e.Error()
	}
	ui.Use(ctx).Display(ui.MsgError, msg)
}

func DisplayScanInvalidPathMustFile(ctx context.Context, e error) {
	msg := "您指定了一个无效路径，请指定一个文件"
	if e != nil {
		msg += "：" + e.Error()
	}
	ui.Use(ctx).Display(ui.MsgError, msg)
}

func DisplayTLSNotice(ctx context.Context) {
	ui.Use(ctx).Display(ui.MsgWarn, "当前建立的网络连接不安全，通讯已断开，您可以通过 -x 或 --allow-insecure 选项忽略这个错误")
}

func DisplayCreateSubtaskErr(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "扫描失败："+e.Error())
}

func DisplayScanning(ctx context.Context) {
	ui.Use(ctx).UpdateStatus(ui.StatusRunning, "正在扫描...")
}

func DisplayWaitingResponse(ctx context.Context) {
	ui.Use(ctx).UpdateStatus(ui.StatusRunning, "正在等待服务器返回扫描结果...")
}

func DisplayReportUrl(ctx context.Context, serverUrl string, taskId, subtaskId string) {
	ui.Use(ctx).Display(ui.MsgInfo, fmt.Sprintf("%s/console/report/%s/%s?allow=1", serverUrl, taskId, subtaskId))
}

func DisplaySubtaskCreated(ctx context.Context, projectName string, subtaskId string) {
	ui.Use(ctx).Display(ui.MsgInfo, "项目名称："+projectName)
	ui.Use(ctx).Display(ui.MsgInfo, fmt.Sprintf("检测历史ID: %s", subtaskId))
	ui.Use(ctx).Display(ui.MsgInfo, "")
}

func DisplayScanFailed(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "扫描失败"+e.Error())
}

func DisplaySubmitSBOMErr(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "上传SBOM信息失败："+e.Error())
}

func DisplayStatusClear(ctx context.Context) {
	ui.Use(ctx).ClearStatus()
}

func DisplayScanResultSummary(ctx context.Context, totalDep int, totalVulnDep int, totalVuln int) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgNotice, fmt.Sprint(
		"项目扫描完成，依赖数：",
		termenv.String(strconv.Itoa(totalDep)).Foreground(termenv.ANSIBrightCyan),
		"，缺陷组件数：",
		termenv.String(strconv.Itoa(totalVulnDep)).Foreground(termenv.ANSIBrightRed),
		"，漏洞数",
		termenv.String(strconv.Itoa(totalVuln)).Foreground(termenv.ANSIBrightRed),
	))
}
func DisplayUploading(ctx context.Context) {
	ui.Use(ctx).UpdateStatus(ui.StatusRunning, "正在上传...")
}

func DisplayUploadErr(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "上传失败："+e.Error())
}

func DisplayAlertMessage(ctx context.Context, msg string) {
	if msg == "" {
		return
	}
	ui.Use(ctx).Display(ui.MsgNotice, msg)
}

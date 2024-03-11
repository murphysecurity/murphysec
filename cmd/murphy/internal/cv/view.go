package cv

import (
	"context"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
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

func DisplayReportUrl(ctx context.Context, response model.ScanResultResponse) {
	ui.Use(ctx).Display(ui.MsgInfo, "任务详情："+response.DetailURL)
	var s1 string
	if response.ExpireDay == -1 {
		s1 = "永久有效"
	} else {
		s1 = strconv.Itoa(response.ExpireDay) + "天有效"
	}
	var s2 string
	if response.AllowAction == 12 {
		s2 = "可匿名查看报告"
	} else {
		s2 = "加入团队后可编辑报告"
	}
	if response.ShareURL != "" {
		ui.Use(ctx).Display(ui.MsgInfo, fmt.Sprintf("分享链接(%s-%s)：%s", s2, s1, response.ShareURL))
	}
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

func DisplayStartCheckErr(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "开始检测失败："+e.Error())
}

func DisplayStatusClear(ctx context.Context) {
	ui.Use(ctx).ClearStatus()
}

func DisplayScanResultSummary(ctx context.Context, totalDep int, totalVulnDep int, totalVuln int) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgNotice, fmt.Sprint("项目扫描完成，依赖数：", ui.Term.String(strconv.Itoa(totalDep)).Foreground(termenv.ANSIBrightCyan), "，缺陷组件数：", ui.Term.String(strconv.Itoa(totalVulnDep)).Foreground(termenv.ANSIBrightRed), "，漏洞数", ui.Term.String(strconv.Itoa(totalVuln)).Foreground(termenv.ANSIBrightRed)))
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

package view

import (
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/display"
	"strconv"
)

func TLSAlert(ui display.UI, e error) {
	ui.Display(display.MsgError, "当前建立的网络连接不安全，您可以通过 -x 或 --allow-insecure 选项忽略这个错误")
	ui.Display(display.MsgError, e.Error())
}

func TokenInvalid(ui display.UI) {
	ui.Display(display.MsgError, "任务创建失败，Token 无效")
}

func TaskCreateFailed(ui display.UI, e error) {
	ui.Display(display.MsgError, fmt.Sprintf("任务创建失败："+e.Error()))
}

func TaskCreating(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")
	return func() { ui.ClearStatus() }
}

func ScanCompleteSubmitting(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "项目扫描结束，正在提交信息...")
	return func() { ui.ClearStatus() }
}

func SubmitError(ui display.UI, err error) {
	ui.Display(display.MsgError, "信息提交失败："+err.Error())
}

func WaitingServerResponse(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "正在等待服务器返回结果")
	return func() { ui.ClearStatus() }
}

func GetScanResultFailed(ui display.UI, e error) {
	ui.Display(display.MsgError, "获取扫描结果失败："+e.Error())
}

func StartingInspection(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "正在启动检测")
	return func() { ui.ClearStatus() }
}
func StartingInspectionFailed(ui display.UI, e error) {
	ui.Display(display.MsgError, "启动检测失败："+e.Error())
}

func DisplayScanResultSummary(ui display.UI, totalDep int, totalVuln int) {
	ui.Display(display.MsgNotice, fmt.Sprint(
		"项目扫描完成，依赖数：",
		termenv.String(strconv.Itoa(totalDep)).Foreground(termenv.ANSIBrightCyan),
		"，漏洞数：",
		termenv.String(strconv.Itoa(totalVuln)).Foreground(termenv.ANSIBrightRed),
	))
}

func DisplayScanResultReport(ui display.UI, r string) {
	if r == "" {
		return
	}
	ui.Display(display.MsgNotice, "检测报告详见："+r)
}

func ProjectName(ui display.UI, n string) {
	ui.Display(display.MsgInfo, "项目名称："+n)
}

func FileUploading(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "正在上传文件...")
	return func() { ui.ClearStatus() }
}

func FileUploadSucceeded(ui display.UI) {
	ui.Display(display.MsgInfo, "文件上传成功")
}

func FileUploadFailed(ui display.UI, e error) {
	ui.Display(display.MsgError, "文件上传失败："+e.Error())
}

func ProjectScanComplete(ui display.UI) {
	ui.Display(display.MsgInfo, "项目扫描完成")
}

func ProjectScanning(ui display.UI) func() {
	ui.UpdateStatus(display.StatusRunning, "正在进行扫描...")
	return func() { ui.ClearStatus() }
}

func HashingFileFailed(ui display.UI, e error) {
	ui.Display(display.MsgInfo, "文件哈希计算失败："+e.Error())
}

func CodeFileUploadingForDeep(ui display.UI) func() {
	ui.Display(display.MsgInfo, "正在上传代码进行深度检测")
	ui.UpdateStatus(display.StatusRunning, "代码上传中")
	return func() { ui.ClearStatus() }
}

func CodeFileUploadErr(ui display.UI, e error) {
	ui.Display(display.MsgError, "代码上传失败："+e.Error())
}

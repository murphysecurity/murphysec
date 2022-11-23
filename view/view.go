package view

import (
	"context"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/murphysecurity/murphysec/infra/ui"
	"strconv"
	"sync"
)

func TLSAlert(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "当前建立的网络连接不安全，您可以通过 -x 或 --allow-insecure 选项忽略这个错误")
	u.Display(ui.MsgError, e.Error())
}

func TokenInvalid(ctx context.Context) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "任务创建失败，Token 无效")
}

func TaskCreateFailed(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, fmt.Sprintf("任务创建失败："+e.Error()))
}

func TaskCreating(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "正在创建扫描任务，请稍候······")
	return func() { u.ClearStatus() }
}

func ScanCompleteSubmitting(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "项目扫描结束，正在提交信息...")
	return func() { u.ClearStatus() }
}

func SubmitError(ctx context.Context, err error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "信息提交失败："+err.Error())
}

func WaitingServerResponse(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "正在等待服务器返回结果")
	return func() { u.ClearStatus() }
}

func GetScanResultFailed(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "获取扫描结果失败："+e.Error())
}

func StartingInspection(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "正在启动检测")
	return func() { u.ClearStatus() }
}
func StartingInspectionFailed(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "启动检测失败："+e.Error())
}

func DisplayScanResultSummary(ctx context.Context, totalDep int, totalVuln int) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgNotice, fmt.Sprint(
		"项目扫描完成，依赖数：",
		termenv.String(strconv.Itoa(totalDep)).Foreground(termenv.ANSIBrightCyan),
		"，漏洞数：",
		termenv.String(strconv.Itoa(totalVuln)).Foreground(termenv.ANSIBrightRed),
	))
}

func DisplayScanResultReport(ctx context.Context, r string) {
	var u = ui.Use(ctx)
	if r == "" {
		return
	}
	u.Display(ui.MsgNotice, "检测报告详见："+r)
}

func ProjectName(ctx context.Context, n string) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgInfo, "项目名称："+n)
}

func FileUploading(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "正在上传文件...")
	return func() { u.ClearStatus() }
}

func FileUploadSucceeded(ctx context.Context) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgInfo, "文件上传成功")
}

func FileUploadFailed(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "文件上传失败："+e.Error())
}

func ProjectScanComplete(ctx context.Context) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgInfo, "项目扫描完成")
}

func ProjectScanning(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.UpdateStatus(ui.StatusRunning, "正在进行扫描...")
	var once sync.Once
	return func() { once.Do(func() { u.ClearStatus() }) }
}

func HashingFileFailed(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgInfo, "文件哈希计算失败："+e.Error())
}

func CodeFileUploadingForDeep(ctx context.Context) func() {
	var u = ui.Use(ctx)
	u.Display(ui.MsgInfo, "正在上传代码进行深度检测")
	u.UpdateStatus(ui.StatusRunning, "代码上传中")
	return func() { u.ClearStatus() }
}

func CodeFileUploadErr(ctx context.Context, e error) {
	var u = ui.Use(ctx)
	u.Display(ui.MsgError, "代码上传失败："+e.Error())
}

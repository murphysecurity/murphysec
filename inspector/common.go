package inspector

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Logger = zap.NewNop()

func createTaskC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")
	defer ui.ClearStatus()
	e = createTaskApi(ctx)
	if errors.Is(e, api.ErrTokenInvalid) {
		ui.Display(display.MsgError, "任务创建失败，Token 无效")
	} else if e != nil {
		ui.Display(display.MsgError, fmt.Sprintf("任务创建失败：%s", e.Error()))
	}
	return
}

func submitModuleInfoC(ctx context.Context) (e error) {
	ui := model.UseScanTask(ctx).UI()
	ui.UpdateStatus(display.StatusRunning, "项目扫描结束，正在提交信息...")
	defer ui.ClearStatus()
	e = submitModuleInfoApi(ctx)
	if e != nil {
		ui.Display(display.MsgError, fmt.Sprint("信息提交失败：", e.Error()))
	}
	return
}

func startCheckC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.UpdateStatus(display.StatusRunning, "正在启动检测")
	defer ui.ClearStatus()
	e = api.StartCheckTaskType(scanTask.TaskId, scanTask.Kind)
	if e != nil {
		ui.Display(display.MsgError, "启动检测失败："+e.Error())
	}
	return
}

func queryResultC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.UpdateStatus(display.StatusRunning, "正在等待服务器返回结果")
	defer ui.ClearStatus()
	var resp *model.TaskScanResponse
	resp, e = api.QueryResult(scanTask.TaskId)
	if e != nil {
		ui.Display(display.MsgError, fmt.Sprintf("获取扫描结果失败：%s", e.Error()))
	} else {
		scanTask.ScanResult = resp
	}
	return
}

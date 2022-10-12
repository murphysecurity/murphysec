package inspector

import (
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/view"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Logger = zap.NewNop()

func createTaskC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	defer view.TaskCreating(ui)()

	e = createTaskApi(ctx)
	if errors.Is(e, api.ErrTlsRequest) {
		view.TLSAlert(ui, e)
		return
	}
	if errors.Is(e, api.ErrTokenInvalid) {
		view.TokenInvalid(ui)
		return
	}
	if e != nil {
		view.TaskCreateFailed(ui, e)
	}
	return
}

func submitModuleInfoC(ctx context.Context) (e error) {
	ui := model.UseScanTask(ctx).UI()
	defer view.ScanCompleteSubmitting(ui)()
	e = submitModuleInfoApi(ctx)
	if e != nil {
		view.SubmitError(ui, e)
	}
	return
}

func startCheckC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	defer view.StartingInspection(ui)()
	e = api.StartCheckTaskType(scanTask.TaskId, scanTask.Kind)
	if e != nil {
		view.StartingInspectionFailed(ui, e)
	}
	return
}

func queryResultC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	view.WaitingServerResponse(ui)()
	var resp *model.TaskScanResponse
	resp, e = api.QueryResult(scanTask.TaskId)
	if e != nil {
		view.GetScanResultFailed(ui, e)
	} else {
		scanTask.ScanResult = resp
	}
	return
}

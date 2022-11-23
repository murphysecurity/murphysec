package inspector

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/view"
	"go.uber.org/zap"
)

func Scan(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	view.ProjectName(ctx, scanTask.SubtaskName())
	if e := createTaskC(ctx); e != nil {
		return e
	}

	var scanComplete = view.ProjectScanning(ctx)

	if e := managedInspect(ctx); e != nil {
		Logger.Error("Managed inspect failed", zap.Error(e))
		return e
	}

	scanComplete()
	view.ProjectScanComplete(ctx)

	if err := submitModuleInfoC(ctx); err != nil {
		return err
	}

	if err := startCheckC(ctx); err != nil {
		return err
	}
	if e := queryResultC(ctx); e != nil {
		return e
	}
	view.DisplayScanResultSummary(ctx, scanTask.Result().RelyNum, scanTask.Result().LeakNum)
	return nil
}

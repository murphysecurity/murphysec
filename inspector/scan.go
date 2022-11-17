package inspector

import (
	"context"
	"github.com/murphysecurity/murphysec/build_flags"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/view"
	"go.uber.org/zap"
)

func Scan(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	view.ProjectName(ui, scanTask.ProjectName)
	if e := createTaskC(ctx); e != nil {
		return e
	}

	var scanComplete = view.ProjectScanning(ui)

	if e := managedInspect(ctx); e != nil {
		Logger.Error("Managed inspect failed", zap.Error(e))
		return e
	}

	if build_flags.AllowFileHash {
		Logger.Info("File hash scanning")
		if e := FileHashScan(ctx); e != nil {
			Logger.Error("FileHash calc failed", zap.Error(e))
			view.HashingFileFailed(ui, e)
		}
	}

	if build_flags.AllowDeepScan && scanTask.EnableDeepScan {
		Logger.Info("DeepScan")
		scanComplete = view.CodeFileUploadingForDeep(ui)
		if e := UploadCodeFile(ctx); e != nil {
			Logger.Error("Code upload failed", zap.Error(e))
			view.CodeFileUploadErr(ui, e)
		}
	}

	scanComplete()
	view.ProjectScanComplete(ui)

	if err := submitModuleInfoC(ctx); err != nil {
		return err
	}

	if err := startCheckC(ctx); err != nil {
		return err
	}
	if e := queryResultC(ctx); e != nil {
		return e
	}
	view.DisplayScanResultSummary(ui, scanTask.ScanResult.DependenciesCount, scanTask.ScanResult.IssuesCompsCount)
	view.DisplayScanResultReport(ui, scanTask.ScanResult.ReportURL())

	return nil
}

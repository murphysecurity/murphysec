package model

import (
	"context"
	"path/filepath"
)

type InspectorTask struct {
	*ScanTask
	ScanDir string
}

func (i *InspectorTask) AddModule(module Module) {
	module.FilePath = filepath.ToSlash(module.FilePath)
	module.PackageFile = filepath.ToSlash(module.PackageFile)
	i.Modules = append(i.Modules, module)
}

func WithInspectorTask(ctx context.Context, scanDir string) context.Context {
	p := UseScanTask(ctx)
	if p == nil {
		panic("scan task not exists")
	}
	task := &InspectorTask{
		ScanTask: p,
		ScanDir:  scanDir,
	}
	return context.WithValue(ctx, inspectorTaskKey, task)
}

func UseInspectorTask(ctx context.Context) *InspectorTask {
	p, ok := ctx.Value(inspectorTaskKey).(*InspectorTask)
	if ok {
		return p
	}
	return nil
}

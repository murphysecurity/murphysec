package model

import (
	"context"
	"murphysec-cli-simple/logger"
	"path/filepath"
)

type InspectorTask struct {
	*ScanTask
	ScanDir string
}

func (i *InspectorTask) AddModule(module Module) {
	if filepath.IsAbs(module.FilePath) {
		f, e := filepath.Rel(i.ProjectDir, module.FilePath)
		if e != nil {
			logger.Warn.Println("Calc rel-path failed.", e.Error(), module.FilePath)
		} else {
			module.FilePath = f
		}
	}
	module.FilePath = filepath.ToSlash(module.FilePath)
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

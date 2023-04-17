package model

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"path/filepath"
)

type InspectionTask struct {
	ctx           context.Context
	scanTask      *ScanTask
	inspectionDir string
}

// Dir 返回当前扫描器扫描的路径，绝对路径
func (i *InspectionTask) Dir() string {
	return i.inspectionDir
}

// RelDir 返回相对路径
func (i *InspectionTask) RelDir() string {
	rel, e := filepath.Rel(i.scanTask.ProjectPath, i.inspectionDir)
	if e != nil {
		return ""
	}
	return rel
}

type _InspectionTaskCtxKey struct{}

var inspectionTaskCtxKey _InspectionTaskCtxKey

func WithInspectionTask(ctx context.Context, task *InspectionTask) context.Context {
	if task == nil {
		panic("task == nil")
	}
	return context.WithValue(ctx, inspectionTaskCtxKey, task)
}

func UseInspectionTask(ctx context.Context) *InspectionTask {
	i, _ := ctx.Value(inspectionTaskCtxKey).(*InspectionTask)
	return i
}

func (i *InspectionTask) AddModule(module Module) {
	var logger = logctx.Use(i.ctx).Sugar()
	logger.Infof("add module: %v", module)
	if filepath.IsAbs(module.ModulePath) {
		relPath, e := filepath.Rel(i.scanTask.ProjectPath, module.ModulePath)
		if e != nil {
			logger.Warnf("get module relative-path: %v", e)
		}
		module.ModulePath = relPath
	}
	if module.ModulePath == "." {
		module.ModulePath = "./"
	}
	module.ModulePath = filepath.ToSlash(module.ModulePath)
	i.scanTask.Modules = append(i.scanTask.Modules, module)
}

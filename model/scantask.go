package model

import (
	"context"
	"github.com/murphysecurity/murphysec/gitinfo"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"os"
	"path/filepath"
)

type ScanTaskBuilder struct {
	Ctx         context.Context
	Mode        ScanMode
	AccessType  AccessType
	Path        string
	TaskId      string
	SubtaskName string
}

func (s ScanTaskBuilder) Build() (*ScanTask, error) {
	if s.Ctx == nil {
		panic("context is nil")
	}
	var logger = logctx.Use(s.Ctx).Sugar()

	if !filepath.IsAbs(s.Path) {
		return nil, ErrPathIsNotAbsolute
	}
	stat, e := os.Stat(s.Path)
	if e != nil {
		return nil, e
	}
	var isDir = stat.IsDir()
	if s.Mode == ScanModeSource && isDir {
		return nil, ErrMustBeDirectory
	}
	var task ScanTask
	task.ctx = s.Ctx
	task.projectPath = s.Path
	task.mode = s.Mode
	task.accessType = s.AccessType
	task.taskId = s.TaskId
	task.subtaskName = s.SubtaskName
	if s.SubtaskName == "" {
		return nil, ErrSubtaskNameIsEmpty
	}
	if isDir {
		summary, e := gitinfo.GetSummary(s.Ctx, task.projectPath)
		if e != nil {
			logger.Errorf("Read git info failed, %v", e)
		} else {
			task.gitInfo = summary
		}
	}
	return &task, nil
}

// ScanTask 表示当前扫描任务
type ScanTask struct {
	ctx         context.Context
	projectPath string
	accessType  AccessType
	mode        ScanMode
	gitInfo     *gitinfo.Summary
	taskId      string
	subtaskId   string
	subtaskName string
	modules     []Module
	result      *ScanResultResponse
}

func (s *ScanTask) AccessType() AccessType {
	return s.accessType
}

func (s *ScanTask) ScanMode() ScanMode {
	return s.mode
}

func (s *ScanTask) SubtaskName() string {
	return s.subtaskName
}

func (s *ScanTask) Result() *ScanResultResponse {
	return s.result
}

func (s *ScanTask) ProjectPath() string {
	return s.projectPath
}
func (s *ScanTask) BuildInspectionTask(dir string) *InspectionTask {
	return &InspectionTask{
		ctx:           s.ctx,
		scanTask:      s,
		inspectionDir: dir,
	}
}

type _scanTaskType struct{}

var scanTaskType _scanTaskType

func WithScanTask(ctx context.Context, task *ScanTask) context.Context {
	return context.WithValue(ctx, scanTaskType, task)
}

func UseScanTask(ctx context.Context) *ScanTask {
	return ctx.Value(scanTaskType).(*ScanTask)
}

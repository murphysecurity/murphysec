package model

import (
	"context"
)

// ScanTask 表示当前扫描任务
type ScanTask struct {
	Ctx           context.Context
	ProjectPath   string
	AccessType    AccessType
	Mode          ScanMode
	TaskId        string
	SubtaskId     string
	Modules       []Module
	CodeFragments []ComponentCodeFragment
	Result        *ScanResultResponse
	SubtaskName   string
}

func (s *ScanTask) BuildInspectionTask(dir string) *InspectionTask {
	return &InspectionTask{
		ctx:           s.Ctx,
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

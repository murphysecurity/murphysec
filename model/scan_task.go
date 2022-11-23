package model

import (
	"context"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

type key int

const (
	scanTaskKey key = iota + 1
	inspectorTaskKey
)

type TaskKind string

const (
	TaskKindNormal     TaskKind = "Normal"
	TaskKindBinary     TaskKind = "Binary"
	TaskKindIotScan    TaskKind = "IotScan"
	TaskKindDockerfile TaskKind = "Dockerfile"
	TaskKindHostEnv    TaskKind = "HostEnvironment"
)

type ProjectType string

const (
	ProjectTypeLocal ProjectType = "Local"
	ProjectTypeGit   ProjectType = "Git"
)

type FileHash struct {
	Hash []string `json:"hash"`
	Path string   `json:"path"`
}

type ScanTask struct {
	TaskId            string
	ProjectDir        string
	ProjectName       string
	Kind              TaskKind
	ProjectType       ProjectType
	ProjectId         string
	Username          string
	StartTime         time.Time
	GitInfo           *GitInfo
	TaskType          TaskType
	TotalContributors int
	Modules           []Module
	ScanResult        *TaskScanResponse
}

func CreateScanTask(projectDir string, taskKind TaskKind, taskType TaskType) *ScanTask {
	must.True(projectDir == "" || filepath.IsAbs(projectDir))
	t := &ScanTask{
		ProjectDir:  projectDir,
		ProjectName: filepath.Base(projectDir),
		Kind:        taskKind,
		ProjectType: ProjectTypeLocal,
		ProjectId:   "",
		StartTime:   time.Now(),
		GitInfo:     nil,
		TaskType:    taskType,
	}
	fillScanTaskGitInfo(t)
	return t
}

func fillScanTaskGitInfo(task *ScanTask) {
	if env.DisableGit {
		Logger.Debug("Git info is disabled")
		return
	}
	Logger.Debug("Check git repo", zap.String("dir", task.ProjectDir))
	gitInfo, e := getGitInfo(task.ProjectDir)
	if e != nil {
		Logger.Warn("Read git info failed", zap.Error(e))
		return
	}
	task.GitInfo = gitInfo
	task.ProjectName = gitInfo.ProjectName
	task.ProjectType = ProjectTypeGit
}

func WithScanTask(ctx context.Context, task *ScanTask) context.Context {
	return context.WithValue(ctx, scanTaskKey, task)
}

func UseScanTask(ctx context.Context) *ScanTask {
	t, ok := ctx.Value(scanTaskKey).(*ScanTask)
	if ok {
		return t
	}
	return nil
}

func (s *ScanTask) UI() display.UI {
	return s.TaskType.UI()
}

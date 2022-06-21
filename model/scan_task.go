package model

import (
	"context"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
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
	TaskKindNormal  TaskKind = "Normal"
	TaskKindBinary  TaskKind = "Binary"
	TaskKindIotScan TaskKind = "IotScan"
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
	ContributorList   []Contributor
	TotalContributors int
	Modules           []Module
	ScanResult        *TaskScanResponse
	EnableDeepScan    bool
	FileHashes        []FileHash
}

func CreateScanTask(projectDir string, taskKind TaskKind, taskType TaskType) *ScanTask {
	must.True(filepath.IsAbs(projectDir))
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
		return
	}
	gitInfo, e := getGitInfo(task.ProjectDir)
	if errors.Is(e, ErrNoGitRepo) {
		logger.Info.Println("No git repo information.", e.Error())
		return
	}
	if e != nil {
		logger.Warn.Println("Read git info failed.", e.Error())
		return
	}
	task.GitInfo = gitInfo
	task.ProjectName = gitInfo.ProjectName
	task.ProjectType = ProjectTypeGit
	contributors, e := collectContributor(task.ProjectDir)
	if e != nil {
		logger.Warn.Println("Collect contributors failed.", e.Error())
		return
	}
	task.ContributorList = contributors
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

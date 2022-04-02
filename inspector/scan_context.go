package inspector

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	base2 "murphysec-cli-simple/base"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"os"
	"path/filepath"
	"time"
)

var ErrProjectDirInvalid = errors.New("Project dir invalid.")
var ErrGetProjectInfo = errors.New("Get project info failed.")

type ScanContext struct {
	GitInfo           *GitInfo
	ProjectName       string
	ProjectDir        string
	ManagedModules    []base.Module
	StartTime         time.Time
	TaskId            string
	FileHashes        []string
	ProjectType       string
	EnableDeepScan    bool
	ScanResult        *api.TaskScanResponse
	TaskType          base2.InspectTaskType
	InspectorError    []base.InspectorError
	ContributorList   []api.Contributor
	ProjectId         string
	TotalContributors int
}

func (s *ScanContext) UI() display.UI {
	return s.TaskType.UI()
}

func NewTaskContext(dir string, taskType base2.InspectTaskType) (*ScanContext, error) {
	ctx := &ScanContext{
		TaskType:  taskType,
		StartTime: time.Now(),
	}
	if baseDir, e := filepath.Abs(dir); e != nil {
		return nil, ErrProjectDirInvalid
	} else {
		ctx.ProjectDir = baseDir
	}
	{
		info, e := os.Stat(ctx.ProjectDir)
		if e != nil || info == nil || !info.IsDir() {
			return nil, ErrProjectDirInvalid
		}
	}
	return ctx, nil
}

func (s *ScanContext) AddManagedModule(module base.Module) {
	module.UUID = uuid.New()
	s.ManagedModules = append(s.ManagedModules, module)
}

func (s *ScanContext) FillProjectInfo() error {
	if s.ProjectDir == "" {
		panic("project dir is empty")
	}
	gitInfo, e := getGitInfo(s.ProjectDir)
	if e != nil {
		logger.Warn.Println("Get git info failed", e.Error())
	}
	if gitInfo == nil {
		s.ProjectType = "Local"
		logger.Info.Println("git not detected, fallback")
	} else {
		s.ProjectType = "Git"
		s.GitInfo = gitInfo
		s.ProjectName = gitInfo.ProjectName
	}
	if s.ProjectName == "" {
		logger.Info.Println("get project name failed, use directory name")
		s.ProjectName = filepath.Base(s.ProjectDir)
	}
	if s.ProjectName == "" {
		logger.Warn.Println("Get project name failed")
		s.ProjectName = "<NoTitle>"
	} else {
		logger.Info.Println("Project name:", s.ProjectName)
	}
	return nil
}

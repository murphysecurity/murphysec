package inspector

import (
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"path/filepath"
	"time"
)

// TaskInfo 外部传入标记
var TaskInfo string

var ErrNoEngineMatched = errors.New("ErrNoEngineMatched")
var ErrAPITokenInvalid = errors.New("ErrAPITokenInvalid")
var ErrNoModule = errors.New("ErrNoModule")

type ManagedScanContext struct {
	GitInfo        *GitInfo
	ProjectName    string
	TaskSource     api.InspectTaskSource
	ProjectDir     string
	ManagedModules []base.Module
	StartTime      time.Time
}

func (m *ManagedScanContext) WrapProjectInfo(projectDir string) {
	m.ProjectDir = projectDir
	gitInfo, e := getGitInfo(projectDir)
	if e != nil {
		logger.Err.Printf("Get git info failed: %+v\n", e)
	}
	if gitInfo == nil {
		logger.Info.Println("No valid git info found.")
	} else {
		m.GitInfo = gitInfo
		m.ProjectName = gitInfo.ProjectName
	}
	if m.ProjectName == "" {
		m.ProjectName = filepath.Base(projectDir)
	}
	logger.Info.Println("Project name:", m.ProjectName)
	if m.ProjectName == "" {
		logger.Warn.Println("Resolve project name failed.")
	}
}

func mapVoGitInfoOrNil(g *GitInfo) *api.VoGitInfo {
	if g == nil {
		return nil
	}
	return &api.VoGitInfo{
		Commit:       g.HeadCommitHash,
		GitRef:       g.HeadRefName,
		GitRemoteUrl: g.RemoteURL,
	}
}

package inspector

import (
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TaskInfo 外部传入标记
var TaskInfo string

var ErrNoEngineMatched = errors.New("ErrNoEngineMatched")
var ErrAPITokenInvalid = errors.New("ErrAPITokenInvalid")
var ErrNoModule = errors.New("ErrNoModule")

type ScanContext struct {
	GitInfo        *GitInfo
	ProjectName    string
	TaskSource     api.InspectTaskSource
	ProjectDir     string
	ManagedModules []base.Module
	StartTime      time.Time
}

func (ctx *ScanContext) WrapProjectInfo(projectDir string) {
	ctx.ProjectDir = projectDir
	gitInfo, e := getGitInfo(projectDir)
	if e != nil {
		logger.Err.Printf("Get git info failed: %+v\n", e)
	}
	if gitInfo == nil {
		logger.Info.Println("No valid git info found.")
	} else {
		ctx.GitInfo = gitInfo
		ctx.ProjectName = gitInfo.ProjectName
	}
	if ctx.ProjectName == "" {
		ctx.ProjectName = filepath.Base(projectDir)
	}
	logger.Info.Println("Project name:", ctx.ProjectName)
	if ctx.ProjectName == "" {
		logger.Warn.Println("Resolve project name failed.")
	}
}

func (ctx *ScanContext) getApiRequestObj() *api.UserCliDetectInput {
	// api request object
	req := &api.UserCliDetectInput{
		ApiToken:           conf.APIToken(),
		CliVersion:         version.Version(),
		CmdLine:            strings.Join(os.Args, " "),
		TargetAbsPath:      ctx.ProjectDir,
		TaskConsumeTime:    int(time.Now().Sub(ctx.StartTime).Seconds()),
		TaskInfo:           TaskInfo,
		TaskStartTimestamp: ctx.StartTime.Unix(),
		TaskSource:         ctx.TaskSource,
		UserAgent:          version.UserAgent(),
		ProjectName:        ctx.ProjectName,
	}
	return req
}

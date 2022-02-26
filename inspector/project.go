package inspector

import (
	"murphysec-cli-simple/logger"
	"path/filepath"
)

func readProjectInfo(baseDir string) *ScanContext {
	ctx := new(ScanContext)
	ctx.ProjectDir = baseDir
	gitInfo, e := getGitInfo(ctx.ProjectDir)
	if e != nil {
		logger.Warn.Println("Get git info failed", e.Error())
	}
	if gitInfo == nil {
		logger.Info.Println("git not detected, fallback")
	} else {
		ctx.GitInfo = gitInfo
		ctx.ProjectName = gitInfo.ProjectName
	}
	if ctx.ProjectName == "" {
		logger.Info.Println("get project name failed, use directory name")
		ctx.ProjectName = filepath.Base(ctx.ProjectDir)
	}
	if ctx.ProjectName == "" {
		logger.Warn.Println("Get project name failed")
		ctx.ProjectName = "<NoTitle>"
	} else {
		logger.Info.Println("Project name:", ctx.ProjectName)
	}
	return ctx
}

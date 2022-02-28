package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"time"
)

func ManagedInspect(ctx *ScanContext) (*api.VoDetectResponse, error) {
	logger.Info.Println("Start managed inspect...", ctx.ProjectDir)
	// 包管理器的扫描
	if e := managedInspectScan(ctx); e != nil {
		logger.Info.Printf("Managed inspect failed, %v\n", e)
		return nil, e
	}
	return managedInspectAPIRequest(ctx)
}

func IdeaScan(dir string) (interface{}, error) {
	if info, e := os.Stat(dir); e != nil || !info.IsDir() {
		reportIdeaStatus(1, "directory invalid")
		return nil, errors.Wrap(e, "invalid directory")
	}
	ctx := &ScanContext{TaskSource: api.TaskSourceIdea, StartTime: time.Now()}
	ctx.WrapProjectInfo(dir)
	response, e := ManagedInspect(ctx)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %v\n", e)
		if e == api.ErrTokenInvalid {
			reportIdeaStatus(4, "Token invalid")
			return nil, e
		}
		return nil, e
	}
	if e == ErrNoEngineMatched {
		reportIdeaStatus(IdeaNoEngineMatch, "Engine not match")
		return nil, e
	}
	if e == nil {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
		return nil, nil
	}
	reportIdeaStatus(IdeaUnknownErr, e.Error())
	return nil, e
}
func CliScan(dir string, jsonOutput bool) (interface{}, error) {
	if info, e := os.Stat(dir); e != nil || !info.IsDir() {
		fmt.Println("给定的路径无效")
		return nil, errors.Wrap(e, "invalid directory")
	}
	ctx := &ScanContext{TaskSource: api.TaskSourceCli, StartTime: time.Now()}
	// detect Jenkins environment
	if os.Getenv("JENKINS_HOME") != "" && os.Getenv("JENKINS_URL") != "" {
		ctx.TaskSource = api.TaskSourceJenkins
	}
	ctx.WrapProjectInfo(dir)
	response, e := ManagedInspect(ctx)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %v\n", e)
		if e == api.ErrTokenInvalid {
			fmt.Println("Token 无效")
			return nil, e
		}
		return nil, e
	}
	if e == nil {
		if jsonOutput {
			fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
		} else {
			fmt.Println(fmt.Sprintf("扫描完成，共计%d个组件，%d个漏洞", response.DependenciesCount, response.IssuesCompsCount))
		}
		return nil, nil
	}
	if e == ErrNoEngineMatched {
		if jsonOutput {
			reportIdeaStatus(IdeaNoEngineMatch, e.Error())
		} else {
			fmt.Println("没有找到受支持的模块")
		}
	}
	return nil, nil
}

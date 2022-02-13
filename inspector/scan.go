package inspector

import (
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"time"
)

func ManagedInspect(ctx *ScanContext) (*api.VoDetectResponse, error) {
	logger.Info.Println("Start managed inspect...", ctx.ProjectDir)
	// 包管理器的扫描
	if e := managedInspectScan(ctx); e != nil {
		logger.Info.Printf("Managed inspect failed, %+v\n", e)
		return nil, e
	}
	return managedInspectAPIRequest(ctx)
}

func IdeaScan(dir string) (interface{}, error) {
	ctx := &ScanContext{TaskSource: api.TaskSourceIdea, StartTime: time.Now()}
	ctx.WrapProjectInfo(dir)
	response, e := ManagedInspect(ctx)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %+v\n", e)
		if e == api.ErrTokenInvalid {
			reportIdeaStatus(4, "Token invalid")
			return nil, e
		}
		return nil, e
	}
	if e == nil {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
		return nil, nil
	}
	// 文件哈希扫描
	response, e = FileHashInspect(ctx)
	if e != nil {
		logger.Err.Printf("FileHash scan failed: %+v\n", e)
		if e == api.ErrTokenInvalid {
			reportIdeaStatus(4, "Token invalid")
			return nil, e
		}
		return nil, e
	}
	// 扫描成功
	fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
	return nil, nil
}
func CliScan(dir string, jsonOutput bool) (interface{}, error) {
	ctx := &ScanContext{TaskSource: api.TaskSourceCli, StartTime: time.Now()}
	ctx.WrapProjectInfo(dir)
	response, e := ManagedInspect(ctx)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %+v\n", e)
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
	// 文件哈希扫描
	response, e = FileHashInspect(ctx)
	if e != nil {
		logger.Err.Printf("FileHash scan failed: %+v\n", e)
		if e == api.ErrTokenInvalid {
			fmt.Println("Token 无效")
			return nil, e
		}
		return nil, e
	}
	// 扫描成功
	if jsonOutput {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
	} else {
		fmt.Println(fmt.Sprintf("扫描完成，共计%d个组件，%d个漏洞", response.DependenciesCount, response.IssuesCompsCount))
	}
	return nil, nil
}

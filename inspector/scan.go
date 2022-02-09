package inspector

import (
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"time"
)

func ManagedInspect(dir string, taskSource api.InspectTaskSource) (*api.VoDetectResponse, error) {
	logger.Info.Println("Start managed inspect...", dir)
	// 包管理器的扫描
	ctx := &ManagedScanContext{
		StartTime:  time.Now(),
		TaskSource: taskSource,
	}
	ctx.WrapProjectInfo(dir)
	if e := managedInspectScan(ctx); e != nil {
		logger.Info.Printf("Managed inspect failed, %+v\n", e)
		return nil, e
	}
	response, e := managedInspectAPIRequest(ctx)
	if e != nil {
		return nil, e
	}
	return response, nil
}

func IdeaScan(dir string) (interface{}, error) {
	response, e := ManagedInspect(dir, api.TaskSourceIdea)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %+v\n", e)
		if e == api.ErrTokenInvalid {
			reportIdeaStatus(4, "Token invalid")
			return nil, e
		}
		return nil, e
	}
	// 扫描成功
	if e == nil {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
		return nil, nil
	}
	// todo: cpp
	return nil, nil
}
func CliScan(dir string, jsonOutput bool) (interface{}, error) {
	response, e := ManagedInspect(dir, api.TaskSourceCli)
	// 扫描出错
	if e != nil && e != ErrNoEngineMatched {
		logger.Err.Printf("Managed scan failed: %+v\n", e)
		if e == api.ErrTokenInvalid {
			fmt.Println("Token 无效")
			return nil, e
		}
		return nil, e
	}
	if jsonOutput {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(response)))))
	} else {
		fmt.Println(fmt.Sprintf("扫描完成，共计%d个组件，%d个漏洞", response.DependenciesCount, response.IssuesCompsCount))
	}
	return nil, nil
}

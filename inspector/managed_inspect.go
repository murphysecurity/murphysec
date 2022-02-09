package inspector

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"strings"
	"time"
)

var managedInspector = []base.Inspector{
	go_mod.New(),
	maven.New(),
	npm.New(),
}

func managedInspectAPIRequest(ctx *ManagedScanContext) (*api.VoDetectResponse, error) {
	must.True(len(ctx.ManagedModules) > 0)
	// api request object
	req := api.UserCliDetectInput{
		ApiToken:           conf.APIToken(),
		CliVersion:         version.Version(),
		CmdLine:            strings.Join(os.Args, " "),
		Modules:            []api.VoModule{},
		TargetAbsPath:      ctx.ProjectDir,
		TaskConsumeTime:    int(time.Now().Sub(ctx.StartTime).Seconds()),
		TaskInfo:           TaskInfo,
		TaskStartTimestamp: ctx.StartTime.Unix(),
		TaskSource:         ctx.TaskSource,
		UserAgent:          version.UserAgent(),
		ProjectName:        ctx.ProjectName,
	}
	// 拼请求体
	uuidModuleMap := map[uuid.UUID]base.Module{}
	for _, it := range ctx.ManagedModules {
		_uuid := uuid.Must(uuid.NewRandom())
		uuidModuleMap[_uuid] = it
		voM := it.ApiVo()
		voM.ModuleUUID = _uuid
		req.Modules = append(req.Modules, *voM)
	}
	response, e := api.SendDetect(req)
	if e == api.ErrTokenInvalid {
		return nil, ErrAPITokenInvalid
	}
	if e != nil {
		return nil, errors.Wrap(e, "API request failed")
	}
	return response, nil
}

// 受管理扫描
func managedInspectScan(ctx *ManagedScanContext) error {
	dir := ctx.ProjectDir
	startTime := time.Now()
	logger.Info.Println("Auto scan dir:", dir)
	var inspectors []base.Inspector
	{
		// 尝试匹配检测器
		logger.Debug.Println("Try match managed inspector...")
		for _, it := range managedInspector {
			if it.CheckDir(dir) {
				inspectors = append(inspectors, it)
			}
		}
		logger.Debug.Println("Matched managed inspector:", inspectors)
	}
	if len(inspectors) == 0 {
		logger.Debug.Println("No managed inspector matched")
		return ErrNoEngineMatched
	}

	for _, it := range inspectors {
		rs, e := it.Inspect(dir)
		if e != nil {
			logger.Err.Printf("Engine: %v scan failed. Reason: %+v\n", it, e)
			continue
		}
		fmt.Printf("Inspector terminated %v, total module: %v\n", it, len(rs))
		for _, it := range rs {
			ctx.ManagedModules = append(ctx.ManagedModules, it)
		}
	}
	endTime := time.Now()
	logger.Info.Println("Scan terminated. Cost time:", endTime.Sub(startTime))
	if len(ctx.ManagedModules) < 1 {
		return ErrNoModule
	}
	return nil
}

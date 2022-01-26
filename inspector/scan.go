package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/utils/must"
	"os"
	"time"
)

func CliScan(dir string, jsonOutput bool) (interface{}, error) {
	startTime := time.Now()
	engine := tryMatchInspector(dir)
	if engine == nil {
		return nil, errors.New("Can't inspect project. No inspector supported.")
	}
	// 开始扫描
	logger.Info.Println("IdeaScan dir:", dir, "Inspector:", engine.String(), "Time:", startTime.Format(time.RFC3339))
	modules, e := engine.Inspect(dir)
	if e != nil {
		return nil, errors.Wrap(e, "Engine scan failed.")
	}
	taskType := "client"
	if os.Getenv("CI") != "" {
		taskType = "ci"
	}
	req := getAPIRequest(taskType)
	// 拼凑项目信息
	wrapProjectInfoToReqObj(req, dir)
	logger.Debug.Println("Before scan. projectName:", req.ProjectName, "git:", req.GitInfo != nil)
	// 拼凑请求体 模块
	moduleUUIDMap := map[uuid.UUID]base.Module{}
	for _, it := range modules {
		moduleVo := mapVoModule(it)
		moduleVo.ModuleUUID = uuid.Must(uuid.NewRandom())
		moduleUUIDMap[moduleVo.ModuleUUID] = it
		req.Modules = append(req.Modules, moduleVo)
	}
	// API 请求
	r, e := api.SendDetect(*req)
	if e != nil {
		return nil, errors.Wrap(e, "Server request failed.")
	}
	// 输出 API 响应
	if jsonOutput {
		fmt.Println(string(must.Byte(json.Marshal(mapForIdea(r)))))
	} else {
		fmt.Println(fmt.Sprintf("扫描完成，共计%d个组件，%d个漏洞", r.DependenciesCount, r.IssuesCompsCount))
	}
	if _, ok := engine.(*maven.Inspector); ok {
		javaImportClauseScan(r, dir)
	}
	return nil, nil
}

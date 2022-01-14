package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IdeaScan(dir string, pmType base.PackageManagerType) (interface{}, error) {
	startTime := time.Now()
	engine := getInspectorSupportPkgManagerType(pmType)
	// 检查项目是否可以被指定的扫描器扫描
	if !engine.CheckDir(dir) {
		logger.Err.Println("Dir can't be scan.", dir)
		ideaFail(3, "Dir can't be scan.")
		return nil, errors.New("Can't be scan")
	}
	// 开始扫描
	logger.Info.Println("IdeaScan dir:", dir, "PackageManagerType:", pmType)
	logger.Info.Println("Task start at:", startTime.Format(time.RFC3339))
	modules, e := engine.Inspect(dir)
	if e != nil {
		ideaFail(1, "Engine scan failed.")
		return nil, e
	}
	req := getIdeaRequest()
	// 拼凑项目信息
	wrapProjectInfoToReqObj(req, dir)
	logger.Debug.Println("Before scan. projectName:", req.ProjectName, "git:", req.GitInfo != nil, "packageManager:", pmType)
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
		ideaFail(2, "Server request failed.")
		return nil, e
	}
	// 输出 API 响应
	fmt.Println(string(must.Byte(json.Marshal(mapForIdea(r)))))
	return nil, nil
}

func wrapProjectInfoToReqObj(input *api.UserCliDetectInput, projectDir string) {
	// 获取git信息
	gitInfo, e := getGitInfo(projectDir)
	if e != nil {
		logger.Err.Println("GetGitInfo failed.", e.Error())
	}
	if gitInfo == nil {
		logger.Info.Println("No git repo found")
	}
	input.GitInfo = mapVoGitInfoOrNil(gitInfo)
	input.TargetAbsPath = projectDir
	if input.GitInfo != nil && input.GitInfo.GitRemoteUrl != "" {
		u, e := giturls.Parse(input.GitInfo.GitRemoteUrl)
		if e != nil {
			input.ProjectName = u.Path
		}
	}
	if input.ProjectName == "" {
		input.ProjectName = filepath.Base(projectDir)
	}
}

func getIdeaRequest() *api.UserCliDetectInput {
	return &api.UserCliDetectInput{
		ApiToken:           conf.APIToken(),
		CliVersion:         version.Version(),
		CmdLine:            strings.Join(os.Args, " "),
		GitInfo:            nil,
		Engine:             "",
		Modules:            []api.VoModule{},
		TargetAbsPath:      "",
		TaskStartTimestamp: 0,
		TaskType:           "Plugin",
		UserAgent:          version.UserAgent(),
		TaskInfo:           os.Getenv("TASK_INFO"),
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

func mapVoDependency(d []base.Dependency) []api.VoDependency {
	r := make([]api.VoDependency, 0)
	for _, it := range d {
		r = append(r, api.VoDependency{
			Name:         it.Name,
			Version:      it.Version,
			Dependencies: mapVoDependency(it.Dependencies),
		})
	}
	return r
}
func mapVoModule(m base.Module) api.VoModule {
	r := api.VoModule{
		Dependencies:   mapVoDependency(m.Dependencies),
		Language:       m.Language,
		Name:           m.Name,
		PackageFile:    m.PackageFile,
		PackageManager: m.PackageManager,
		RelativePath:   m.RelativePath,
		RuntimeInfo:    m.RuntimeInfo,
		Version:        m.Version,
	}
	return r
}

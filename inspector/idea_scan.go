package inspector

import (
	"encoding/json"
	"fmt"
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

func wrapProjectNameToDetectObj(input *api.UserCliDetectInput, projectDir string) {
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

func IdeaScan(dir string, pmType base.PackageManagerType) (interface{}, error) {
	startTime := time.Now()
	engine := getInspectorSupportPkgManagerType(pmType)
	logger.Info.Println("IdeaScan dir:", dir, "PackageManagerType:", pmType)
	logger.Info.Println("Task start at:", startTime.Format(time.RFC3339))
	gitInfo, e := getGitInfo(dir)
	if e != nil {
		logger.Err.Println("GetGitInfo failed.", e.Error())
	}
	if gitInfo == nil {
		logger.Info.Println("No git repo found")
	}
	req := api.UserCliDetectInput{
		ApiToken:           conf.APIToken(),
		CliVersion:         version.Version(),
		CmdLine:            strings.Join(os.Args, " "),
		GitInfo:            mapVoGitInfoOrNil(gitInfo),
		Engine:             engine.Version(),
		Modules:            []api.VoModule{},
		TargetAbsPath:      dir,
		TaskStartTimestamp: int(startTime.Unix()),
		TaskType:           "Plugin",
		UserAgent:          version.UserAgent(),
	}
	wrapProjectNameToDetectObj(&req, dir)
	logger.Debug.Println("Before scan. projectName:", req.ProjectName, "git:", gitInfo != nil, "packageManager:", pmType)
	if !engine.CheckDir(dir) {
		logger.Err.Println("Dir can't be scan.", dir)
	}
	modules, e := engine.Inspect(dir)
	if e != nil {
		ideaFail(1, "Engine scan failed.")
		return nil, e
	}
	for _, it := range modules {
		req.Modules = append(req.Modules, mapVoModule(it))
	}
	r, e := api.SendDetect(req)
	if e != nil {
		ideaFail(2, "Server request failed.")
		return nil, e
	}
	fmt.Println(string(must.Byte(json.Marshal(mapForIdea(r)))))
	return nil, nil
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

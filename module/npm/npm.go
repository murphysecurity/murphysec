package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "NpmInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "package.json")) &&
		utils.IsFile(filepath.Join(dir, "package-lock.json"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	m, e := ScanNpmProject(model.UseInspectorTask(ctx).ScanDir)
	if e != nil {
		return e
	}
	for _, it := range m {
		model.UseInspectorTask(ctx).AddModule(it)
	}
	return nil
}

func ScanNpmProject(dir string) ([]model.Module, error) {
	logger.Info.Println("Scan dir, npm.", dir)
	pkgFile := filepath.Join(dir, "package-lock.json")
	logger.Debug.Println("Read package-lock file:", pkgFile)
	data, e := ioutil.ReadFile(pkgFile)
	if e != nil {
		return nil, e
	}
	var lockfile NpmPkgLock
	if e := json.Unmarshal(data, &lockfile); e != nil {
		return nil, e
	}
	logger.Debug.Println("lockfileVersion:", lockfile.LockfileVersion)
	if lockfile.LockfileVersion > 2 {
		return nil, errors.New(fmt.Sprintf("unsupported lockfileVersion: %d", lockfile.LockfileVersion))
	}
	for s := range lockfile.Dependencies {
		if strings.HasPrefix(s, "node_modules/") {
			delete(lockfile.Dependencies, s)
		}
	}
	var rootComp []string
	{
		// kahn
		indegree := map[string]int{}
		for s := range lockfile.Dependencies {
			indegree[s] = 0
		}
		for _, it := range lockfile.Dependencies {
			for d := range it.Requires {
				indegree[d] = indegree[d] + 1
			}
		}

		for k, i := range indegree {
			if i == 0 {
				rootComp = append(rootComp, k)
			}
		}
	}

	module := model.Module{
		PackageManager: model.PMNpm,
		Language:       model.JavaScript,
		PackageFile:    "package-lock.json",
		Name:           lockfile.Name,
		Version:        lockfile.Version,
		FilePath:       filepath.Join(dir, "package.json"),
		Dependencies:   []model.Dependency{},
		RuntimeInfo:    nil,
	}
	m := map[string]int{}
	for _, it := range rootComp {
		if d := _convDep(it, lockfile, m, 1); d != nil {
			module.Dependencies = append(module.Dependencies, *d)
		}
	}
	return []model.Module{module}, nil
}

func _convDep(root string, m NpmPkgLock, visited map[string]int, deep int) *model.Dependency {
	if deep > 5 {
		return nil
	}
	if _, ok := visited[root]; ok {
		return nil
	}
	visited[root] = deep
	defer delete(visited, root)
	d, ok := m.Dependencies[root]
	if !ok {
		return nil
	}
	r := model.Dependency{
		Name:         root,
		Version:      d.Version,
		Dependencies: nil,
	}
	for depName := range d.Requires {
		cd := _convDep(depName, m, visited, deep+1)
		if cd == nil {
			continue
		}
		r.Dependencies = append(r.Dependencies, *cd)
	}
	return &r
}

//goland:noinspection GoNameStartsWithPackageName
type NpmPkgLock struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	LockfileVersion int    `json:"LockfileVersion"`
	Dependencies    map[string]struct {
		Version  string                 `json:"version"`
		Requires map[string]interface{} `json:"requires"`
	} `json:"dependencies"`
}

package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Npm"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "package-lock.json"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	m, e := ScanNpmProject(ctx)
	if e != nil {
		return e
	}
	for _, it := range m {
		model.UseInspectionTask(ctx).AddModule(it)
	}
	return nil
}

func ScanNpmProject(ctx context.Context) ([]model.Module, error) {
	dir := model.UseInspectionTask(ctx).Dir()
	pkgFile := filepath.Join(dir, "package-lock.json")
	module := model.Module{
		PackageManager: "npm",
		ModuleName:     "",
		ModuleVersion:  "",
		ModulePath:     pkgFile,
	}
	logger := logctx.Use(ctx)
	logger.Debug("Read package-lock.json", zap.String("path", pkgFile))
	data, e := os.ReadFile(pkgFile)
	if e != nil {
		return nil, errors.WithMessage(e, "Errors when reading package-lock.json")
	}
	lockfileVer, e := parseLockfileVersion(data)
	if e != nil {
		return nil, e
	}
	if lockfileVer == 3 {
		parsed, e := parseLockfileV3(data)
		if e != nil {
			return nil, fmt.Errorf("v3lockfile: %w", e)
		}
		module.ModuleName = parsed.Name
		module.ModuleVersion = parsed.Version
		module.Dependencies = parsed.Deps
		return []model.Module{module}, nil
	}
	var lockfile NpmPkgLock
	if e := json.Unmarshal(data, &lockfile); e != nil {
		return nil, e
	}
	if lockfile.LockfileVersion > 2 || lockfile.LockfileVersion < 1 {
		return nil, errors.New(fmt.Sprintf("unsupported lockfileVersion: %d", lockfile.LockfileVersion))
	}
	module.ModuleName = lockfile.Name
	module.ModuleVersion = lockfile.Version
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
	if len(rootComp) == 0 {
		logger.Warn("Not found root component")
	}

	m := map[string]int{}
	for _, it := range rootComp {
		if d := _convDep(it, lockfile, m, 0); d != nil {
			module.Dependencies = append(module.Dependencies, *d)
		}
	}
	return []model.Module{module}, nil
}

func _convDep(root string, m NpmPkgLock, visited map[string]int, deep int) *model.DependencyItem {
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
	r := model.DependencyItem{
		Component: model.Component{
			CompName:    root,
			CompVersion: d.Version,
			EcoRepo:     EcoRepo,
		},
		IsDirectDependency: deep == 0,
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
	LockfileVersion int    `json:"lockfileVersion"`
	Dependencies    map[string]struct {
		Version  string                 `json:"version"`
		Requires map[string]interface{} `json:"requires"`
	} `json:"dependencies"`
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

func parseLockfileVersion(data []byte) (int, error) {
	type unknownVersionLockfile struct {
		LockfileVersion int `json:"lockfileVersion"`
	}
	var u unknownVersionLockfile
	if e := json.Unmarshal(data, &u); e != nil {
		return 0, fmt.Errorf("parse lockfile version failed: %w", e)
	}
	return u.LockfileVersion, nil
}

package npm

import (
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"path"
)

type v3Lockfile struct {
	Name            string               `json:"name"`
	Version         string               `json:"version"`
	LockfileVersion int                  `json:"lockfileVersion"`
	Packages        map[string]v3Package `json:"packages"`
}

type v3Package struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Dev             bool              `json:"dev"`
}

type v3ParsedLockfile struct {
	Name    string
	Version string
	Deps    []model.DependencyItem
}

func processLockfileV3(data []byte) (r *v3ParsedLockfile, e error) {
	var lockfile v3Lockfile
	if e := json.Unmarshal(data, &lockfile); e != nil {
		return nil, fmt.Errorf("parse lockfile failed: %w", e)
	}
	if lockfile.LockfileVersion != 3 {
		return nil, fmt.Errorf("unsupported lockfile version: %d", lockfile.LockfileVersion)
	}
	if lockfile.Packages == nil {
		lockfile.Packages = make(map[string]v3Package)
	}
	parsedLockfile := v3ParsedLockfile{
		Name:    lockfile.Name,
		Version: lockfile.Version,
		Deps:    make([]model.DependencyItem, 0),
	}
	for i := range parsedLockfile.Deps {
		parsedLockfile.Deps[i].IsDirectDependency = true
	}
	root := lockfile._v3Conv("", "", make(map[string]struct{}))
	if root != nil {
		parsedLockfile.Deps = root.Dependencies
	}
	return &parsedLockfile, nil
}

func (v *v3Lockfile) _v3Conv(rp string, name string, visited map[string]struct{}) *model.DependencyItem {
	if _, ok := visited[name]; ok {
		return nil
	}
	visited[name] = struct{}{}
	defer func() {
		delete(visited, name)
	}()
	key := path.Join(rp, "node_modules", name)
	pkg, ok := v.Packages[key]
	if !ok {
		key = path.Join(rp, name)
		pkg, ok = v.Packages[key]
	}
	if !ok {
		key = path.Join("node_modules", name)
		pkg, ok = v.Packages[key]
	}
	if !ok {
		key = name
		pkg, ok = v.Packages[name]
	}
	if !ok {
		return nil
	}
	if pkg.Dev {
		return nil
	}
	var item = &model.DependencyItem{
		Component: model.Component{
			CompName:    pkg.Name,
			CompVersion: pkg.Version,
			EcoRepo:     EcoRepo,
		},
	}
	if item.CompName == "" {
		item.CompName = name
	}
	if pkg.Dependencies != nil {
		for s := range pkg.Dependencies {
			if s == "" {
				continue
			}
			r := v._v3Conv(key, s, visited)
			if r == nil {
				continue
			}
			item.Dependencies = append(item.Dependencies, *r)
		}
	}
	return item
}

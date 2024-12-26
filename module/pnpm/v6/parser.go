package v6

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	"iter"
	"strings"
)

type _LockFile struct {
	LockfileVersion string                  `json:"lockfile_version" yaml:"lockfileVersion"`
	Importers       map[string]_Importer    `json:"importers" yaml:"importers"`
	Dependencies    map[string]_TopLevelPkg `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]_TopLevelPkg `json:"dev_dependencies" yaml:"devDependencies"`
	Packages        map[string]_Package     `json:"packages" yaml:"packages"`
}

type _TopLevelPkg struct {
	Specifier string `json:"specifier" yaml:"specifier"`
	Version   string `json:"version" yaml:"version"`
}

type _Package struct {
	Name                 string            `json:"name" yaml:"name"`
	Version              string            `json:"version" yaml:"version"`
	Dev                  bool              `json:"dev" yaml:"dev"`
	Optional             bool              `json:"optional" yaml:"optional"`
	Dependencies         map[string]string `json:"dependencies" yaml:"dependencies"`
	OptionalDependencies map[string]string `json:"optional_dependencies" yaml:"optionalDependencies"`
}

type _Importer struct {
	Dependencies    map[string]_TopLevelPkg `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]_TopLevelPkg `json:"dev_dependencies" yaml:"devDependencies"`
}

type visitContext struct {
	context.Context
	Lockfile *_LockFile
	Strict   bool
	Prune    map[[2]string]struct{}
	cd       shared.CircularDetector
	Dev      bool
}

func Process(_ctx context.Context, data []byte, strict bool) ([]shared.DepTree, error) {
	var lockfile _LockFile
	var e = shared.ParseYaml(data, &lockfile)
	if e != nil {
		return nil, e
	}
	var f = func(m map[string]_TopLevelPkg, dev bool, tree *shared.DepTree) error {
		var ctx = &visitContext{
			Context:  _ctx,
			Lockfile: &lockfile,
			Strict:   strict,
			Prune:    make(map[[2]string]struct{}),
			Dev:      dev,
		}
		for name, it := range m {
			var dep model.DependencyItem
			e := _visit(ctx, name, it.Version, "", &dep)
			if e != nil {
				return e
			}
			if e == nil {
				tree.Dependencies = append(tree.Dependencies, dep)
			}
		}
		return nil
	}
	if len(lockfile.DevDependencies) != 0 || len(lockfile.Dependencies) != 0 {
		var root = shared.DepTree{Name: ""}
		if e := f(lockfile.Dependencies, false, &root); e != nil {
			return nil, e
		}
		if e := f(lockfile.DevDependencies, true, &root); e != nil {
			return nil, e
		}
		return []shared.DepTree{root}, nil
	} else {
		var r []shared.DepTree
		for relPath, importer := range lockfile.Importers {
			var tree shared.DepTree
			tree.Name = relPath
			if e := f(importer.Dependencies, false, &tree); e != nil {
				return nil, e
			}
			if e := f(importer.DevDependencies, true, &tree); e != nil {
				return nil, e
			}
			r = append(r, tree)
		}
		return r, nil
	}
}

func _visit(ctx *visitContext, name, version, path string, node *model.DependencyItem) error {
	var cDetected = ctx.cd.Contains(name, version)
	if cDetected && ctx.Strict {
		return ctx.cd.Error(name, version)
	}
	if cDetected {
		defer func() { ctx.cd.Leave(name, version) }()
	}
	node.CompName = name
	node.CompVersion = getPureVersion(version)
	if ctx.Dev {
		node.IsOnline.SetOnline(false)
	}
	node.EcoRepo = shared.EcoRepo
	if _, found := ctx.Prune[[2]string{name, version}]; found {
		return nil
	}
	ctx.Prune[[2]string{name, version}] = struct{}{}
	// iterate all possible paths
	for key := range searchPathSeq(path, name, version) {
		var p, ok = ctx.Lockfile.Packages[key]
		if !ok {
			continue
		}
		if p.Version != "" {
			node.CompVersion = p.Version
		}
		if !cDetected {
			var f = func(m map[string]string) error {
				for n, v := range m {
					node.Dependencies = append(node.Dependencies, model.DependencyItem{})
					var key = key + "/" + n + "/" + v
					var e = _visit(ctx, n, v, key, &node.Dependencies[len(node.Dependencies)-1])
					if e != nil {
						node.Dependencies = node.Dependencies[:len(node.Dependencies)-1]
						if ctx.Strict {
							return e
						}
					}
				}
				return nil
			}
			if e := f(p.Dependencies); e != nil {
				return e
			}
			if e := f(p.OptionalDependencies); e != nil {
				return e
			}
		}
		return nil
	}
	if ctx.Strict {
		return fmt.Errorf("package %s@%s not found", name, version)
	}
	return nil
}

func searchPathSeq(lastPath string, name, version string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if !yield(lastPath + name + "@" + version) {
			return
		}
		if !yield(lastPath + "/" + name + "@" + version) {
			return
		}
		var i = strings.Index(lastPath, "(")
		if i > -1 {
			lastPath = lastPath[:i]
			if !yield(lastPath + "/" + name + "@" + version) {
				return
			}
		}
		for len(lastPath) > 0 {
			var i = strings.LastIndex(lastPath, "/")
			if i > -1 {
				lastPath = lastPath[:i]
				if !yield(lastPath + "/" + name + "@" + version) {
					return
				}
			}
		}
		if !yield(name + "@" + version) {
			return
		}
	}
}

func getPureVersion(input string) string {
	var i = strings.Index(input, "(")
	if i > -1 {
		input = input[:i]
	}
	i = strings.LastIndex(input, "@")
	if i > -1 {
		input = input[i+1:]
	}
	return input
}

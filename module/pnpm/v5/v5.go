package v5

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	"gopkg.in/yaml.v3"
	"strings"
)

type Pkg struct {
	Name         string            `json:"name" yaml:"name"`
	Version      string            `json:"version" yaml:"version"`
	Dependencies map[string]string `json:"dependencies" yaml:"dependencies"`
	Dev          bool              `json:"dev" yaml:"dev"`
}

type Importer struct {
	Dependencies    map[string]string `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies" yaml:"devDependencies"`
}

type Lockfile struct {
	Importers       map[string]*Importer `json:"importers" yaml:"importers"`
	Dependencies    map[string]string    `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]string    `json:"devDependencies" yaml:"devDependencies"`
	Packages        map[string]*Pkg      `json:"packages" yaml:"packages"`
}

func ParseLockfile(data []byte) (*Lockfile, error) {
	var lockfile Lockfile
	if e := yaml.Unmarshal(data, &lockfile); e != nil {
		return nil, e
	}
	return &lockfile, nil
}

func AnalyzeDepTree(lock *Lockfile) []shared.DepTree {
	var ctx = treeBuildingCtx{
		cache:  make(map[string]model.DependencyItem),
		pkgSet: *buildPkgPathTree(lock),
	}
	var r []shared.DepTree
	{
		var tree shared.DepTree
		for name, version := range lock.Dependencies {
			r := _buildTree(ctx, name, version)
			if r != nil {
				tree.Dependencies = append(tree.Dependencies, *r)
			}
		}
		for name, version := range lock.DevDependencies {
			r := _buildTree(ctx, name, version)
			if r != nil {
				tree.Dependencies = append(tree.Dependencies, *r)
			}
		}
		r = append(r, tree)
	}
	for name, importer := range lock.Importers {
		var tree shared.DepTree
		tree.Name = name
		for name, version := range importer.Dependencies {
			r := _buildTree(ctx, name, version)
			if r != nil {
				tree.Dependencies = append(tree.Dependencies, *r)
			}
		}
		for name, version := range importer.DevDependencies {
			r := _buildTree(ctx, name, version)
			if r != nil {
				tree.Dependencies = append(tree.Dependencies, *r)
			}
		}
		r = append(r, tree)
	}
	return r
}

type treeBuildingCtx struct {
	cd     *circularDetector
	cache  map[string]model.DependencyItem
	pkgSet pkgSet
}

func _buildTree(ctx treeBuildingCtx, name, version string) *model.DependencyItem {
	if ctx.cd.Has(name, version) {
		return nil
	}
	ctx.cd = ctx.cd.With(name, version)
	var path = "/" + name + "/" + version
	if cached, ok := ctx.cache[path]; ok {
		return &cached
	}
	pkg := ctx.pkgSet.GetByPath(path)
	if pkg == nil {
		return nil
	}
	var dep = model.DependencyItem{
		Component: model.Component{
			CompName:    name,
			CompVersion: version,
			EcoRepo:     model.EcoRepo{Ecosystem: "npm"},
		},
		IsOnline: model.IsOnline{Valid: true, Value: pkg.Dev},
	}
	if pkg.Dependencies != nil {
		for name, version := range pkg.Dependencies {
			r := _buildTree(ctx, name, version)
			if r != nil {
				dep.Dependencies = append(dep.Dependencies, *r)
			}
		}
	}
	ctx.cache[path] = dep
	return &dep
}

type circularDetector struct {
	Name    string
	Version string
	parent  *circularDetector
}

func (c *circularDetector) With(name, version string) *circularDetector {
	return &circularDetector{
		Name:    name,
		Version: version,
		parent:  c,
	}
}

func (c *circularDetector) Has(name, version string) bool {
	if c == nil {
		return false
	}
	if c.Name == name && c.Version == version {
		return true
	}
	return c.parent.Has(name, version)
}

type pkgSet struct {
	tree  pathTreeElement
	nvMap map[[2]string]*Pkg
}

func (p *pkgSet) GetByPath(path string) *Pkg {
	var curr = &p.tree
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
o:
	for _, ps := range segments {
		for _, child := range curr.children {
			if child.Segment == ps {
				curr = &child
				continue o
			}
		}
		for _, child := range curr.children {
			if strings.HasPrefix(child.Segment, ps) {
				curr = &child
				continue o
			}
		}
		curr = nil
		break
	}
	if curr != nil {
		return curr.Pkg
	}
	name, version := getNameVersionFromPath0(path)
	return p.nvMap[[2]string{name, version}]
}

type pathTreeElement struct {
	Segment  string
	Pkg      *Pkg
	children []pathTreeElement
}

func buildPkgPathTree(lock *Lockfile) *pkgSet {
	var root pathTreeElement
	for path, pkg := range lock.Packages {
		var curr = &root
	o:
		for _, ps := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
			for i := range curr.children {
				if curr.children[i].Segment == ps {
					curr = &curr.children[i]
					continue o
				}
			}
			curr.children = append(curr.children, pathTreeElement{Segment: ps})
			curr = &curr.children[len(curr.children)-1]
		}
		curr.Pkg = pkg
	}
	nvMap := make(map[[2]string]*Pkg)
	for path, pkg := range lock.Packages {
		var name, version = getNameVersionFromPath0(path)
		if pkg.Name != "" {
			name = pkg.Name
		}
		if pkg.Version != "" {
			version = pkg.Version
		}
		nvMap[[2]string{name, version}] = pkg
	}
	return &pkgSet{
		tree:  root,
		nvMap: nvMap,
	}
}

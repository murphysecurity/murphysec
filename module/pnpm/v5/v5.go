package v5

import (
	"sort"
	"strings"

	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
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
	pkgIndexes      map[[2]string]*Pkg
}

func Visit[T any](l *Lockfile, importer *Importer, callback shared.GVisitor[T], arg T) (e error) {
	if importer == nil {
		e = _visit(l, nil, l.Dependencies, nil, callback, arg)
		if e != nil {
			return
		}
		e = _visit(l, nil, l.DevDependencies, nil, callback, arg)
	} else {
		e = _visit(l, nil, importer.Dependencies, nil, callback, arg)
		if e != nil {
			return
		}
		e = _visit(l, nil, importer.DevDependencies, nil, callback, arg)
	}
	return
}

func nextCallVisitor[T any](l *Lockfile, parent *shared.GComponent, m map[string]string, cd *circleDetector, visitor shared.GVisitor[T]) shared.DoVisit[T] {
	return func(arg T) error {
		return _visit[T](l, parent, m, cd, visitor, arg)
	}
}

func _visit[T any](l *Lockfile, parent *shared.GComponent, m map[string]string, cd *circleDetector, visitor shared.GVisitor[T], arg T) error {
	keys := make([]string, 0, len(m))
	for n := range m {
		keys = append(keys, n)
	}
	sort.Strings(keys)
	for _, n := range keys {
		v := m[n]
		var pkg = l.findPkg(n, v)
		if pkg == nil {
			continue
		}
		if cd.Has(n, v) {
			continue
		}
		var c = &shared.GComponent{
			Name:    n,
			Version: v,
			Dev:     pkg.Dev,
		}
		var ncv = nextCallVisitor(l, c, pkg.Dependencies, cd.With(n, v), visitor)
		if e := visitor(ncv, parent, c, arg); e != nil {
			return e
		}
	}
	return nil
}

func (p *Pkg) adjustByPath(path string) {
	if p.Name != "" && p.Version != "" {
		return
	}
	name, version, ok := parseNameVersionFromPackagePath(path)
	if !ok {
		return
	}
	if p.Name == "" {
		p.Name = name
	}
	if p.Version == "" {
		p.Version = version
	}
}

func parseNameVersionFromPackagePath(path string) (name, version string, ok bool) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", "", false
	}
	slash := strings.LastIndex(trimmed, "/")
	if slash <= 0 || slash == len(trimmed)-1 {
		return "", "", false
	}

	name = trimmed[:slash]
	versionWithSuffix := trimmed[slash+1:]
	underscore := strings.Index(versionWithSuffix, "_")
	if underscore >= 0 {
		versionWithSuffix = versionWithSuffix[:underscore]
	}
	if name == "" || versionWithSuffix == "" {
		return "", "", false
	}
	return name, versionWithSuffix, true
}

func (l *Lockfile) buildIndexes() {
	paths := make([]string, 0, len(l.Packages))
	for path := range l.Packages {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		pkg := l.Packages[path]
		pkg.adjustByPath(path)
	}
	l.pkgIndexes = make(map[[2]string]*Pkg, len(l.Packages))
	for _, path := range paths {
		pkg := l.Packages[path]
		var name, version = pkg.Name, pkg.Version
		if name == "" {
			continue
		}
		key := [2]string{name, version}
		if merged, ok := l.pkgIndexes[key]; ok {
			mergePkg(merged, pkg)
		} else {
			l.pkgIndexes[key] = clonePkg(pkg)
		}
	}
}

func clonePkg(pkg *Pkg) *Pkg {
	out := &Pkg{
		Name:    pkg.Name,
		Version: pkg.Version,
		Dev:     pkg.Dev,
	}
	if len(pkg.Dependencies) == 0 {
		return out
	}
	out.Dependencies = make(map[string]string, len(pkg.Dependencies))
	for n, v := range pkg.Dependencies {
		out.Dependencies[n] = v
	}
	return out
}

func mergePkg(merged, pkg *Pkg) {
	// If any context is non-dev, keep the merged node as non-dev.
	merged.Dev = merged.Dev && pkg.Dev
	if len(pkg.Dependencies) == 0 {
		return
	}
	if merged.Dependencies == nil {
		merged.Dependencies = make(map[string]string, len(pkg.Dependencies))
	}
	for n, v := range pkg.Dependencies {
		if _, exists := merged.Dependencies[n]; !exists {
			merged.Dependencies[n] = v
		}
	}
}

func (l *Lockfile) findPkg(name, version string) (p *Pkg) {
	p = l.pkgIndexes[[2]string{name, version}]
	if p == nil {
		p = l.Packages["/"+name+"/"+version]
	}
	if p == nil {
		p = l.Packages[version]
	}
	if p == nil {
		var us = strings.LastIndex(version, "_")
		if us > -1 {
			p = l.findPkg(name, version[:us])
		}
	}
	return
}

type circleDetector struct {
	Name    string
	Version string
	Parent  *circleDetector
}

func (c *circleDetector) Has(name, version string) bool {
	if c == nil {
		return false
	}
	if c.Name == name && c.Version == version {
		return true
	}
	return c.Parent.Has(name, version)
}

func (c *circleDetector) With(name, version string) *circleDetector {
	return &circleDetector{
		Name:    name,
		Version: version,
		Parent:  c,
	}
}

func ParseLockfile(data []byte) (*Lockfile, error) {
	var lockfile Lockfile
	// workaround for unquoted question-mark in inline object
	if e := shared.ParseYaml(data, &lockfile); e != nil {
		return nil, e
	}
	lockfile.buildIndexes()
	return &lockfile, nil
}

func BuildDepTree(l *Lockfile, importer *Importer, importName string) *shared.DepTree {
	var di model.DependencyItem
	_ = Visit(l, importer, v5Visitor(), &di)
	var d = shared.DepTree{Name: importName, Dependencies: di.Dependencies}
	if len(d.Dependencies) == 0 {
		return nil
	}
	return &d
}

func v5Visitor() func(visitor shared.DoVisit[*model.DependencyItem], _ *shared.GComponent, child *shared.GComponent, arg *model.DependencyItem) error {
	var pruneSet = make(map[shared.GComponent]struct{})
	return func(visitor shared.DoVisit[*model.DependencyItem], _ *shared.GComponent, child *shared.GComponent, arg *model.DependencyItem) error {
		var c = model.DependencyItem{
			Component: model.Component{
				CompName:    child.Name,
				CompVersion: child.Version,
				EcoRepo:     shared.EcoRepo,
			},
		}
		c.IsOnline.SetOnline(!child.Dev)
		var e error
		if _, ok := pruneSet[*child]; !ok {
			pruneSet[*child] = struct{}{}
			e = visitor(&c)
			arg.Dependencies = append(arg.Dependencies, c)
		}
		return e
	}
}

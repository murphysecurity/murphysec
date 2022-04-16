package maven

import (
	"fmt"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"path/filepath"
	"sync"
)

type Dependency struct {
	Coordinate
	Children []Dependency `json:"children,omitempty"`
}

func (d Dependency) String() string {
	return fmt.Sprintf("%v: %v", d.Coordinate, d.Children)
}

var MvnSkipped = base.NewInspectError("java", "Mvn inspect is skipped, please check you maven environment.")

func ScanMavenProject(dir string) ([]base.Module, error) {
	var modules []base.Module
	var deps map[Coordinate][]Dependency
	moduleFileMapping := map[Coordinate]string{}
	var e error
	// check maven version, skip maven scan if check fail
	doMvnScan, mvnVer := checkMvnEnv()
	if doMvnScan {
		deps, e = scanMvnDependency(dir)
		if e != nil {
			logger.Err.Printf("mvn scan failed: %+v\n", e)
		}
	}
	// analyze pom file
	{
		if deps == nil {
			deps = map[Coordinate][]Dependency{}
		}
		pomFiles := InspectModule(dir)
		logger.Info.Printf("scanned pom modules: %d", len(pomFiles))
		resolver := NewResolver()
		for _, builder := range pomFiles {
			{
				pf := resolver.ResolveLocally(builder, nil)
				if pf == nil {
					continue
				}
				moduleFileMapping[pf.coordinate] = pf.path
				if len(deps[pf.coordinate]) > 0 {
					continue
				}
			}
			pf := resolver.Resolve(builder, nil)
			if pf == nil {
				continue
			}
			if !pf.coordinate.Complete() {
				logger.Info.Println("local pom coordinate can't be resolve", pf.coordinate)
				continue
			}
			analyzer := NewDepTreeAnalyzer(resolver)
			graph := analyzer.Resolve(pf)
			logger.Info.Println("dep graph")
			logger.Info.Println(graph.DOT())
			deps[pf.coordinate] = graph.Tree(pf.coordinate)
		}
	}
	for coordinate, dependencies := range deps {
		modules = append(modules, base.Module{
			PackageManager: "maven",
			Language:       "java",
			PackageFile:    "pom.xml",
			Name:           coordinate.Name(),
			Version:        coordinate.Version,
			FilePath:       filepath.Join(moduleFileMapping[coordinate], "pom.xml"),
			Dependencies:   convDeps(dependencies),
			RuntimeInfo:    mvnVer,
		})
	}
	if len(modules) == 0 && !doMvnScan {
		return nil, MvnSkipped
	}
	return modules, nil
}

func convDeps(deps []Dependency) []base.Dependency {
	rs := make([]base.Dependency, 0)
	for _, it := range deps {
		d := _convDep(it)
		if d == nil {
			continue
		}
		rs = append(rs, *d)
	}
	return rs
}

func _convDep(dep Dependency) *base.Dependency {
	if dep.GroupId == "" || dep.ArtifactId == "" || dep.Version == "" {
		return nil
	}
	d := &base.Dependency{
		Name:         dep.Name(),
		Version:      dep.Version,
		Dependencies: []base.Dependency{},
	}
	for _, it := range dep.Children {
		dd := _convDep(it)
		if dd == nil {
			continue
		}
		d.Dependencies = append(d.Dependencies, *dd)
	}
	return d
}

type DepTreeCacheMap struct {
	m sync.Map
}

func (d *DepTreeCacheMap) Get(coor Coordinate) *Dependency {
	v, _ := d.m.Load(coor)
	if vv, ok := v.(*Dependency); ok {
		return vv
	}
	return nil
}

func (d *DepTreeCacheMap) Put(coor Coordinate, tree *Dependency) {
	d.m.Store(coor, tree)
}

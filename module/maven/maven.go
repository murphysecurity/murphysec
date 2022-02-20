package maven

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
	"sync"
)

type Dependency struct {
	fmt.Stringer
	Coordinate
	Children []Dependency `json:"children,omitempty"`
}

func (d Dependency) String() string {
	return fmt.Sprintf("%v: %v", d.Coordinate, d.Children)
}

func ScanMavenProject(dir string) ([]base.Module, error) {
	var modules []base.Module
	var deps map[Coordinate][]Dependency
	moduleFileMapping := map[Coordinate]string{}
	var e error
	// check maven version, skip maven scan if check fail
	skipMvnScan, mvnVer := checkMvnEnv()
	if skipMvnScan {
		deps, e = scanMvnDependency(dir)
		if e != nil {
			logger.Err.Printf("mvn scan failed: %+v\n", e)
		}
	}
	// analyze pom file
	repo, e := NewProjectRepoFromDir(dir)
	if e != nil {
		logger.Err.Println("Scan pom file failed")
		return nil, errors.Wrap(e, "scan project pom failed")
	}
	// fill moduleFileMapping
	for _, info := range repo.ListModuleInfo() {
		relPath := must.String(filepath.Rel(dir, info.FilePath))
		moduleFileMapping[info.Coordinate()] = relPath
		logger.Debug.Println("Module path mapping:", info.Coordinate().String(), relPath)
	}
	if len(repo.ListModuleInfo()) == 0 {
		logger.Debug.Println("No module found")
	}
	resolver := NewResolver(repo)
	m2Settings := ReadM2SettingMirror()
	if m2Settings == nil {
		resolver.AddRepo(NewLocalRepo(""))
	} else {
		resolver.AddRepo(NewLocalRepo(m2Settings.RepoPath))
	}
	if m2Settings == nil || len(m2Settings.Mirrors) == 0 {
		resolver.AddRepo(DefaultMavenRepo()...)
	} else {
		for _, it := range m2Settings.Mirrors {
			resolver.AddRepo(MustNewHttpRepo(it))
		}
	}
	if len(deps) == 0 {
		logger.Warn.Println("Maven command execute failed, use another tools")
		deps = map[Coordinate][]Dependency{}
		cacheMap := &DepTreeCacheMap{}
		for _, info := range repo.ListModuleInfo() {
			logger.Debug.Println("Resolving module", info.Coordinate())
			pomFile, e := resolver.ResolvePomFile(nil, info.Coordinate())
			if e != nil {
				logger.Err.Println("Resolve local module failed", info.Coordinate(), e.Error())
				continue
			}
			p := _resolve(nil, resolver, pomFile, cacheMap, nil, 3)
			if p == nil {
				logger.Info.Println("Resolve pom dependency failed.", info.PomFile.Coordinate().String())
			} else {
				deps[info.Coordinate()] = p.Children
			}
		}
	}
	for coordinate, dependencies := range deps {
		modules = append(modules, base.Module{
			PackageManager: "Maven",
			Language:       "java",
			PackageFile:    "pom.xml",
			Name:           coordinate.Name(),
			Version:        coordinate.Version,
			RelativePath:   moduleFileMapping[coordinate],
			Dependencies:   convDeps(dependencies),
			RuntimeInfo:    mvnVer,
		})
	}
	return modules, nil
}

func convDeps(deps []Dependency) []base.Dependency {
	var rs []base.Dependency
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

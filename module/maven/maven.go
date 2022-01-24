package maven

import (
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "MavenInspector@" + i.Version()
}

func (i *Inspector) Version() string {
	return "v0.0.1"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pom.xml"))
}

func (i *Inspector) Inspect(dir string) ([]base.Module, error) {
	return ScanMavenProject(dir)
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMMaven
}

type Dependency struct {
	Coordinate
	Children []Dependency
}

func ScanMavenProject(dir string) ([]base.Module, error) {
	var modules []base.Module
	var deps map[Coordinate][]Dependency
	mvnVer, e := checkMvnVersion()
	skipMvnScan := false
	if e != nil {
		logger.Err.Println("Check mvn version failed", e.Error())
		logger.Err.Println("Skip maven scan")
		skipMvnScan = true
	}
	if !skipMvnScan {
		deps, e = scanMvnDependency(dir)
		if e != nil {
			logger.Err.Println("mvn scan failed.", e.Error())
		}
	}
	poms, e := ReadMavenProject(dir)
	if e != nil {
		logger.Err.Println("Read mvn project failed.", e.Error())
		return nil, errors.Wrap(e, "read maven project failed")
	}
	if deps == nil {
		deps = map[Coordinate][]Dependency{}
		for _, it := range poms {
			for _, d := range it.Dependencies() {
				deps[it.Coordination()] = append(deps[it.Coordination()], Dependency{
					Coordinate: d,
				})
			}
		}
	}
	for _, it := range poms {
		coor := it.Coordination()
		modules = append(modules, base.Module{
			PackageManager: "Maven",
			Language:       "Java",
			PackageFile:    "pom.xml",
			Name:           coor.GroupId + ":" + coor.ArtifactId,
			Version:        coor.Version,
			RelativePath:   must.String(filepath.Rel(dir, it.Path)),
			Dependencies:   _convDepToModule(deps[coor]),
			RuntimeInfo:    mvnVer,
		})
	}
	return modules, nil
}

func _convDepToModule(deps []Dependency) []base.Dependency {
	var rs []base.Dependency
	for _, it := range deps {
		rs = append(rs, base.Dependency{
			Name:         it.GroupId + ":" + it.ArtifactId,
			Version:      it.Version,
			Dependencies: _convDepToModule(it.Children),
		})
	}
	return rs
}

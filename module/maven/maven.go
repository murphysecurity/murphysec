package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
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

var MvnSkipped = model.NewInspectError(model.Java, "Mvn inspect is skipped, please check you maven environment.")

func ScanMavenProject(ctx context.Context, task *model.InspectorTask) ([]model.Module, error) {
	log := utils.UseLogger(ctx)
	dir := task.ScanDir
	var modules []model.Module
	var deps map[Coordinate][]Dependency
	moduleFileMapping := map[Coordinate]string{}
	var e error
	var doMvnScan bool
	// check maven version, skip maven scan if check fail
	mvnCmdInfo, e := CheckMvnCommand()
	if e != nil {
		log.Sugar().Warnf("Mvn command not found %v", e)
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到您的环境中 Maven 无法正常运行，可能会导致检测结果不完整，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
	} else {
		log.Sugar().Infof("Mvn command found: %s", mvnCmdInfo)
		doMvnScan = true
	}

	var useBackupResolver = !doMvnScan
	if doMvnScan {
		deps, e = ScanMvnDeps(ctx, mvnCmdInfo)
		if e != nil {
			task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 Maven获取依赖信息失败，可能会导致检测结果不完整或失败，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
			logger.Err.Printf("mvn scan failed: %+v\n", e)
			useBackupResolver = true
		}
	}
	// analyze pom file
	if useBackupResolver {
	}
	for coordinate, dependencies := range deps {
		modules = append(modules, model.Module{
			PackageManager: model.PMMaven,
			Language:       model.Java,
			PackageFile:    "pom.xml",
			Name:           coordinate.Name(),
			Version:        coordinate.Version,
			FilePath:       filepath.Join(moduleFileMapping[coordinate], "pom.xml"),
			Dependencies:   convDeps(dependencies),
			RuntimeInfo:    mvnCmdInfo,
		})
	}
	if len(modules) == 0 && !doMvnScan {
		return nil, MvnSkipped
	}
	return modules, nil
}

func convDeps(deps []Dependency) []model.Dependency {
	rs := make([]model.Dependency, 0)
	for _, it := range deps {
		d := _convDep(it)
		if d == nil {
			continue
		}
		rs = append(rs, *d)
	}
	return rs
}

func _convDep(dep Dependency) *model.Dependency {
	if dep.GroupId == "" || dep.ArtifactId == "" || dep.Version == "" {
		return nil
	}
	d := &model.Dependency{
		Name:         dep.Name(),
		Version:      dep.Version,
		Dependencies: []model.Dependency{},
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

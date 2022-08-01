package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"sync"
)

type Dependency struct {
	Coordinate
	Children []Dependency `json:"children,omitempty"`
}

func (d Dependency) String() string {
	return fmt.Sprintf("%v: %v", d.Coordinate, d.Children)
}

func ScanMavenProject(ctx context.Context, task *model.InspectorTask) ([]model.Module, error) {
	log := utils.UseLogger(ctx)
	dir := task.ScanDir
	var modules []model.Module
	var e error
	var useBackupResolver = false
	var deps *DepsMap
	// check maven version, skip maven scan if check fail
	mvnCmdInfo, e := CheckMvnCommand()
	if e != nil {
		useBackupResolver = true
		log.Sugar().Warnf("Mvn command not found %v", e)
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到您的环境中 Maven 无法正常运行，可能会导致检测结果不完整，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
	} else {
		log.Sugar().Infof("Mvn command found: %s", mvnCmdInfo)
		var e error
		deps, e = ScanDepsByPluginCommand(ctx, dir, mvnCmdInfo)
		if e != nil {
			log.Error("Scan maven dependencies failed", zap.Error(e))
			useBackupResolver = true
			task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 Maven获取依赖信息失败，可能会导致检测结果不完整或失败，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
		}
	}

	// analyze pom file
	if useBackupResolver {
		var e error
		deps, e = BackupResolve(ctx, dir)
		if e != nil {
			log.Error("Use backup resolver failed", zap.Error(e))
		}
	}
	if deps == nil {
		return nil, ErrInspection
	}
	for _, entry := range deps.ListAllEntries() {
		modules = append(modules, model.Module{
			PackageManager: model.PMMaven,
			Language:       model.Java,
			PackageFile:    "pom.xml",
			Name:           entry.coordinate.Name(),
			Version:        entry.coordinate.Version,
			FilePath:       entry.relativePath,
			Dependencies:   convDeps(entry.children),
			RuntimeInfo:    mvnCmdInfo,
		})
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

package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
	"path/filepath"
)

type Dependency struct {
	Coordinate
	Children []Dependency `json:"children,omitempty"`
	Scope    string       `json:"scope"`
}

func (d Dependency) IsZero() bool {
	return len(d.Children) == 0 && d.ArtifactId == "" && d.GroupId == "" && d.Version == ""
}

func (d Dependency) String() string {
	return fmt.Sprintf("%v: %v", d.Coordinate, d.Children)
}

func ScanMavenProject(ctx context.Context, task *model.InspectionTask) ([]model.Module, error) {
	log := logctx.Use(ctx)
	dir := task.Dir()
	var modules []model.Module
	var e error
	var useBackupResolver = false
	var deps *DepsMap
	// check maven version, skip maven scan if check fail
	mvnCmdInfo, e := CheckMvnCommand(ctx)
	if e != nil {
		useBackupResolver = true
		log.Sugar().Warnf("Mvn command not found %v", e)
	} else {
		log.Sugar().Infof("Mvn command found: %s", mvnCmdInfo)
		var e error
		deps, e = ScanDepsByPluginCommand(ctx, dir, mvnCmdInfo)
		if e != nil {
			log.Error("Scan maven dependencies failed", zap.Error(e))
			useBackupResolver = true
		}
	}

	// analyze pom file
	env.ScannerShouldEnableMavenBackupScan = useBackupResolver || deps == nil || deps.allEmpty()
	if useBackupResolver {
		if env.ScannerScan {
			return nil, nil
		} else {
			ui.Use(ctx).Display(ui.MsgWarn, "通过 Maven获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")
			var e error
			deps, e = BackupResolve(ctx, dir)
			if e != nil {
				log.Error("Use backup resolver failed", zap.Error(e))
			}
		}
	}
	if deps == nil {
		return nil, ErrInspection
	}

	var strategy = model.ScanStrategyNormal
	if useBackupResolver {
		strategy = model.ScanStrategyBackup
	}
	for _, entry := range deps.ListAllEntries() {
		task.AddModule(model.Module{
			PackageManager: "maven",
			ModuleName:     entry.coordinate.Name(),
			ModuleVersion:  entry.coordinate.Version,
			ModulePath:     filepath.Join(dir, entry.relativePath),
			Dependencies:   convDeps(entry.children),
			ScanStrategy:   strategy,
		})
	}
	return modules, nil
}

func convDeps(deps []Dependency) []model.DependencyItem {
	rs := make([]model.DependencyItem, 0)
	for _, it := range deps {
		d := _convDep(it)
		if d == nil {
			continue
		}
		d.IsDirectDependency = true
		rs = append(rs, *d)
	}
	return rs
}

func _convDep(dep Dependency) *model.DependencyItem {
	if dep.IsZero() {
		return nil
	}
	d := &model.DependencyItem{
		Component: model.Component{
			CompName:    dep.Name(),
			CompVersion: dep.Version,
			EcoRepo:     EcoRepo,
		},
		IsOnline:   model.IsOnlineTrue(),
		MavenScope: dep.Scope,
	}
	if d.MavenScope == "test" || d.MavenScope == "provided" || d.MavenScope == "system" {
		d.IsOnline.SetOnline(false)
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

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}

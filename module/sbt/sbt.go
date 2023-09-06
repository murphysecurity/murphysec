package sbt

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "SBT"
}

func (i Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "build.sbt"))
}

func (i Inspector) InspectProject(ctx context.Context) error {
	if env.DoNotBuild {
		return nil
	}
	task := model.UseInspectionTask(ctx)
	dep, e := sbtDependencyTree(ctx, task.Dir())
	if e != nil {
		return fmt.Errorf("sbt command: %w", e)
	}
	module := model.Module{
		PackageManager: "sbt",
		ModulePath:     filepath.Join(task.Dir(), "build.sbt"),
		Dependencies:   mapToModel(dep),
		ScanStrategy:   model.ScanStrategyNormal,
	}
	if len(dep) > 0 {
		module.ModuleName = dep[0].Name
		module.ModuleVersion = dep[0].Version
	}
	task.AddModule(module)
	return nil
}

func (i Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}

type Dep struct {
	Name     string
	Version  string
	Children []Dep
}

func mapToModel(deps []Dep) []model.DependencyItem {
	r := make([]model.DependencyItem, len(deps))
	for i := range deps {
		r[i] = model.DependencyItem{
			Component: model.Component{
				CompName:    deps[i].Name,
				CompVersion: deps[i].Version,
				EcoRepo:     EcoRepo,
			},
			Dependencies: mapToModel(deps[i].Children),
		}
	}
	return r
}

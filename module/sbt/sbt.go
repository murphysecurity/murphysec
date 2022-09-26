package sbt

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
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
	task := model.UseInspectorTask(ctx)
	dep, e := sbtDependencyTree(ctx, task.ScanDir)
	if e != nil {
		return fmt.Errorf("sbt command: %w", e)
	}
	module := model.Module{
		PackageManager: model.PmSbt,
		Language:       model.Scala,
		Name:           "",
		Version:        "",
		RelativePath:   filepath.Join(task.ScanDir, "build.sbt"),
		Dependencies:   dep,
		ScanStrategy:   model.ScanStrategyNormal,
	}
	if len(dep) > 0 {
		module.Name = dep[0].Name
		module.Version = dep[0].Version
	}
	task.AddModule(module)
	return nil
}

func (i Inspector) SupportFeature(feature base.InspectorFeature) bool {
	return false
}

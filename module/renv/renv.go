package renv

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/simplejson"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (Inspector) String() string {
	return "REnv"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "renv.lock"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	inspectTask := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	data, e := os.ReadFile(filepath.Join(inspectTask.ScanDir, "renv.lock"))
	if e != nil {
		return fmt.Errorf("read renv.lock: %w", e)
	}
	j, e := simplejson.NewJSON(data)
	if e != nil || j == nil {
		return fmt.Errorf("parse renv.lock: %w", e)
	}
	var deps []model.Dependency
	for _, it := range j.Get("Packages").JSONMap() {
		if it == nil {
			continue
		}
		var name = it.Get("Package").String()
		var version = it.Get("Version").String()
		deps = append(deps, model.Dependency{
			Name:    name,
			Version: version,
		})
	}
	if len(deps) == 0 {
		logger.Warn("No valid package item found")
		return nil
	}
	inspectTask.AddModule(model.Module{
		PackageManager: "",
		Language:       "R",
		Name:           "RProject",
		Version:        "",
		RelativePath:   filepath.Join(inspectTask.ScanDir, "renv.lock"),
		Dependencies:   deps,
		RuntimeInfo:    nil,
		ScanStrategy:   "",
	})
	return nil
}

func (Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

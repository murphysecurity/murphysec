package go_mod

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type goPkgLock struct {
	Projects []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `toml:"projects"`
}

// ParseGopkgLock 解析Gopkg.lock文件
func parseGopkgLock(dir string) (goPkgLock, error) {
	var m goPkgLock
	file, err := os.Open(dir)
	if err != nil {
		return m, err
	}
	_, err = toml.NewDecoder(file).Decode(&m)
	if err != nil {
		return m, err
	}
	return m, nil
}
func parserGoPkgLock(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	goPkgLockPath := filepath.Join(task.Dir(), "Gopkg.lock")
	pkg, e := parseGopkgLock(goPkgLockPath)
	if e != nil {
		logger.Error("parse gopkg.lock error :", zap.Error(e))
		return e
	}
	var dependencies []model.DependencyItem
	for _, j := range pkg.Projects {
		dependencies = append(dependencies, model.DependencyItem{
			Component: model.Component{
				CompName:    j.Name,
				CompVersion: j.Version,
				EcoRepo:     EcoRepo,
			},
			IsDirectDependency: true,
		})
	}

	m := model.Module{
		PackageManager: "gopkg",
		ModulePath:     goPkgLockPath,
		ModuleName:     task.Dir(),
		Dependencies:   dependencies,
	}
	task.AddModule(m)
	return nil
}

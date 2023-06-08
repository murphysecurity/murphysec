package bundler

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Bundler"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Gemfile")) && utils.IsFile(filepath.Join(dir, "Gemfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	scanDir := task.Dir()
	gemFile := filepath.Join(scanDir, "Gemfile")
	gemLockFile := filepath.Join(scanDir, "Gemfile.lock")
	if !utils.IsFile(gemFile) || !utils.IsFile(gemLockFile) {
		return nil
	}
	logger.Debug("Reading Gemfile.lock", zap.String("path", gemLockFile))
	data, e := utils.ReadFileLimited(gemLockFile, 1024*1024*4)
	if e != nil {
		return errors.WithMessage(e, "Read Gemfile.lock failed")
	}
	tree, e := getDepGraph(string(data))
	if e != nil {
		return errors.WithMessage(e, "Parse Gemfile.lock failed")
	}
	task.AddModule(model.Module{
		PackageManager: "bundler",
		ModuleName:     tree[0].CompName,
		Dependencies:   tree,
		ModulePath:     gemFile,
	})
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "bundler",
	Repository: "",
}

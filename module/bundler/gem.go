package bundler

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Bundler"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Gemfile")) && utils.IsFile(filepath.Join(dir, "Gemfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	scanDir := task.ScanDir
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
		PackageManager: model.PMBundler,
		Language:       model.Ruby,
		Name:           tree[0].Name,
		Dependencies:   tree,
		RelativePath:   scanDir,
	})
	return nil
}

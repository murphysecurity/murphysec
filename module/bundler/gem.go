package bundler

import (
	"context"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"path/filepath"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "BundlerInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Gemfile")) && utils.IsFile(filepath.Join(dir, "Gemfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	scanDir := task.ScanDir
	gemFile := filepath.Join(scanDir, "Gemfile")
	gemLockFile := filepath.Join(scanDir, "Gemfile.lock")
	if !utils.IsFile(gemFile) || !utils.IsFile(gemLockFile) {
		return nil
	}
	logger.Info.Println("RubyGems inspect: ", scanDir)
	data, e := utils.ReadFileLimited(gemLockFile, 1024*1024*4)
	if e != nil {
		return errors.Wrap(e, "ReadRubyGemsLockFile")
	}
	tree, e := getDepGraph(string(data))
	if e != nil {
		return errors.Wrap(e, "ParseGemLockFile")
	}
	task.AddModule(model.Module{
		PackageManager: model.PMBundler,
		Language:       model.Ruby,
		PackageFile:    "Gemfile.lock",
		Name:           tree[0].Name,
		Dependencies:   tree,
	})
	return nil
}

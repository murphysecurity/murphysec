package bundler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Bundler"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Gemfile")) || utils.IsFile(filepath.Join(dir, "Gemfile.lock"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	scanDir := task.Dir()
	gemFile := filepath.Join(scanDir, "Gemfile")
	gemLockFile := filepath.Join(scanDir, "Gemfile.lock")
	if utils.IsFile(gemFile) && !utils.IsFile(gemLockFile) {
		var m gemfile
		var dep []model.DependencyItem
		file, _ := os.Open(gemFile)
		m.Parse(file)
		for _, j := range m.Gems {
			dep = append(dep, model.DependencyItem{
				Component: model.Component{
					CompName:    j.Name,
					CompVersion: j.Version,
					EcoRepo:     EcoRepo,
				},
				DependencyRelation: model.DependencyRelationDirect,
			})
		}
		task.AddModule(model.Module{
			PackageManager: "bundler",
			ModuleName:     "Gemfile",
			Dependencies:   dep,
			ModulePath:     gemFile,
		})
		return nil
	}
	if !utils.IsFile(gemFile) && !utils.IsFile(gemLockFile) {
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

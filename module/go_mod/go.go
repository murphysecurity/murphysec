package go_mod

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (i *Inspector) String() string {
	return "GoMod"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := utils.UseLogger(ctx)
	modFilePath := filepath.Join(task.Dir(), "go.mod")
	logger.Debug("Reading go.mod", zap.String("path", modFilePath))
	data, e := utils.ReadFileLimited(modFilePath, 1024*1024*4)
	if e != nil {
		return errors.WithMessage(e, "Open GoMod file")
	}
	logger.Debug("Parsing go.mod")
	f, e := modfile.ParseLax(filepath.Base(modFilePath), data, nil)
	if e != nil {
		return errors.WithMessage(e, "Parse go mod failed")
	}
	m := model.Module{
		PackageManager: "gomod",
		ModulePath:     modFilePath,
		ModuleName:     "<NoNameModule>",
	}
	if f.Module != nil {
		m.ModuleVersion = f.Module.Mod.Version
		m.ModulePath = f.Module.Mod.Path
	}

	for _, it := range f.Require {
		if it == nil {
			continue
		}
		m.Dependencies = append(m.Dependencies, model.DependencyItem{
			Component: model.Component{
				CompName:    it.Mod.Path,
				CompVersion: it.Mod.Version,
				EcoRepo:     EcoRepo,
			},
		})
	}
	task.AddModule(m)
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "gomod",
	Repository: "",
}

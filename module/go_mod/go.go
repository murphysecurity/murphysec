package go_mod

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) String() string {
	return "GoModInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	modFilePath := filepath.Join(task.ScanDir, "go.mod")
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
		PackageManager: model.PMGoMod,
		Language:       model.Go,
		PackageFile:    "go.mod",
		FilePath:       modFilePath,
		Name:           "<NoNameModule>",
	}
	if f.Module != nil {
		m.Version = f.Module.Mod.Version
		m.Name = f.Module.Mod.Path
	}

	for _, it := range f.Require {
		if it == nil {
			continue
		}
		m.Dependencies = append(m.Dependencies, model.Dependency{
			Name:    it.Mod.Path,
			Version: it.Mod.Version,
		})
	}
	task.AddModule(m)
	return nil
}

func New() base.Inspector {
	return &Inspector{}
}

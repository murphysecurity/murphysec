package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (*Inspector) String() string {
	return "Conan"
}

func (*Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "conanfile.txt")) ||
		utils.IsFile(filepath.Join(dir, "conanfile.py")) ||
		utils.IsFile(filepath.Join(dir, "conan.txt")) ||
		utils.IsFile(filepath.Join(dir, "conan.py"))
}
func (*Inspector) InspectProject(ctx context.Context) error {
	if env.DoNotBuild {
		return nil
	}
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	cmdInfo, e := getConanInfo(ctx)
	if e != nil {
		return e
	}
	jsonFilePath, e := ExecuteConanInfoCmd(ctx, cmdInfo.Path, task.Dir())

	var conanErr conanError
	if errors.As(e, &conanErr) {
		badConanView(ctx)
		printConanError(ctx, &conanErr)
		return e
	}
	if e != nil {
		return e
	}
	defer func() {
		if e := os.Remove(jsonFilePath); e != nil {
			logger.Error("Can't remove temp file", zap.Error(e), zap.Any("path", jsonFilePath))
		}
	}()
	var conanJson _ConanInfoJsonFile
	if e := conanJson.ReadFromFile(jsonFilePath); e != nil {
		return e
	}
	t, e := conanJson.Tree()
	if e != nil {
		return e
	}
	task.AddModule(model.Module{
		PackageManager: "conan",
		ModuleName:     "conanfile.txt",
		ModulePath:     filepath.Join(task.Dir(), "conanfile.txt"),
		Dependencies:   t.Dependencies,
	})
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "conan",
	Repository: "",
}

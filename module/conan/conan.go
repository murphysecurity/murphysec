package conan

import (
	"context"
	"github.com/google/uuid"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (*Inspector) String() string {
	return "ConanInspector"
}

func (*Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "conanfile.txt"))
}
func (*Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	cmdInfo, e := getConanInfo(ctx)
	if e != nil {
		return e
	}
	jsonFilePath, e := ExecuteConanInfoCmd(ctx, cmdInfo.Path, task.ScanDir)
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
		PackageManager: model.PmConan,
		Language:       model.Cxx,
		PackageFile:    "conanfile.txt",
		Name:           "conanfile.txt",
		Version:        "",
		FilePath:       filepath.Join(task.ScanDir, "conanfile.txt"),
		Dependencies:   t.Dependencies,
		RuntimeInfo:    cmdInfo,
		UUID:           uuid.UUID{},
	})
	return nil
}

func New() base.Inspector {
	return &Inspector{}
}

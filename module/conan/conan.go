package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Conan"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, "conanfile.txt")) ||
		utils.IsFile(filepath.Join(dir, "conanfile.py")) ||
		utils.IsFile(filepath.Join(dir, "conan.txt")) ||
		utils.IsFile(filepath.Join(dir, "conan.py"))
}
func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	registeredAutoBuild := task.RegisterAutoBuild()
	if env.DoNotBuild {
		scanerr.Add(ctx, scanerr.Param{Kind: scanerr.KindBuildDisabled})
		registeredAutoBuild.MarkDisabled()
		return nil
	}
	logger := logctx.Use(ctx)
	cmdInfo, e := getConanInfo(ctx)
	if e != nil {
		registeredAutoBuild.MarkFailed()
		kind := scanerr.KindConanFailed
		if errors.Is(e, ErrConanNotFound) {
			kind = scanerr.KindConanNotFound
		}
		scanerr.Add(ctx, scanerr.Param{
			Kind:    kind,
			Content: e.Error(),
		})
		return e
	}
	jsonFilePath, jsonKind, e := ExecuteConanInfoCmd(ctx, cmdInfo, task.Dir())

	var conanErr conanError
	if errors.As(e, &conanErr) {
		registeredAutoBuild.MarkFailed()
		scanerr.Add(ctx, scanerr.Param{
			Kind:    scanerr.KindConanFailed,
			Content: conanErr.Error(),
		})
		if !env.ScannerScan {
			badConanView(ctx)
			printConanError(ctx, &conanErr)
		}
		return e
	}
	if e != nil {
		registeredAutoBuild.MarkFailed()
		scanerr.Add(ctx, scanerr.Param{
			Kind:    scanerr.KindConanFailed,
			Content: e.Error(),
		})
		return e
	}
	defer func() {
		if e := os.Remove(jsonFilePath); e != nil {
			logger.Error("Can't remove temp file", zap.Error(e), zap.Any("path", jsonFilePath))
		}
	}()
	var t *model.DependencyItem
	switch jsonKind {
	case ConanJsonKindGraph:
		var conanGraphJson _ConanGraphInfoJsonFile
		if e := conanGraphJson.ReadFromFile(jsonFilePath); e != nil {
			registeredAutoBuild.MarkFailed()
			scanerr.Add(ctx, scanerr.Param{
				Kind:    scanerr.KindConanFailed,
				Content: e.Error(),
			})
			return e
		}
		t, e = conanGraphJson.Tree()
	default:
		var conanJson _ConanInfoJsonFile
		if e := conanJson.ReadFromFile(jsonFilePath); e != nil {
			registeredAutoBuild.MarkFailed()
			scanerr.Add(ctx, scanerr.Param{
				Kind:    scanerr.KindConanFailed,
				Content: e.Error(),
			})
			return e
		}
		t, e = conanJson.Tree()
	}
	if e != nil {
		registeredAutoBuild.MarkFailed()
		scanerr.Add(ctx, scanerr.Param{
			Kind:    scanerr.KindConanFailed,
			Content: e.Error(),
		})
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

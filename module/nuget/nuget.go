package nuget

import (
	"context"

	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Nuget"
}

func (Inspector) CheckDir(dir string) bool {
	// return utils.IsFile(filepath.Join(dir, "packages.config"))
	return utils.IsDir(dir)
}

func (Inspector) InspectProject(ctx context.Context) error {

	task := model.UseInspectionTask(ctx)
	packagesFilePath := filepath.Join(task.Dir(), "packages.config")
	if checkPackagesIsExistence(packagesFilePath) {
		return scanPackage(task, packagesFilePath)
	}
	switch !task.IsNoBuild() {
	case true:
		if multipleBuilds(ctx, task) != nil {
			ui.Use(ctx).Display(ui.MsgWarn, "通过 Nuget获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")
			return noBuildEntrance(ctx, task)
		}
	default:
		return noBuildEntrance(ctx, task)
	}

	return nil
}

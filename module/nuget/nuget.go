package nuget

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/utils"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Nuget"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	// return utils.IsFile(filepath.Join(dir, "packages.config"))
	return utils.IsDir(dir)
}

func (Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx)
	task := model.UseInspectionTask(ctx)
	allowFallback := !env.ScannerScan
	doOld := false
	registeredAutoBuild := task.RegisterAutoBuild()

	var e error
	if !task.IsNoBuild() {
		if buildErr := multipleBuilds(ctx, task); buildErr != nil {
			registeredAutoBuild.MarkFailed()
			kind := scanerr.KindNugetFailed
			if errors.Is(buildErr, _ErrDotnetNotFound) {
				kind = scanerr.KindDotnetNotFound
			}
			scanerr.Add(ctx, scanerr.Param{
				Kind:    kind,
				Content: buildErr.Error(),
			})
			logger.Warn("multipleBuilds no build")
			ui.Use(ctx).Display(ui.MsgWarn, "通过 Nuget获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")
			if allowFallback {
				e = noBuildEntrance(ctx, task, &doOld)
			}
		} else if allowFallback {
			e = noBuildEntrance(ctx, task, &doOld)
		}
	} else {
		registeredAutoBuild.MarkDisabled()
		scanerr.Add(ctx, scanerr.Param{Kind: scanerr.KindBuildDisabled})
		logger.Warn("multipleBuilds no build")
		if allowFallback {
			e = noBuildEntrance(ctx, task, &doOld)
		}
	}
	if e != nil {
		logger.Sugar().Error(e)
	}
	if doOld && allowFallback {
		packagesFilePath := filepath.Join(task.Dir(), "packages.config")
		if checkPackagesIsExistence(packagesFilePath) {
			return scanPackage(task, packagesFilePath)
		}
	}

	return nil
}

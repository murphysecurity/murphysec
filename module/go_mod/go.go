package go_mod

import (
	"context"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "GoMod"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod")) || utils.IsFile(filepath.Join(dir, "Gopkg.lock"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	if utils.IsFile(filepath.Join(task.Dir(), "go.mod")) {
		// 新版本
		if task.IsNoBuild() {
			ui.Use(ctx).Display(ui.MsgWarn, "通过 go build获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")

			if err := baseScan(ctx); err != nil {
				return err
			}
		} else {
			if err := buildScan(ctx); err != nil {
				ui.Use(ctx).Display(ui.MsgWarn, "通过 go build获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")
				if err := baseScan(ctx); err != nil {
					return err
				}
			}
		}
	} else {
		// 旧版本Gopkg
		if err := parserGoPkgLock(ctx); err != nil {
			return nil
		}
	}

	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "go",
	Repository: "",
}

var _ErrGoNotFound = errors.New("go not found")

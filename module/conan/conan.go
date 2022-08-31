package conan

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (*Inspector) String() string {
	return "ConanInspector"
}

func (*Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "conanfile.txt")) ||
		utils.IsFile(filepath.Join(dir, "conanfile.py")) ||
		utils.IsFile(filepath.Join(dir, "conan.txt")) ||
		utils.IsFile(filepath.Join(dir, "conan.py"))
}
func (*Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	cmdInfo, e := getConanInfo(ctx)
	if e != nil {
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 conan 获取依赖信息失败，可能会导致检测结果不完整或失败，访问 https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", task.ScanDir))
		return e
	}
	jsonFilePath, e := ExecuteConanInfoCmd(ctx, cmdInfo.Path, task.ScanDir)

	var conanErr conanError
	if errors.As(e, &conanErr) {
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到您的环境中 conan 无法正常运行，可能会导致检测结果不完整或失败，访问 https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", task.ScanDir))
		for _, it := range conanErr.ErrorMultiLine() {
			task.UI().Display(display.MsgWarn, it)
		}
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

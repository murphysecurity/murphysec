package cocoapods

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

var Instance = &Inspector{}

func (i *Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Pod"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Podfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	projectDir := task.ScanDir
	podLockPath := filepath.Join(projectDir, "Podfile.lock")
	logger.Debug("Reading Podfile.lock", zap.String("path", podLockPath))
	data, e := os.ReadFile(filepath.Join(projectDir, "Podfile.lock"))
	if e != nil {
		return errors.WithMessage(e, "ReadPodLock")
	}
	tree, e := getDepFromLock(string(data))
	if e != nil {
		return errors.WithMessage(e, "Parse Podfile.lock failed")
	}
	task.AddModule(model.Module{
		PackageManager: model.PMCocoaPods,
		Language:       model.ObjectiveC,
		PackageFile:    "Podfile.lock",
		Name:           tree[0].Name,
		Dependencies:   tree,
	})
	return nil
}

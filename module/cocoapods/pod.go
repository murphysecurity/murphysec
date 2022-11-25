package cocoapods

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Pod"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Podfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	projectDir := task.Dir()
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
		PackageManager: "cocoapods",
		ModuleName:     tree[0].CompName,
		Dependencies:   tree,
		ModulePath:     podLockPath,
	})
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "cocoapods",
	Repository: "",
}

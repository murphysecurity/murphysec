package cargo

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (Inspector) String() string {
	return "Cargo"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "cargo.lock"))
}

func (Inspector) InspectProject(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("CargoInspector: %w", err)
		}
	}()
	task := model.UseInspectorTask(ctx)
	cargoLockPath := filepath.Join(task.ScanDir, "cargo.lock")
	data, e := os.ReadFile(cargoLockPath)
	if e != nil {
		return e
	}
	tree, e := analyzeCargoLock(data)
	if e != nil {
		return e
	}
	deps := tree.Dependencies
	task.AddModule(model.Module{
		PackageManager: model.PmCargo,
		Language:       model.Rust,
		Name:           tree.Name,
		Version:        tree.Version,
		RelativePath:   cargoLockPath,
		Dependencies:   deps,
	})
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

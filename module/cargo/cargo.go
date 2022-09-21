package cargo

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "Cargo"
}

func (i Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "cargo.lock"))
}

func (i Inspector) InspectProject(ctx context.Context) (err error) {
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

func (i Inspector) SupportFeature(feature base.Feature) bool {
	return base.FeatureAllowNested&feature > 0
}

var Instance base.Inspector = &Inspector{}

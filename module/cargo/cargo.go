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
	task := model.UseInspectionTask(ctx)
	cargoLockPath := filepath.Join(task.Dir(), "cargo.lock")
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
		PackageManager: "cargo",
		ModuleName:     tree.CompName,
		ModuleVersion:  tree.CompVersion,
		ModulePath:     cargoLockPath,
		Dependencies:   deps,
	})
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "cargo",
	Repository: "",
}

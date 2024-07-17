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

var _cargoLockNameList = []string{"Cargo.lock", "cargo.lock"}

func (Inspector) CheckDir(dir string) bool {
	for _, it := range _cargoLockNameList {
		if utils.IsFile(filepath.Join(dir, it)) {
			return true
		}
	}
	return false
}

func (Inspector) InspectProject(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("CargoInspector: %w", err)
		}
	}()
	task := model.UseInspectionTask(ctx)
	var cargoLockPath string
	var data []byte
	for _, it := range _cargoLockNameList {
		var e error
		cargoLockPath = filepath.Join(task.Dir(), it)
		data, e = os.ReadFile(cargoLockPath)
		if e == nil {
			break
		}
		if os.IsNotExist(e) {
			continue
		}
	}
	if data == nil {
		return fmt.Errorf("CargoInspector: Cargo.lock not found")
	}
	trees, e := analyzeCargoLock(data)
	if e != nil {
		return e
	}
	var m = model.Module{
		PackageManager: "cargo",
		ModulePath:     cargoLockPath,
	}
	if len(trees) == 1 {
		m.ModuleName = trees[0].CompName
		m.ModuleVersion = trees[0].CompVersion
		m.Dependencies = trees[0].Dependencies
	} else {
		for i := range trees {
			m.Dependencies = append(m.Dependencies, *trees[i])
		}
	}
	task.AddModule(m)
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "cargo",
	Repository: "",
}

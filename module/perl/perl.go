package perl

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "Perl"
}

func (i Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "MYMETA.json")) ||
		utils.IsFile(filepath.Join(dir, "META.json"))
}

func (i Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	var scanDir = task.Dir()
	var metaFile = filepath.Join(scanDir, "MYMETA.json")
	if !utils.IsFile(metaFile) {
		metaFile = filepath.Join(scanDir, "META.json")
	}
	if !utils.IsFile(metaFile) {
		return fmt.Errorf("no files valid")
	}
	data, e := os.ReadFile(metaFile)
	if e != nil {
		return e
	}
	deps, e := parseMeta(data)
	if e != nil {
		return e
	}
	task.AddModule(model.Module{
		PackageManager: "PerlEnv",
		ModuleName:     deps.Name,
		ModuleVersion:  deps.Version,
		ModulePath:     metaFile,
		Dependencies:   deps.deps(),
	})
	return nil
}

func (i Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "perl",
	Repository: "",
}

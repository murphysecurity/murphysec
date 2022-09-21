package perl

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
	return "Perl"
}

func (i Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "MYMETA.json")) ||
		utils.IsFile(filepath.Join(dir, "META.json"))
}

func (i Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	var scanDir = task.ScanDir
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
		Language:       "Perl",
		Name:           deps.Name,
		Version:        deps.Version,
		RelativePath:   metaFile,
		Dependencies:   deps.deps(),
	})
	return nil
}

func (i Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

var Instance base.Inspector = &Inspector{}

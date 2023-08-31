package pnpm

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	v5 "github.com/murphysecurity/murphysec/module/pnpm/v5"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
	"strings"
)

const LockfileName = "pnpm-lock.yaml"

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

type Inspector struct{}

func (Inspector) String() string {
	return "PNPMInspector"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, LockfileName)) && utils.IsFile(filepath.Join(dir, "package.json"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	inspectionTask := model.UseInspectionTask(ctx)
	dir := inspectionTask.Dir()
	lockfilePath := filepath.Join(dir, LockfileName)
	data, e := os.ReadFile(lockfilePath)
	if e != nil {
		return fmt.Errorf("PNPMInspector: read lockfile failed, %w", e)
	}
	version, e := parseLockfileVersion(data)
	if e != nil {
		return fmt.Errorf("PNPMInspector: parse lockfile version failed, %w", e)
	}
	versionNumber := matchLockfileVersion(version)
	var treeList []shared.DepTree
	if versionNumber == 5 {
		lockfile, e := v5.ParseLockfile(data)
		if e != nil {
			return fmt.Errorf("PNPMInspector(v5): %w", e)
		}
		treeList = v5.AnalyzeDepTree(lockfile)
	} else {
		return fmt.Errorf("PNPMInspector: unsupported version \"%s\"", version)
	}
	for _, tree := range treeList {
		var module = model.Module{
			ModulePath:     lockfilePath,
			PackageManager: "pnpm",
			Dependencies:   tree.Dependencies,
			ScanStrategy:   model.ScanStrategyNormal,
		}
		tree.Name = strings.TrimPrefix(tree.Name, "./")
		if tree.Name == "" || tree.Name == "." {
			module.ModulePath = lockfilePath
		} else {
			module.ModulePath = filepath.Join(dir, tree.Name, "<pnpm-module>")
		}
		inspectionTask.AddModule(module)
	}

	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var _ model.Inspector = (*Inspector)(nil)

package pnpm

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
)

const LockfileName = "pnpm-lock.yaml"

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

type _PnpmInspector struct{}

func (_PnpmInspector) String() string {
	return "PNPMInspector"
}

func (_PnpmInspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, LockfileName)) && utils.IsFile(filepath.Join(dir, "package.json"))
}

func (_PnpmInspector) InspectProject(ctx context.Context) error {
	inspectionTask := model.UseInspectionTask(ctx)
	dir := inspectionTask.Dir()
	lockfilePath := filepath.Join(dir, LockfileName)
	data, e := os.ReadFile(lockfilePath)
	if e != nil {
		return fmt.Errorf("PNPMInspector: read lockfile failed, %w", e)
	}
	lockfile, e := parseV6Lockfile(data, false)
	if e != nil {
		return e
	}
	deps, e := lockfile.buildDependencyTree(false)
	if e != nil {
		return e
	}
	inspectionTask.AddModule(model.Module{
		ModuleName:     dir,
		ModuleVersion:  "",
		ModulePath:     lockfilePath,
		PackageManager: "pnpm",
		Dependencies:   deps,
		ScanStrategy:   model.ScanStrategyNormal,
	})
	return nil
}

func (_PnpmInspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var _ model.Inspector = (*_PnpmInspector)(nil)

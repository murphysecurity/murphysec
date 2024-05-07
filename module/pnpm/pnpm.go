package pnpm

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func (Inspector) String() string {
	return "PNPMInspector"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, LockfileName)) && utils.IsFile(filepath.Join(dir, "package.json"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	inspectionTask := model.UseInspectorTask(ctx)
	dir := inspectionTask.ScanDir
	pResult := processDir(ctx, dir)
	if pResult.e != nil {
		return fmt.Errorf("PNPMInspector: %w", pResult.e)
	}
	for _, tree := range pResult.trees {
		var module = model.Module{
			Name:           "<pnpm-root-module>",
			RelativePath:   pResult.lockfile,
			PackageManager: "pnpm",
			Dependencies:   tree.Dependencies,
			ScanStrategy:   model.ScanStrategyNormal,
		}
		tree.Name = strings.TrimPrefix(tree.Name, "./")
		if tree.Name == "" || tree.Name == "." {
			module.RelativePath = pResult.lockfile
		} else {
			if len(tree.Dependencies) == 0 {
				continue
			}
			module.Name = fmt.Sprintf("<pnpm-module>/%s", tree.Name)
			module.RelativePath = filepath.Join(dir, tree.Name, "<pnpm-module>")
		}
		inspectionTask.AddModule(module)
	}

	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

var _ model.Inspector = (*Inspector)(nil)

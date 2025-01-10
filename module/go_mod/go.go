package go_mod

import (
	"context"
	"path/filepath"

	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "GoMod"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	if task.IsNoBuild() {
		if err := baseScan(ctx); err != nil {
			return err
		}
	} else {
		if err := buildScan(ctx); err != nil {
			if err := baseScan(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "go",
	Repository: "",
}

var _ErrGoNotFound = errors.New("go not found")

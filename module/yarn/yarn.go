package yarn

import (
	"context"
	"os"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
)

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

type Dep struct {
	Name     string
	Version  string
	Children []Dep
}

func mapToModel(deps []Dep) []model.DependencyItem {
	r := _mapToModel(deps)
	for i := range r {
		r[i].DependencyRelation = model.DependencyRelationDirect
		r[i].IsOnline.SetOnline(false)
	}
	return r
}
func _mapToModel(deps []Dep) []model.DependencyItem {
	var r = make([]model.DependencyItem, len(deps))
	for i := range deps {
		r[i] = model.DependencyItem{
			Component: model.Component{
				CompName:    deps[i].Name,
				CompVersion: deps[i].Version,
				EcoRepo:     EcoRepo,
			},
			Dependencies:       _mapToModel(deps[i].Children),
			DependencyRelation: model.DependencyRelationTransitive,
		}
	}
	return r
}

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "Yarn"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	info, e := os.Stat(filepath.Join(dir, "yarn.lock"))
	return e == nil && !info.IsDir()
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	var logger = logctx.Use(ctx).Sugar()
	dir := task.Dir()
	logger.Info("yarn inspect.", dir)
	rs, e := analyzeYarnDep(ctx, dir)

	if e != nil {
		return e
	}
	m := model.Module{
		PackageManager: "yarn",
		ModuleName:     filepath.Base(dir),
		ModulePath:     filepath.Join(dir, "yarn.lock"),
		Dependencies:   mapToModel(rs),
	}
	if n, v := readModuleName(dir); n != "" {
		m.ModuleName = n
		m.ModuleVersion = v
	}
	task.AddModule(m)
	return nil
}

package rebar3

import (
	"context"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Rebar3"
}
func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "rebar.config"))
}
func (Inspector) InspectProject(ctx context.Context) error {
	if env.DoNotBuild {
		return nil
	}
	task := model.UseInspectionTask(ctx)
	_, e := GetRebar3Version(ctx)
	if e != nil {
		return e
	}
	tree, e := EvaluateRebar3Tree(ctx, task.Dir())
	if e != nil {
		return e
	}
	if len(tree) == 0 {
		return nil
	}

	task.AddModule(model.Module{
		PackageManager: "rebar3",
		ModuleName:     tree[0].Name,
		ModuleVersion:  tree[0].Version,
		ModulePath:     filepath.Join(task.Dir(), "rebar.config"),
		Dependencies:   _mapDepNodes(tree),
	})
	return nil
}

func _mapDepNodes(node []depNode) (r []model.DependencyItem) {
	for _, it := range node {
		var di model.DependencyItem
		di.CompName = it.Name
		di.CompVersion = it.Version
		di.EcoRepo = EcoRepo
		r = append(r, di)
	}
	return
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "rebar",
	Repository: "",
}

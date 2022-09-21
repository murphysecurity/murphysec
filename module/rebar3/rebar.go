package rebar3

import (
	"context"
	"github.com/google/uuid"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
)

type Inspector struct{}

var Instance = &Inspector{}

func (Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (Inspector) String() string {
	return "Rebar3"
}
func (Inspector) CheckDir(dir string) bool {
	if utils.IsFile(filepath.Join(dir, "rebar.config")) {
		return true
	}
	return false
}
func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	ver, e := GetRebar3Version(ctx)
	if e != nil {
		return e
	}
	tree, e := EvaluateRebar3Tree(ctx, task.ScanDir)
	if e != nil {
		return e
	}
	if len(tree) == 0 {
		return nil
	}

	task.AddModule(model.Module{
		PackageManager: model.PmRebar3,
		Language:       model.Erlang,
		Name:           tree[0].Name,
		Version:        tree[0].Version,
		RelativePath:   filepath.Join(task.ScanDir, "rebar.config"),
		Dependencies:   _mapDepNodes(tree),
		RuntimeInfo:    ver,
		UUID:           uuid.New(),
	})
	return nil
}

func (Inspector) PackageManagerType() model.PackageManagerType {
	return model.PmRebar3
}

func _mapDepNodes(node []depNode) (r []model.Dependency) {
	for _, it := range node {
		r = append(r, model.Dependency{
			Name:         it.Name,
			Version:      it.Version,
			Dependencies: _mapDepNodes(it.Children),
		})
	}
	return
}

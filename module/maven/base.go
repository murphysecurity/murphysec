package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/vifraa/gopom"
	"path/filepath"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "MavenInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pom.xml"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	modules, e := ScanMavenProject(ctx, task)
	if e != nil {
		return e
	}
	for _, it := range modules {
		task.AddModule(it)
	}
	return nil
}

type Repo interface {
	Fetch(coordinate Coordinate) (*gopom.Project, error)
}

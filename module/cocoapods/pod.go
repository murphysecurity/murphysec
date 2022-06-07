package cocoapods

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "PodInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Podfile.lock"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	projectDir := task.ScanDir
	data, e := os.ReadFile(filepath.Join(projectDir, "Podfile.lock"))
	if e != nil {
		return errors.Wrap(e, "ReadPodLock")
	}
	tree, e := getDepFromLock(string(data))
	if e != nil {
		return errors.Wrap(e, "ParsePodLock")
	}
	task.AddModule(model.Module{
		PackageManager: model.PMCocoaPods,
		Language:       model.ObjectiveC,
		PackageFile:    "Podfile.lock",
		Name:           tree[0].Name,
		Dependencies:   tree,
	})
	return nil
}

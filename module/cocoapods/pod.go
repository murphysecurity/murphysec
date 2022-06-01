package cocoapods

import (
	"github.com/pkg/errors"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMCocoaPods
}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "PodInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "Podfile.lock"))
}
func (i *Inspector) Inspect(task *base.ScanTask) ([]base.Module, error) {
	data, e := os.ReadFile(filepath.Join(task.ProjectDir, "Podfile.lock"))
	if e != nil {
		return nil, errors.Wrap(e, "Open Gemfile.lock failed")
	}
	tree, e := getDepFromLock(string(data))
	if e != nil {
		return nil, errors.Wrap(e, "Bundler")
	}
	return []base.Module{{
		PackageManager: "cocoapods",
		Language:       "Objective-C",
		PackageFile:    "Podfile.lock",
		Name:           tree[0].Name,
		Dependencies:   tree,
	}}, nil
}

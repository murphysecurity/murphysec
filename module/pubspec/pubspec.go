package pubspec

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "Pubspec"
}

func (i Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, lockfileName))
}

func (i Inspector) InspectProject(ctx context.Context) (e error) {
	var logger = logctx.Use(ctx)
	var inspTask = model.UseInspectionTask(ctx)
	lockPath := filepath.Join(inspTask.Dir(), lockfileName)
	var f *os.File
	f, e = os.Open(lockPath)
	if e != nil {
		return
	}
	defer func() { utils.CloseLogErrZap(f, logger) }()
	var deps []model.DependencyItem
	deps, e = parseFile(ctx, f)
	if e != nil {
		return
	}
	var module = model.Module{
		ModulePath:     lockPath,
		PackageManager: "pubspec",
		Dependencies:   deps,
		ScanStrategy:   model.ScanStrategyNormal,
	}
	inspTask.AddModule(module)
	return
}

func parseFile(ctx context.Context, reader io.Reader) (r []model.DependencyItem, e error) {
	var parsedLockfile *pubspecLockfile
	e = yaml.NewDecoder(bufio.NewReader(reader)).Decode(&parsedLockfile)
	if e != nil {
		return
	}
	r = fp.Map(func(it lo.Entry[string, pubspecLockPackage]) model.DependencyItem {
		dep := model.DependencyItem{
			Component: model.Component{
				CompName:    it.Key,
				CompVersion: it.Value.Version,
				EcoRepo: model.EcoRepo{
					Ecosystem: "pubspec",
				},
			},
			IsOnline:           model.IsOnlineTrue(),
			DependencyRelation: model.DependencyRelationTransitive,
		}
		for _, flag := range strings.Split(it.Value.Dependency, " ") {
			switch flag {
			case "dev":
				dep.IsOnline.SetOnline(false)
			case "direct":
				dep.DependencyRelation = model.DependencyRelationDirect
			}
		}
		return dep
	})(lo.Entries(parsedLockfile.Packages))
	return
}

func (i Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var _ model.Inspector = &Inspector{}

const lockfileName = "pubspec.lock"

type pubspecLockPackage struct {
	Version    string `yaml:"version"`
	Dependency string `yaml:"dependency"`
}
type pubspecLockfile struct {
	Packages map[string]pubspecLockPackage `yaml:"packages,omitempty"`
}

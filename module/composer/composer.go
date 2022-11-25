package composer

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io/fs"
	"path/filepath"
	"strings"
)

const _ComposerManifestFileSizeLimit = 4 * 1024 * 1024 // 4MiB
const _ComposerLockFileSizeLimit = _ComposerManifestFileSizeLimit

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Composer"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "composer.json"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx)
	task := model.UseInspectionTask(ctx)
	dir := task.Dir()
	manifest, e := readManifest(ctx, filepath.Join(dir, "composer.json"))
	if e != nil {
		return e
	}
	module := &model.Module{
		PackageManager: "composer",
		ModuleName:     manifest.Name,
		ModuleVersion:  manifest.Version,
		ModulePath:     filepath.Join(dir, "composer.json"),
	}
	lockfilePkgs := map[string]Package{}

	{
		if !utils.IsPathExist(filepath.Join(dir, "composer.lock")) {
			logger.Info("composer.lock doesn't exists. Try to generate it")
			if e := doComposerInstall(context.TODO(), dir); e != nil {
				logger.Sugar().Warnf("Do composer install fail. %s", e.Error())
			} else {
				logger.Sugar().Info("Do composer install succeeded")
			}
		}
		composerLockFilePath := filepath.Join(dir, "composer.lock")
		logger.Debug("Reading composer.lock", zap.String("path", composerLockFilePath))
		pkgs, e := readComposerLockFile(composerLockFilePath)
		if e != nil {
			logger.Sugar().Infof("Read composer lock file failed: %s", e.Error())
		}
		pkgs = append(pkgs, vendorScan(ctx, filepath.Join(dir, "vendor"))...)
		for _, it := range pkgs {
			lockfilePkgs[it.Name] = it
		}
	}

	for _, requiredPkg := range manifest.Require {
		node := _buildDepTree(lockfilePkgs, map[string]struct{}{}, requiredPkg.Name, requiredPkg.Version)
		if node != nil {
			module.Dependencies = append(module.Dependencies, *node)
		}
	}
	if module.IsZero() {
		return nil
	}
	task.AddModule(*module)
	return nil
}

func _buildDepTree(lockfile map[string]Package, visitedDep map[string]struct{}, targetName string, versionConstraint string) *model.DependencyItem {
	if _, ok := visitedDep[targetName]; ok || len(visitedDep) > 3 {
		return nil
	}
	visitedDep[targetName] = struct{}{}
	defer delete(visitedDep, targetName)
	rs := &model.DependencyItem{
		Component: model.Component{
			CompName:    targetName,
			CompVersion: versionConstraint,
			EcoRepo:     EcoRepo,
		},
	}
	pkg := lockfile[rs.CompName]
	if targetName == "php" || (strings.HasPrefix(targetName, "ext-") && (pkg.Version == "*" || pkg.Version == "" || versionConstraint == "*")) {
		return nil
	}
	if pkg.Version == "" {
		return rs // fallback
	}
	rs.CompVersion = pkg.Version
	for _, requiredPkgName := range pkg.Require {
		node := _buildDepTree(lockfile, visitedDep, requiredPkgName, "") // ignore transitive dependency version constraint
		if node != nil {
			rs.Dependencies = append(rs.Dependencies, *node)
		}
	}
	return rs
}

type Element struct {
	Name    string
	Version string
}

type Package struct {
	Element
	Require []string
}

type Manifest struct {
	Element
	Require []Element
}

func vendorScan(ctx context.Context, dir string) []Package {
	logger := logctx.Use(ctx)
	logger.Debug("vendorScan", zap.String("dir", dir))
	defer logger.Debug("vendorScan terminated")
	var rs []Package
	e := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.Name() == "composer.json" {
			m, e := readManifest(ctx, path)
			if e != nil {
				return nil
			}
			var p Package
			p.Name = m.Name
			p.Version = m.Version
			for _, it := range m.Require {
				p.Require = append(p.Require, it.Name)
			}
			rs = append(rs, p)
		}
		return nil
	})
	if e != nil {
		logger.Sugar().Warnf("Walk: %v", e)
	}
	return rs
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "composer",
	Repository: "",
}

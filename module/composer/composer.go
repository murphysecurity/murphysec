package composer

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

const _ComposerManifestFileSizeLimit = 4 * 1024 * 1024 // 4MiB
const _ComposerLockFileSizeLimit = _ComposerManifestFileSizeLimit

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Composer"
}

func (Inspector) CheckDir(ctx context.Context, dir string) bool {
	return utils.IsFile(filepath.Join(dir, "composer.json")) ||
		utils.IsFile(filepath.Join(dir, "composer.lock")) ||
		utils.IsFile(filepath.Join(dir, "installed.json")) ||
		utils.IsFile(filepath.Join(dir, "vendor", "composer", "installed.json"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx)
	task := model.UseInspectionTask(ctx)
	dir := task.Dir()
	manifestPath := filepath.Join(dir, "composer.json")
	modulePath := manifestPath
	manifest := &Manifest{}
	var e error
	if utils.IsFile(manifestPath) {
		manifest, e = readManifest(ctx, manifestPath)
		if e != nil {
			return e
		}
	} else {
		logger.Sugar().Infof("composer.json not found, fallback to installed/lock files. dir=%s", dir)
		modulePath = filepath.Join(dir, "installed.json")
		if !utils.IsFile(modulePath) {
			modulePath = filepath.Join(dir, "vendor", "composer", "installed.json")
			if !utils.IsFile(modulePath) {
				modulePath = filepath.Join(dir, "composer.lock")
			}
		}
	}
	module := &model.Module{
		PackageManager: "composer",
		ModuleName:     manifest.Name,
		ModuleVersion:  manifest.Version,
		ModulePath:     modulePath,
	}
	if module.ModuleName == "" {
		module.ModuleName = filepath.Base(dir)
	}
	lockfilePkgs := map[string]Package{}

	{
		if utils.IsFile(manifestPath) && !utils.IsPathExist(filepath.Join(dir, "composer.lock")) {
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
		installedPaths := []string{
			filepath.Join(dir, "installed.json"),
			filepath.Join(dir, "vendor", "composer", "installed.json"),
		}
		for _, installedPath := range installedPaths {
			installedPkgs, ie := readComposerInstalledFile(installedPath)
			if ie != nil {
				logger.Sugar().Debugf("Read installed.json failed: %s", ie.Error())
				continue
			}
			pkgs = append(pkgs, installedPkgs...)
		}
		pkgs = append(pkgs, vendorScan(ctx, filepath.Join(dir, "vendor"))...)
		for _, it := range pkgs {
			if it.Version == "" || isVersionConstrain(it.Version) {
				continue
			}
			lockfilePkgs[it.Name] = it
		}
	}

	roots := map[string]string{}
	for _, requiredPkg := range manifest.Require {
		roots[requiredPkg.Name] = requiredPkg.Version
	}
	if len(roots) == 0 {
		for name := range lockfilePkgs {
			roots[name] = ""
		}
	}
	var rootNames []string
	for name := range roots {
		rootNames = append(rootNames, name)
	}
	sort.Strings(rootNames)
	for _, name := range rootNames {
		node := _buildDepTree(lockfilePkgs, map[string]struct{}{}, name, roots[name])
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

func isVersionConstrain(v string) bool {
	if v == "" {
		return false
	}
	if v[0] == '^' || v[0] == '~' || v[0] == '>' || v[0] == '<' || v[0] == '=' {
		return true
	}
	if strings.Contains(v, "||") || strings.Contains(v, ",") {
		return true
	}
	return false
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

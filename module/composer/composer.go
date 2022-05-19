package composer

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/fs"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/simplejson"
	"os/exec"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func (i *Inspector) String() string {
	return "ComposerInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "composer.json"))
}

func (i *Inspector) Inspect(task *base.ScanTask) ([]base.Module, error) {
	dir := task.ProjectDir
	manifest, e := readManifest(filepath.Join(dir, "composer.json"))
	if e != nil {
		return nil, e
	}
	module := &base.Module{
		PackageManager: "composer",
		Language:       "php",
		PackageFile:    "composer.json",
		Name:           manifest.Name,
		Version:        manifest.Version,
		FilePath:       filepath.Join(dir, "composer.json"),
		Dependencies:   []base.Dependency{},
		RuntimeInfo:    nil,
		UUID:           uuid.UUID{},
	}
	lockfilePkgs := map[string]Package{}

	{
		if !utils.IsPathExist(filepath.Join(dir, "composer.lock")) {
			logger.Info.Println("composer.lock doesn't exists. Try generate it")
			c := exec.Command("composer", "update", "--ignore-platform-req=*", "--no-dev", "--no-progress")
			logger.Info.Println("Command:", c.String())
			_, e := c.Output()
			if e != nil {
				logger.Info.Println("composer update exit with no errors")
			} else {
				logger.Warn.Warn("composer update exit with error:" + e.Error())
			}
		}
		pkgs, e := readComposerLockFile(filepath.Join(dir, "composer.lock"))
		if e != nil {
			logger.Info.Println("Composer:", e.Error())
		}
		pkgs = append(pkgs, vendorScan(filepath.Join(dir, "vendor"))...)
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
	return []base.Module{*module}, nil
}

func _buildDepTree(lockfile map[string]Package, visitedDep map[string]struct{}, targetName string, versionConstraint string) *base.Dependency {
	if _, ok := visitedDep[targetName]; ok || len(visitedDep) > 3 {
		return nil
	}
	visitedDep[targetName] = struct{}{}
	defer delete(visitedDep, targetName)
	rs := &base.Dependency{
		Name:    targetName,
		Version: versionConstraint,
	}
	pkg := lockfile[rs.Name]
	if targetName == "php" || (strings.HasPrefix(targetName, "ext-") && (pkg.Version == "*" || pkg.Version == "" || versionConstraint == "*")) {
		return nil
	}
	if pkg.Version == "" {
		return rs // fallback
	}
	rs.Version = pkg.Version
	for _, requiredPkgName := range pkg.Require {
		node := _buildDepTree(lockfile, visitedDep, requiredPkgName, "") // ignore transitive dependency version constraint
		if node != nil {
			rs.Dependencies = append(rs.Dependencies, *node)
		}
	}
	return rs
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMComposer
}

func New() base.Inspector {
	return &Inspector{}
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

func readComposerLockFile(path string) ([]Package, error) {
	lockFileData, e := utils.ReadFileLimited(path, 4*1024*1024)
	if e != nil {
		return nil, errors.Wrap(e, "Read composer.lock failed")
	}
	pkgs, e := parseComposerLock(lockFileData)
	if e != nil {
		return nil, errors.Wrap(e, "Parse composer.lock failed")
	}
	return pkgs, nil
}

func parseComposerLock(data []byte) ([]Package, error) {
	var j simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return nil, errors.Wrap(e, "ParseComposerLock:")
	}
	pkgList := make([]Package, 0)
	for _, pkg := range j.Get("packages").JSONArray() {
		p := Package{}
		p.Name = pkg.Get("name").String()
		p.Version = pkg.Get("version").String()
		if p.Name == "" || p.Version == "" {
			continue
		}
		for s := range pkg.Get("require").JSONMap() {
			p.Require = append(p.Require, s)
		}
		pkgList = append(pkgList, p)
	}
	return pkgList, nil
}

func readManifest(path string) (*Manifest, error) {
	composerFileData, e := utils.ReadFileLimited(path, 4*1024*1024)
	if e != nil {
		return nil, errors.Wrap(e, "Read composer.json failed.")
	}
	manifest, e := parseComposeManifest(composerFileData)
	if e != nil {
		return nil, errors.Wrap(e, "Parse composer.json failed.")
	}
	return manifest, nil
}

func parseComposeManifest(data []byte) (*Manifest, error) {
	var j simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return nil, errors.Wrap(e, "ParseComposeManifest:")
	}
	m := &Manifest{}
	m.Name = j.Get("name").String()
	m.Version = j.Get("version").String()
	for name, versionConstraint := range j.Get("require").JSONMap() {
		m.Require = append(m.Require, Element{
			Name:    name,
			Version: versionConstraint.String(),
		})
	}
	return m, nil
}

func vendorScan(dir string) []Package {
	var rs []Package
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.Name() == "composer.json" {
			m, e := readManifest(path)
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
	return rs
}

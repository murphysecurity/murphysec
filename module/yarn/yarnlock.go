package yarn

import (
	"github.com/iseki0/go-yarnlock"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"os"
	"path/filepath"
	"regexp"
)

type Inspector struct{}

func (i *Inspector) String() string {
	return "YarnInspector@v0.0.1"
}

func (i *Inspector) Version() string {
	return "0.0.1"
}

func (i *Inspector) CheckDir(dir string) bool {
	info, e := os.Stat(filepath.Join(dir, "yarn.lock"))
	return e == nil && !info.IsDir()
}

func (i *Inspector) Inspect(dir string) ([]base.Module, error) {
	logger.Info.Println("yarn inspect.", dir)
	rs, e := analyzeYarnDep(dir)
	if e != nil {
		return nil, e
	}
	m := base.Module{
		PackageManager: "yarn",
		Language:       "javascript",
		PackageFile:    "yarn.lock",
		Name:           filepath.Base(dir),
		Version:        "",
		RelativePath:   "",
		Dependencies:   rs,
	}
	return []base.Module{m}, nil
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMYarn
}

func New() base.Inspector {
	return &Inspector{}
}

func analyzeYarnDep(dir string) ([]base.Dependency, error) {
	f, e := os.Open(filepath.Join(dir, "yarn.lock"))
	if e != nil {
		return nil, errors.Wrap(e, "Open yarn.lock failed.")
	}
	defer f.Close()
	data, e := io.ReadAll(io.LimitReader(f, 16*1024*1024))
	if e != nil {
		return nil, errors.Wrap(e, "Read yarn.lock failed.")
	}
	lockfile, e := yarnlock.ParseLockFileData(data)
	if e != nil {
		return nil, errors.Wrap(e, "Parse lockfile failed.")
	}
	return buildDepTree(lockfile), nil
}

func buildDepTree(lkFile yarnlock.LockFile) []base.Dependency {
	type id struct {
		name    string
		version string
	}
	var rs []base.Dependency
	repeatedElement := map[id]struct{}{}
	for _, key := range lkFile.RootElement() {
		node := _buildDepTree(lkFile, key, map[string]struct{}{})
		if node == nil {
			continue
		}
		if _, ok := repeatedElement[id{node.Name, node.Version}]; ok {
			continue
		}
		rs = append(rs, *node)
	}
	return rs
}

func _buildDepTree(lkFile yarnlock.LockFile, element string, visitedKey map[string]struct{}) *base.Dependency {
	{
		// avoid circle dependency
		if _, ok := visitedKey[element]; ok {
			return nil
		}
		visitedKey[element] = struct{}{}
		defer delete(visitedKey, element)
	}
	info, ok := lkFile[element]
	if !ok {
		return nil
	}
	pkgName, pkgVer := parsePkgName(element)
	if pkgName == "" || pkgVer == "" {
		return nil
	}
	node := &base.Dependency{
		Name:         pkgName,
		Version:      info.Version, // use real version
		Dependencies: nil,
	}
	type id struct {
		name    string
		version string
	}
	repeatedElement := map[id]struct{}{}
	for childComp, childVer := range lkFile[element].Dependencies {
		childKey := childComp + "@" + childVer
		c := _buildDepTree(lkFile, childKey, visitedKey)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Dependencies = append(node.Dependencies, *c)
	}
	return node
}

var __parsePkgNamePattern = regexp.MustCompile("(@?[^@]+)@(.+)")

func parsePkgName(input string) (pkgName string, pkgVersion string) {
	m := __parsePkgNamePattern.FindStringSubmatch(input)
	if m == nil {
		return "", ""
	} else {
		return m[1], m[2]
	}
}

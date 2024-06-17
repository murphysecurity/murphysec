package yarn

import (
	"context"
	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/module/pkgjs"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func readModuleName(dir string) (string, string) {
	f, e := pkgjs.ReadDir(dir)
	if e != nil {
		return "", ""
	}
	return f.Name, f.Version
}

func yarnFallback(dir string) ([]Dep, error) {
	var rs []Dep
	pkg0, e := pkgjs.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	distinct := map[string]string{}
	for k, v := range pkg0.DevDependencies {
		distinct[k] = v
	}
	for k, v := range pkg0.Dependencies {
		distinct[k] = v
	}

	for k, v := range distinct {
		var di Dep
		di.Name = k
		di.Version = v
		rs = append(rs, di)
	}
	return rs, nil
}

func analyzeYarnDep(ctx context.Context, dir string) ([]Dep, error) {
	var logger = logctx.Use(ctx).Sugar()
	f, e := os.Open(filepath.Join(dir, "yarn.lock"))
	if e != nil {
		logger.Infof("Open yarn.lock failed. %v", e)
		return yarnFallback(dir)
	}
	defer func() { _ = f.Close() }()
	data, e := io.ReadAll(io.LimitReader(f, 16*1024*1024))
	if e != nil {
		return nil, errors.Wrap(e, "Read yarn.lock failed.")
	}
	lockfile, e := yarnlock.ParseLockFileData(data)
	if e != nil {
		return nil, errors.Wrap(e, "Parse lockfile failed.")
	}
	pkg, e := pkgjs.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	return buildDepTree(lockfile, pkg), nil
}

func buildDepTree(lkFile yarnlock.LockFile, pkg *pkgjs.Pkg) []Dep {
	type id struct {
		name    string
		version string
	}
	var rs []Dep
	repeatedElement := map[id]struct{}{}
	for n, v := range pkg.Dependencies {
		node := _buildDepTree(lkFile, n+"@"+v, map[string]struct{}{}, 5)
		if node == nil {
			continue
		}
		if _, ok := repeatedElement[id{node.Name, node.Version}]; ok {
			continue
		}
		rs = append(rs, *node)
	}
	for n, v := range pkg.DevDependencies {
		node := _buildDepTree(lkFile, n+"@"+v, map[string]struct{}{}, 5)
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

var versionNpmPattern = regexp.MustCompile(`^npm:(.+?)@(.+)`)

func _buildDepTree(lkFile yarnlock.LockFile, element string, visitedKey map[string]struct{}, depth int) *Dep {
	if depth < 0 {
		return nil
	}
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
	node := &Dep{
		Name:    pkgName,
		Version: info.Version, // use real version
	}
	type id struct {
		name    string
		version string
	}
	repeatedElement := map[id]struct{}{}
	for childComp, childVer := range lkFile[element].Dependencies {
		childKey := childComp + "@" + childVer
		c := _buildDepTree(lkFile, childKey, visitedKey, depth-1)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Children = append(node.Children, *c)
	}
	for childComp, childVer := range lkFile[element].OptionalDependencies {
		childKey := childComp + "@" + childVer
		c := _buildDepTree(lkFile, childKey, visitedKey, depth-1)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Children = append(node.Children, *c)
	}
	return node
}

var __parsePkgNamePattern = regexp.MustCompile("(@?[^@]+)@(.+)")

func parsePkgName(input string) (pkgName string, pkgVersion string) {
	m := __parsePkgNamePattern.FindStringSubmatch(input)
	if m == nil {
		return "", ""
	} else {
		if m := versionNpmPattern.FindStringSubmatch(m[2]); m != nil {
			return m[1], m[2]
		}
		return m[1], m[2]
	}
}

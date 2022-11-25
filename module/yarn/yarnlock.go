package yarn

import (
	"context"
	"encoding/json"
	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/simplejson"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type pkgFile struct {
	DevDependencies map[string]string `json:"dev_dependencies,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
}

func readModuleName(dir string) (string, string) {
	f, e := os.Open(filepath.Join(dir, "package.json"))
	if e == nil {
		return "", ""
	}
	defer f.Close()
	r := io.LimitReader(f, 1024*1024)
	data, e := io.ReadAll(r)
	if e != nil {
		return "", ""
	}
	var j *simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return "", ""
	}
	return j.Get("name").String(), j.Get("version").String()
}

func yarnFallback(dir string) ([]Dep, error) {
	f, e := os.Open(filepath.Join(dir, "package.json"))
	if e != nil {
		return nil, errors.Wrap(e, "Open package.json failed.")
	}
	defer f.Close()
	r := io.LimitReader(f, 1024*1024)
	data, e := io.ReadAll(r)
	if e != nil {
		return nil, e
	}
	var pkg pkgFile
	if e := json.Unmarshal(data, &pkg); e != nil {
		return nil, errors.Wrap(e, "parse failed")
	}
	var rs []Dep
	distinct := map[string]string{}
	for k, v := range pkg.DevDependencies {
		distinct[k] = v
	}
	for k, v := range pkg.Dependencies {
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

func buildDepTree(lkFile yarnlock.LockFile) []Dep {
	type id struct {
		name    string
		version string
	}
	var rs []Dep
	repeatedElement := map[id]struct{}{}
	for _, key := range lkFile.RootElement() {
		node := _buildDepTree(lkFile, key, map[string]struct{}{}, 5)
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

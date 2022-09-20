package yarn

import (
	"context"
	"encoding/json"
	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils/simplejson"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type Inspector struct{}

var Instance = &Inspector{}

func (i *Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Yarn"
}

func (i *Inspector) CheckDir(dir string) bool {
	info, e := os.Stat(filepath.Join(dir, "yarn.lock"))
	return e == nil && !info.IsDir()
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	dir := task.ScanDir
	logger.Info.Println("yarn inspect.", dir)
	rs, e := analyzeYarnDep(dir)

	if e != nil {
		return e
	}
	m := model.Module{
		PackageManager: model.PMYarn,
		Language:       model.JavaScript,
		Name:           filepath.Base(dir),
		Version:        "",
		FilePath:       filepath.Join(dir, "yarn.lock"),
		Dependencies:   rs,
	}
	if n, v := readModuleName(dir); n != "" {
		m.Name = n
		m.Version = v
	}
	task.AddModule(m)
	return nil
}

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

func yarnFallback(dir string) ([]model.Dependency, error) {
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
	rs := make([]model.Dependency, 0)
	distinct := map[string]string{}
	for k, v := range pkg.DevDependencies {
		distinct[k] = v
	}
	for k, v := range pkg.Dependencies {
		distinct[k] = v
	}
	for k, v := range distinct {
		rs = append(rs, model.Dependency{
			Name:    k,
			Version: v,
		})
	}
	return rs, nil
}

func analyzeYarnDep(dir string) ([]model.Dependency, error) {
	f, e := os.Open(filepath.Join(dir, "yarn.lock"))
	if e != nil {
		logger.Info.Println("Open yarn.lock failed.", e.Error())
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

func buildDepTree(lkFile yarnlock.LockFile) []model.Dependency {
	type id struct {
		name    string
		version string
	}
	var rs []model.Dependency
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

func _buildDepTree(lkFile yarnlock.LockFile, element string, visitedKey map[string]struct{}, depth int) *model.Dependency {
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
	node := &model.Dependency{
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
		c := _buildDepTree(lkFile, childKey, visitedKey, depth-1)
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

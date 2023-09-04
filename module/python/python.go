package python

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/pathignore"
	"github.com/murphysecurity/murphysec/model"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Inspector struct{}

func (i Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i Inspector) String() string {
	return "Python"
}

func (i Inspector) CheckDir(dir string) bool {
	r, e := os.ReadDir(dir)
	if e == nil {
		for _, it := range r {
			if it.IsDir() {
				continue
			}
			name := it.Name()
			if name == "conanfile.py" {
				continue
			}
			if isRequirementsFile(name) {
				return true
			}
			if filepath.Ext(name) == ".py" {
				return true
			}
		}
	}
	return false
}

func (i Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx).Sugar()
	dir := model.UseInspectionTask(ctx).Dir()
	info, e := collectDepsInfo(ctx, dir)
	if e != nil {
		return e
	}
	if len(info) == 0 {
		logger.Infof("found no deps, omit module")
		return nil
	}

	m := model.Module{
		ModuleName:     filepath.ToSlash(model.UseInspectionTask(ctx).RelDir()),
		PackageManager: "pip",
		ModulePath:     dir,
	}
	if m.ModuleName == "." {
		m.ModuleName = "Python"
	}
	for _, it := range info {
		k, v := it[0], it[1]
		var di model.DependencyItem
		di.CompName = k
		di.CompVersion = v
		di.EcoRepo = EcoRepo
		m.Dependencies = append(m.Dependencies, di)
	}
	model.UseInspectionTask(ctx).AddModule(m)
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "pip",
	Repository: "",
}

var requirementPattern = regexp.MustCompile(`^requirement.*\.txt$`)

func isRequirementsFile(filename string) bool {
	return requirementPattern.MatchString(strings.ToLower(filename))
}

func isDockerfile(filename string) bool {
	return strings.Contains(strings.ToLower(filename), "dockerfile")
}

func dirIgnore(name string) bool {
	return name == "" || name[0] == '.' || pathignore.DirName(name)
}

func collectDepsInfo(ctx context.Context, dir string) ([][2]string, error) {
	var logger = logctx.Use(ctx).Sugar()
	if !filepath.IsAbs(dir) {
		panic("dir must be absolute")
	}
	var noVersionComps = make(map[string]struct{})
	var versionedComps = make(map[string]string)
	var unknownVersionComps = make(map[string]struct{})
	e := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return err
		}
		filename := d.Name()
		if d.IsDir() {
			if dirIgnore(filename) {
				return filepath.SkipDir
			}
			return nil
		}
		// ignore conan.py
		if filename == "" || filename[0] == '.' || filename == "conan.py" {
			return nil
		}
		if isDockerfile(filename) {
			data, e := readFile(path, 64*1024) // Dockerfile max 64K
			if e != nil {
				logger.Warnf("read dockerfile: %s %v", path, e)
				return nil
			}
			for _, s := range parseDockerfilePipInstall(string(data)) {
				noVersionComps[s] = struct{}{}
			}
			return nil
		}
		if isRequirementsFile(filename) {
			data, e := readFile(path, 64*1024)
			if e != nil {
				logger.Warnf("read requirement: %s %v", path, e)
				return nil
			}
			for k, v := range parseRequirements(string(data)) {
				if v == "" {
					unknownVersionComps[k] = struct{}{}
				} else {
					versionedComps[k] = v
				}
			}
			return nil
		}
		if filepath.Ext(filename) == ".py" {
			data, e := readFile(path, 256*1024)
			if e != nil {
				logger.Warnf("read py: %s %v", path, e)
				return nil
			}
			for _, s := range parsePyImport(string(data)) {
				if pyPkgBlackList[s] {
					continue
				}
				unknownVersionComps[s] = struct{}{}
			}
		}
		return nil
	})
	if e != nil {
		logger.Warnf("walk error: %v", e)
	}
	// merge unknown version components
	for s := range unknownVersionComps {
		if versionedComps[s] != "" {
			delete(unknownVersionComps, s)
		}
	}
	if len(unknownVersionComps) != 0 {
		// try to resolve version from pip list
		m, e := getEnvPipListMap(ctx)
		if e != nil {
			m = make(map[string]string)
		}
		for s := range unknownVersionComps {
			versionedComps[s] = m[s]
		}
	}
	for s := range noVersionComps {
		if versionedComps[s] != "" {
			continue
		}
		versionedComps[s] = ""
	}

	var result [][2]string
	for k, v := range versionedComps {
		result = append(result, [2]string{k, v})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0] // sort by name is enough
	})

	return result, nil
}

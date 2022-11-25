package python

import (
	"bufio"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io"
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
			if name == "pyproject.toml" {
				return true
			}
			if strings.HasPrefix(name, "requirements") {
				return true
			}
			if filepath.Ext(name) == ".py" {
				return true
			}
		}
	}
	return false
}

func parseDockerFile(dir, path string, m map[string]string) {
	// find all PipManagerFiles from dockerfile
	var regexpToFindPipManagerFiles = `pip\d?\s+install.*?\s-r\s+([^\s&|;"']+)`

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	result := regexp.MustCompile(regexpToFindPipManagerFiles).FindAllStringSubmatch(string(data), -1)
	for _, item := range result {
		if len(item) != 2 {
			continue
		}
		m[item[1]] = filepath.Join(dir, item[1])
	}
}

func scanDepFile(ctx context.Context, dir string) (bool, error) {
	var logger = logctx.Use(ctx)
	var found = false
	var task = model.UseInspectionTask(ctx)

	var tomlFile = filepath.Join(dir, "pyproject.toml")
	var waitingScanPipManagerFiles = make(map[string]string)
	if utils.IsFile(tomlFile) {
		list, e := tomlBuildSysFile(ctx, tomlFile)
		if e != nil {
			logger.Sugar().Warnf("Analyze pyproject.toml failed: %s", e.Error())
		} else if len(list) > 0 {
			task.AddModule(model.Module{
				ModuleName:     "Python-pyprojects.toml",
				ModulePath:     tomlFile,
				PackageManager: "pip",
				Dependencies:   list,
			})
		}
	}

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if strings.Contains(info.Name(), "requirement") {
			waitingScanPipManagerFiles[info.Name()] = path
		}
		if info.Name() == "Dockerfile" {
			parseDockerFile(dir, path, waitingScanPipManagerFiles)
		}
		return nil
	})

	// distinct waitingScanPipManagerFiles
	var pendingScanRequirements []string
	var requirementAbsPathSet = make(map[string]struct{})
	for _, s := range waitingScanPipManagerFiles {
		abs, e := filepath.Abs(filepath.Join(dir, s))
		if e != nil {
			continue
		}
		requirementAbsPathSet[abs] = struct{}{}
	}
	for it := range requirementAbsPathSet {
		r, e := filepath.Rel(dir, it)
		if e != nil {
			continue
		}
		pendingScanRequirements = append(pendingScanRequirements, r)
	}
	sort.Strings(pendingScanRequirements)

	logger.Sugar().Infof("total found pipManagerFiles: %d", len(pendingScanRequirements))
	for _, fp := range pendingScanRequirements {
		logger.Info("start readRequirements...", zap.String("relativePath", fp))
		deps, e := readRequirements(fp)
		if e != nil {
			logger.Sugar().Errorf("Parse requirements file failed[%s]: %s", fp, e.Error())
			continue
		}
		if len(deps) == 0 {
			continue
		}
		found = true
		m := model.Module{
			ModuleName:     fmt.Sprintf("Python-%s", filepath.Base(fp)),
			PackageManager: "pip",
			Dependencies:   deps,
			ModulePath:     fp,
		}
		task.AddModule(m)
	}
	return found, nil
}

func (i Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx)
	dir := model.UseInspectionTask(ctx).Dir()

	foundDepFiles, e := scanDepFile(ctx, dir)
	if e != nil {
		logger.Sugar().Warnf("Scan deps failed: %s", e.Error())
	}
	if foundDepFiles {
		return nil
	}

	componentMap := map[string]string{}
	ignoreSet := map[string]struct{}{}
	logger.Debug("Start walk python project dir", zap.String("dir", dir))
	e = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return nil
		}
		if d.Name() == "venv" && d.IsDir() {
			logger.Debug("Found venv dir, skip", zap.String("dir", path))
			return fs.SkipDir
		}
		if d.IsDir() {
			ignoreSet[d.Name()] = struct{}{}
			return nil
		}
		if filepath.Ext(path) != ".py" {
			return nil
		}
		f, e := os.Open(path)
		if e != nil {
			logger.Sugar().Warnf("Open python file failed: %s, path: %s", e.Error(), path)
			return e
		}
		defer f.Close()
		scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
		scanner.Split(bufio.ScanLines)
		scanner.Buffer(make([]byte, 16*1024), 16*1024)
		for scanner.Scan() {
			if scanner.Err() != nil {
				logger.Sugar().Warnf("Scan python file failed, path: %s, error: %s", path, e.Error())
				return nil
			}
			t := strings.TrimSpace(scanner.Text())
			for _, pkg := range parsePyImport(t) {
				if pyPkgBlackList[pkg] {
					continue
				}
				componentMap[pkg] = ""
			}
		}
		return nil
	})
	if e != nil {
		logger.Warn("Walk error", zap.Error(e))
	}

	if pipListDeps, e := executePipList(ctx, dir); e != nil {
		logger.Warn("pip list execution failed", zap.Error(e))
	} else {
		mergeComponentVersionOnly(componentMap, pipListDeps)
	}

	for s := range ignoreSet {
		delete(componentMap, s)
	}
	if len(componentMap) == 0 {
		logger.Warn("No components valid, omit module")
		return nil
	}
	{
		m := model.Module{
			ModuleName:     filepath.ToSlash(model.UseInspectionTask(ctx).RelDir()),
			PackageManager: "pip",
			ModulePath:     dir,
		}
		if m.ModuleName == "." {
			m.ModuleName = "Python"
		}
		for k, v := range componentMap {
			var di model.DependencyItem
			di.CompName = k
			di.CompVersion = v
			di.EcoRepo = EcoRepo
			m.Dependencies = append(m.Dependencies, di)
		}
		model.UseInspectionTask(ctx).AddModule(m)
		return nil
	}
}

func mergeComponentVersionOnly(target map[string]string, deps []model.DependencyItem) {
	for _, it := range deps {
		v, ok := target[strings.ToLower(it.CompName)]
		if v == "" && ok && it.CompVersion != "" {
			target[it.CompName] = it.CompVersion
		}
	}
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "pip",
	Repository: "",
}

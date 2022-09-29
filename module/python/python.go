package python

import (
	"bufio"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
			if (it.Name() != "conanfile.py" && filepath.Ext(it.Name()) == ".py") || strings.HasPrefix(it.Name(), "requirements") || it.Name() == "pyproject.toml" {
				return true
			}
		}
	}
	return false
}

func scanDepFile(ctx context.Context, dir string) (bool, error) {
	var logger = utils.UseLogger(ctx)
	var found = false
	var task = model.UseInspectorTask(ctx)

	var tomlFile = filepath.Join(dir, "pyproject.toml")
	if utils.IsFile(tomlFile) {
		list, e := tomlBuildSysFile(ctx, tomlFile)
		if e != nil {
			logger.Sugar().Warnf("Analyze pyproject.toml failed: %s", e.Error())
		} else if len(list) > 0 {
			task.AddModule(model.Module{
				Name:           "Python-pyprojects.toml",
				RelativePath:   tomlFile,
				PackageManager: model.PMPip,
				Language:       model.Python,
				Dependencies:   list,
			})
		}
	}

	entries, e := os.ReadDir(dir)
	if e != nil {
		return false, e
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var entryName = entry.Name()
		if !strings.HasPrefix(entryName, "requirements") {
			continue
		}

		fp := filepath.Join(dir, entryName)
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
			Name:           fmt.Sprintf("Python-%s", entryName),
			PackageManager: model.PMPip,
			Language:       model.Python,
			Dependencies:   deps,
			RelativePath:   fp,
		}
		task.AddModule(m)
	}
	return found, nil
}

func (i Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	var relativeDir string
	if s, e := filepath.Rel(task.ProjectDir, task.ScanDir); e == nil {
		relativeDir = filepath.ToSlash(s)
	}
	dir := model.UseInspectorTask(ctx).ScanDir

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
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
			Name:           relativeDir,
			PackageManager: model.PMPip,
			Language:       model.Python,
			Dependencies:   []model.Dependency{},
			RelativePath:   filepath.Join(dir),
		}
		if m.Name == "." {
			m.Name = "Python"
		}
		for k, v := range componentMap {
			m.Dependencies = append(m.Dependencies, model.Dependency{
				Name:    k,
				Version: v,
			})
		}
		model.UseInspectorTask(ctx).AddModule(m)
		return nil
	}
}

func mergeComponentVersionOnly(target map[string]string, deps []model.Dependency) {
	for _, it := range deps {
		v, ok := target[it.Name]
		if v == "" && ok && it.Version != "" {
			target[it.Name] = it.Version
		}
	}
}

func mergeComponentInto(source map[string]string, append []model.Dependency) {
	for _, it := range append {
		name, version := it.Name, it.Version
		if version == "" && source[name] != "" {
			continue
		}
		source[name] = version
	}
}

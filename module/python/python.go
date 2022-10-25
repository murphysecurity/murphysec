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
	"regexp"
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
	var regexpToFindPipManagerFiles = `.*?pip.*?install.*?-r[ ]*(.*?)\s+`

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
	var logger = utils.UseLogger(ctx)
	var found = false
	var task = model.UseInspectorTask(ctx)

	var tomlFile = filepath.Join(dir, "pyproject.toml")
	var waitingScanPipManagerFiles = make(map[string]string)
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
		if strings.HasPrefix(entry.Name(), ".") {
			continue //ignore hide file/folder
		}
		err := filepath.Walk(filepath.Join(dir, entry.Name()), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if strings.Contains(info.Name(), "requirement") {
					waitingScanPipManagerFiles[info.Name()] = path
				}
				if info.Name() == "Dockerfile" {
					parseDockerFile(dir, path, waitingScanPipManagerFiles)
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
	logger.Info(fmt.Sprintf("total found pipManagerFiles: %d", len(waitingScanPipManagerFiles)))
	for entryName, fp := range waitingScanPipManagerFiles {
		logger.Info("start readRequirements...", zap.String("name", entryName), zap.String("relativePath", fp))
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
		v, ok := target[strings.ToLower(it.Name)]
		if v == "" && ok && it.Version != "" {
			target[it.Name] = it.Version
		}
	}
}

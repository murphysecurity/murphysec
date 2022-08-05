package python

import (
	"bufio"
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func (i Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (i Inspector) String() string {
	return "PythonInspector"
}

func (i Inspector) CheckDir(dir string) bool {
	r, e := os.ReadDir(dir)
	if e == nil {
		for _, it := range r {
			if filepath.Ext(it.Name()) == ".py" || strings.HasPrefix(it.Name(), "requirements") || it.Name() == "pyproject.toml" {
				return true
			}
		}
	}
	return false
}

func (i Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	logger := utils.UseLogger(ctx)
	var relativeDir string
	if s, e := filepath.Rel(task.ProjectDir, task.ScanDir); e == nil {
		relativeDir = filepath.ToSlash(s)
	}
	dir := model.UseInspectorTask(ctx).ScanDir
	componentMap := map[string]string{}
	requirementsFiles := map[string]struct{}{}
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
		if (filepath.Ext(path) == ".txt" || filepath.Ext(path) == "") && strings.HasPrefix(d.Name(), "requirements") {
			requirementsFiles[path] = struct{}{}
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
	for fp := range requirementsFiles {
		logger.Debug("Merge requirements file", zap.String("path", fp))
		mergeComponentInto(componentMap, parsePythonRequirements(ctx, fp))
	}

	tomlPath := filepath.Join(dir, "pyproject.toml")
	if utils.IsFile(tomlPath) {
		if list, e := tomlBuildSysFile(ctx, tomlPath); e != nil {
			logger.Sugar().Warnf("Analyze pyproject.toml failed: %s", e.Error())
		} else {
			logger.Sugar().Debug("Merge components from toml build file, total: %d", len(list))
			mergeComponentInto(componentMap, list)
		}
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
			FilePath:       filepath.Join(dir),
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

func mergeComponentInto(source map[string]string, append []model.Dependency) {
	for _, it := range append {
		name, version := it.Name, it.Version
		if version == "" && source[name] != "" {
			continue
		}
		source[name] = version
	}
}

func New() base.Inspector {
	return &Inspector{}
}

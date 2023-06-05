package gradle

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Gradle"
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	var logger = utils.UseLogger(ctx).Sugar()
	var rs []model.Module
	task := model.UseInspectorTask(ctx)
	dir := task.ScanDir
	logger.Debugf("gradle inspect dir: %s", dir)
	if !env.SkipGradleExecution {
		useGradle := true
		gradleEnv, e := DetectGradleEnv(ctx, dir)
		if e != nil {
			task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到目录下没有 gradlew 文件或您的环境中 Gradle 无法正常运行，可能会导致检测结果不完整，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
			logger.Infof("check gradle failed: %s", e.Error())
			logger.Warnf("Gradle disabled")
			useGradle = false
		}
		if useGradle {
			logger.Info(gradleEnv.Version.String())
			var projects []string
			if env.GradleProjects == "" {
				projects, e = fetchGradleProjects(ctx, dir, gradleEnv)
				if e != nil {
					logger.Infof("fetch gradle projects failed: %s", e.Error())
				}
			} else {
				for _, p := range strings.Split(env.GradleProjects, ",") {
					p = strings.TrimSpace(p)
					if p == "" {
						continue
					}
					projects = append(projects, strings.TrimSuffix(p, ":"))
				}
			}

			logger.Debugf("Gradle projects: %s", strings.Join(projects, ", "))

			{
				depInfo, e := evalGradleDependencies(ctx, dir, "", gradleEnv)
				if e != nil {
					logger.Info("evalGradleDependencies failed. <root> ", e.Error())
				} else {
					rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
				}
			}
			for _, projectId := range projects {
				depInfo, e := evalGradleDependencies(ctx, dir, projectId, gradleEnv)
				if e != nil {
					task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 Gradle 获取依赖信息失败，可能会导致检测结果不完整或失败，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
					logger.Infof("evalGradleDependencies failed: %s - %s", projectId, e.Error())
				} else {
					rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
				}
			}
		}
	}
	if len(rs) == 0 {
		// if no module find, use backup solution
		if m := backupParser(ctx, dir); m != nil {
			tm := m.BaseModule(dir)
			tm.ScanStrategy = model.ScanStrategyBackup
			rs = append(rs, tm)
		}
	}
	for _, it := range rs {
		task.AddModule(it)
	}
	return nil
}

func backupParser(ctx context.Context, dir string) *GradleDependencyInfo {
	var logger = utils.UseLogger(ctx).Sugar()
	var dep []DepElement
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || !d.Type().IsRegular() || !strings.HasPrefix(d.Name(), "build.gradle") {
			return nil
		}
		if d.Name() == "build.gradle.kts" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Errorf("Read gradle file failed: %s", e.Error())
				return nil
			}
			dep = append(dep, parseGradleKts(string(data))...)
			return nil
		}
		if d.Name() == "build.gradle" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Errorf("Read gradle file failed: %s", e.Error())
				return nil
			}
			dep = append(dep, parseGradleGroovy(string(data))...)
			return nil
		}
		return nil
	})
	if len(dep) != 0 {
		return &GradleDependencyInfo{
			ProjectName:  fmt.Sprintf("GradleProject-%s", filepath.Base(dir)),
			Dependencies: dep,
		}
	}
	return nil
}

var gradleBuildFiles = []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"}

func (i *Inspector) CheckDir(dir string) bool {
	for _, it := range gradleBuildFiles {
		info, e := os.Stat(filepath.Join(dir, it))
		if e == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

// fetchGradleProjects evaluate `gradle projects` and parse the result, then returns a project identifier list.
func fetchGradleProjects(ctx context.Context, projectDir string, info *GradleEnv) ([]string, error) {
	c := info.ExecuteContext(ctx, "projects")
	c.Dir = projectDir
	pattern := regexp.MustCompile("Project\\s+'(:.+?)'")
	output, e := c.Output()
	if e != nil {
		return nil, e
	}
	m := map[string]struct{}{}
	for _, match := range pattern.FindAllStringSubmatch(string(output), -1) {
		if len(match) < 2 || match[1] == "" {
			continue
		}
		m[match[1]] = struct{}{}
	}
	var rs []string
	for s := range m {
		rs = append(rs, s)
	}
	return rs, nil
}

func evalGradleDependencies(ctx context.Context, projectDir string, projectName string, info *GradleEnv) (*GradleDependencyInfo, error) {
	var logger = utils.UseLogger(ctx).Sugar()
	c := info.ExecuteContext(ctx, fmt.Sprintf("%s:dependencies", projectName))
	logger.Infof("Execute: %s", c.String())
	c.Dir = projectDir
	data, e := c.Output()
	logger.Debugf("GradleOutput: %s", string(data))
	if e != nil {
		logger.Errorf("Gradle output: %s", string(e.(*exec.ExitError).Stderr))
		return nil, e
	}
	var lines []string
	for _, it := range strings.Split(string(data), "\n") {
		lines = append(lines, strings.TrimSpace(it))
	}
	depInfo := parseGradleDependencies(lines)
	if depInfo == nil {
		return nil, errors.New("parse dep info failed.")
	}
	return depInfo, nil
}

//goland:noinspection GoNameStartsWithPackageName
type GradleDependencyInfo struct {
	ProjectName  string       `json:"project_name"`
	Dependencies []DepElement `json:"dependencies,omitempty"`
}

func (g *GradleDependencyInfo) BaseModule(path string) model.Module {
	return model.Module{
		PackageManager: model.PMGradle,
		Language:       model.Java,
		Name:           g.ProjectName,
		Dependencies:   _convDep(g.Dependencies),
		RelativePath:   path,
	}
}

func _convDep(input []DepElement) []model.Dependency {
	var rs []model.Dependency
	for _, it := range input {
		rs = append(rs, model.Dependency{
			Name:         it.CompName(),
			Version:      it.Version,
			Dependencies: _convDep(it.Children),
		})
	}
	return rs
}

func parseGradleDependencies(lines []string) *GradleDependencyInfo {
	info := &GradleDependencyInfo{
		ProjectName:  "",
		Dependencies: []DepElement{},
	}
	taskPattern := regexp.MustCompile("^\\w+$|^\\w+\\s-")
	projectPattern := regexp.MustCompile("(?:Root project|[Pp]roject) ([':A-Za-z0-9._-]+)")
	type task struct {
		name  string
		lines []string
	}
	var tasks []task
	{
		var currTaskName string
		var currTaskLines []string
		for _, it := range lines {
			if m := projectPattern.FindStringSubmatch(it); len(m) > 0 && info.ProjectName == "" {
				info.ProjectName = strings.TrimSpace(strings.Trim(m[1], "'"))
				continue
			}
			if it == "" {
				if currTaskName != "" {
					tasks = append(tasks, task{currTaskName, currTaskLines})
					currTaskLines = nil
					currTaskName = ""
				}
				continue
			}
			if m := taskPattern.FindString(it); m != "" {
				if currTaskName != "" {
					tasks = append(tasks, task{currTaskName, currTaskLines})
					currTaskLines = nil
					currTaskName = ""
				}
				currTaskName = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(m), "-"))
				continue
			}
			if currTaskName == "" {
				continue
			}
			currTaskLines = append(currTaskLines, it)
		}
	}
	{
		for _, task := range tasks {
			if task.name == "runtimeClasspath" {
				parser := blockParser{lines: task.lines}
				info.Dependencies = parser._parse()
			}
		}
	}
	return info
}

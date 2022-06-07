package gradle

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "GradleInspector"
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	var rs []model.Module
	task := model.UseInspectorTask(ctx)
	dir := task.ScanDir
	logger.Debug.Println("gradle inspect dir:", dir)
	useGradle := true
	ctx, cf := context.WithTimeout(context.TODO(), time.Second*time.Duration(env.GradleExecutionTimeoutSecond))
	defer cf()
	info, e := evalGradleInfo(ctx, dir)
	if e != nil {
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到目录下没有 gradlew 文件或您的环境中 Gradle 无法正常运行，可能会导致检测结果不完整，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
		logger.Info.Println("check gradle failed", e.Error())
		logger.Warn.Println("Gradle disabled")
		useGradle = false
	}
	if useGradle {
		logger.Info.Println(info)
		projects, e := fetchGradleProjects(ctx, dir, info)
		if e != nil {
			logger.Info.Println("fetch gradle projects failed.", e.Error())
		}
		logger.Debug.Println("Gradle projects:", strings.Join(projects, ", "))

		{
			depInfo, e := evalGradleDependencies(ctx, dir, "", info)
			if e != nil {
				logger.Info.Println("evalGradleDependencies failed. <root>", e.Error())
			} else {
				rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
			}
		}
		for _, projectId := range projects {
			depInfo, e := evalGradleDependencies(ctx, dir, projectId, info)
			if e != nil {
				task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 Gradle 获取依赖信息失败，可能会导致检测结果不完整或失败，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
				logger.Info.Println("evalGradleDependencies failed.", projectId, e.Error())
			} else {
				rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
			}
		}
	}
	if len(rs) == 0 {
		// if no module find, use backup solution
		if m := backupParser(dir); m != nil {
			rs = append(rs, m.BaseModule(dir))
		}
	}
	for _, it := range rs {
		task.AddModule(it)
	}
	return nil
}

func backupParser(dir string) *GradleDependencyInfo {
	var dep []DepElement
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || !d.Type().IsRegular() || !strings.HasPrefix(d.Name(), "build.gradle") {
			return nil
		}
		if d.Name() == "build.gradle.kts" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Err.Println("Read gradle file failed.", e.Error())
				return nil
			}
			dep = append(dep, parseGradleKts(string(data))...)
			return nil
		}
		if d.Name() == "build.gradle" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Err.Println("Read gradle file failed.", e.Error())
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
func fetchGradleProjects(ctx context.Context, projectDir string, info *GradleInfo) ([]string, error) {
	c := info.CallCmd(ctx, "--console", "plain", "-q", "projects")
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

func evalGradleDependencies(ctx context.Context, projectDir string, projectName string, info *GradleInfo) (*GradleDependencyInfo, error) {
	c := info.CallCmd(ctx, fmt.Sprintf("%s:dependencies", projectName), "--console", "plain", "-q", "--configuration=runtimeClasspath")
	logger.Debug.Println("Execute:", c.String())
	c.Dir = projectDir
	data, e := c.Output()
	logger.Debug.Println("GradleOutput:", string(data))
	if e != nil {
		logger.Debug.Println("Gradle output", string(e.(*exec.ExitError).Stderr))
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
		FilePath:       path,
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
	projectPattern := regexp.MustCompile("(?:Root project|project) '([A-Za-z0-9._-]+)'")
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
				info.ProjectName = strings.TrimSpace(strings.TrimPrefix(m[1], "Project"))
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

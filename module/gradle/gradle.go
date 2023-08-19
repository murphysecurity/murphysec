package gradle

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/sl"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"github.com/repeale/fp-go"
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
	var logger = logctx.Use(ctx).Sugar()
	var rs []model.Module
	task := model.UseInspectionTask(ctx)
	dir := task.Dir()
	logger.Debugf("gradle inspect dir: %s", dir)
	useGradle := true
	gradleEnv, e := DetectGradleEnv(ctx, dir)
	if e != nil {
		logger.Infof("check gradle failed: %s", e.Error())
		logger.Warnf("Gradle disabled")
		useGradle = false
	}
	if useGradle {
		logger.Info(gradleEnv.Version.String())
		projects, e := fetchGradleProjects(ctx, dir, gradleEnv)
		if e != nil {
			logger.Infof("fetch gradle projects failed: %s", e.Error())
		}
		logger.Debugf("Gradle projects: %s", strings.Join(projects, ", "))

		{
			depInfo, e := evalGradleDependencies(ctx, dir, "", gradleEnv)
			if e != nil {
				logger.Info("evalGradleDependencies failed. <root> ", e.Error())
			} else {
				rs = append(rs, depInfo.BaseModule(dir))
			}
		}
		for _, projectId := range projects {
			depInfo, e := evalGradleDependencies(ctx, dir, projectId, gradleEnv)
			if e != nil {
				logger.Infof("evalGradleDependencies failed: %s - %s", projectId, e.Error())
			} else {
				rs = append(rs, depInfo.BaseModule(dir))
			}
		}
	}
	if len(rs) == 0 && !env.ScannerScan {
		// if no module find, use backup solution
		if m := backupParser(ctx, dir); m != nil {
			tm := m.BaseModule(dir)
			tm.ScanStrategy = model.ScanStrategyBackup
			rs = append(rs, tm)
		}
	}
	env.ScannerShouldEnableGradleBackupScan = true
	for _, i := range rs {
		if len(i.Dependencies) == 0 {
			env.ScannerShouldEnableGradleBackupScan = false
			break
		}
	}
	for _, it := range rs {
		task.AddModule(it)
	}
	return nil
}

func gradleProjectName2PathComponent(baseDir string, projectName string) string {
	part2 := strings.Join(fp.Filter(sl.StringNotEmpty)(strings.Split(projectName, ":")), "/")
	if f := filepath.Join(baseDir, part2, "build.gradle.kts"); utils.IsFile(f) {
		return f
	}
	return filepath.Join(baseDir, part2, "build.gradle")
}

func backupParser(ctx context.Context, dir string) *GradleDependencyInfo {
	var logger = logctx.Use(ctx).Sugar()
	var dep []DepElement
	e := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
	if e != nil {
		logger.Warnf("Walk: %v", e)
	}
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
	pattern := regexp.MustCompile(`Project\s+'(:.+?)'`)
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
	var logger = logctx.Use(ctx).Sugar()
	c := info.ExecuteContext(ctx, fmt.Sprintf("%s:dependencies", projectName))
	logger.Infof("Execute: %s", c.String())
	c.Dir = projectDir
	data, e := c.Output()
	logger.Debugf("GradleOutput: %s", string(data))
	if e != nil {
		logger.Errorf("Gradle output: %s", string(e.(*exec.ExitError).Stderr))
		return nil, e
	}
	var lines = fp.Map(strings.TrimSpace)(strings.Split(string(data), "\n"))
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

func (g *GradleDependencyInfo) BaseModule(basePath string) model.Module {
	return model.Module{
		PackageManager: "gradle",
		ModuleName:     g.ProjectName,
		Dependencies:   convDep(g.Dependencies),
		ModulePath:     gradleProjectName2PathComponent(basePath, g.ProjectName),
	}
}

func convDep(input []DepElement) []model.DependencyItem {
	var r = _convDep(input)
	for i := range r {
		r[i].IsDirectDependency = true
	}
	return r
}

func _convDep(input []DepElement) []model.DependencyItem {
	var rs []model.DependencyItem
	for _, it := range input {
		rs = append(rs, model.DependencyItem{
			Component: model.Component{
				CompName:    it.CompName(),
				CompVersion: it.Version,
				EcoRepo:     EcoRepo,
			},
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
	taskPattern := regexp.MustCompile(`^\w+$|^\w+\s-`)
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

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}

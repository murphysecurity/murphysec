package gradle

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "GradleInspector"
}

func (i *Inspector) Version() string {
	return "0.0.1"
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMGradle
}

func (i *Inspector) Inspect(dir string) ([]base.Module, error) {
	logger.Debug.Println("gradle inspect dir:", dir)
	info, e := evalGradleInfo(dir)
	if e != nil {
		logger.Info.Println("check gradle failed", e.Error())
		return nil, e
	}
	logger.Info.Println(info)
	projects, e := fetchGradleProjects(dir, info)
	if e != nil {
		logger.Info.Println("fetch gradle projects failed.", e.Error())
	}
	logger.Debug.Println("Gradle projects:", strings.Join(projects, ", "))
	var rs []base.Module
	{
		depInfo, e := evalGradleDependencies(dir, "", info)
		if e != nil {
			logger.Info.Println("evalGradleDependencies failed. <root>", e.Error())
		} else {
			rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
		}
	}
	for _, projectId := range projects {
		depInfo, e := evalGradleDependencies(dir, projectId, info)
		if e != nil {
			logger.Info.Println("evalGradleDependencies failed.", projectId, e.Error())
		} else {
			rs = append(rs, depInfo.BaseModule(filepath.Join(dir, "build.gradle")))
		}
	}
	return rs, nil // todo
}

func (i *Inspector) CheckDir(dir string) bool {
	info, e := os.Stat(filepath.Join(dir, "build.gradle"))
	if e == nil && !info.IsDir() {
		return true
	}
	info, e = os.Stat(filepath.Join(dir, "build.gradle.kts"))
	if e == nil && !info.IsDir() {
		return true
	}
	return false
}

// fetchGradleProjects evaluate `gradle projects` and parse the result, then returns a project identifier list.
func fetchGradleProjects(projectDir string, info *GradleInfo) ([]string, error) {
	c := info.CallCmd("--console", "plain", "-q", "projects")
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

func evalGradleDependencies(projectDir string, projectName string, info *GradleInfo) (*GradleDependencyInfo, error) {
	c := info.CallCmd(fmt.Sprintf("%s:dependencies", projectName), "--console", "plain", "-q", "--configuration=runtimeClasspath")
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

func (g *GradleDependencyInfo) BaseModule(path string) base.Module {
	return base.Module{
		PackageManager: "gradle",
		Language:       "java",
		Name:           g.ProjectName,
		Dependencies:   _convDep(g.Dependencies),
		FilePath:       path,
	}
}

func _convDep(input []DepElement) []base.Dependency {
	var rs []base.Dependency
	for _, it := range input {
		rs = append(rs, base.Dependency{
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

package inspector

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"regexp"
	"sort"
	"strings"
)

func InspectDockerfile(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	dockerfilePath := scanTask.ProjectDir
	data, e := os.ReadFile(dockerfilePath)
	if e != nil {
		ui.Display(display.MsgError, fmt.Sprintf("读取 Dockerfile 失败：%s", e.Error()))
		return e
	}
	if e := createTaskC(ctx); e != nil {
		return e
	}
	r := analyzeDockerfile(string(data))
	module := model.Module{
		Name: "Dockerfile",
	}
	for _, it := range r {
		module.Dependencies = append(module.Dependencies, model.Dependency{
			Name:    it.Name,
			Version: it.Version,
		})
	}
	scanTask.Modules = append(scanTask.Modules, module)
	if e := submitModuleInfoC(ctx); e != nil {
		return nil
	}
	if e := startCheckC(ctx); e != nil {
		return nil
	}
	if e := queryResultC(ctx); e != nil {
		return nil
	}
	if scanTask.ScanResult.ReportURL() != "" {
		ui.Display(display.MsgNotice, fmt.Sprintf("检测报告详见：%s", scanTask.ScanResult.ReportURL()))
	}
	return e
}

type DockerfileResult []DockerfileItem

func (d DockerfileResult) Len() int {
	return len(d)
}

func (d DockerfileResult) Less(i, j int) bool {
	if d[i].Kind != d[j].Kind {
		return d[i].Kind < d[j].Kind
	}
	if d[i].Name != d[j].Name {
		return d[i].Name < d[j].Name
	}
	if d[i].Version != d[j].Version {
		return d[i].Version < d[j].Version
	}
	return false
}

func (d DockerfileResult) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func analyzeDockerfile(data string) DockerfileResult {

	var items DockerfileResult
	for _, line := range strings.Split(replaceLineBreak(data), "\n") {
		line = strings.TrimSuffix(line, "\r")
		hashIndex := strings.Index(line, "#")
		if hashIndex > -1 && hashIndex > len(line) {
			line = line[:hashIndex]
		}
		if i := imageChecker(line); i != nil {
			items = append(items, *i)
		}
		for _, i := range yumChecker(line) {
			items = append(items, i)
		}
		for _, i := range debChecker(line) {
			items = append(items, i)
		}
	}
	sort.Sort(items)

	return items
}

func replaceLineBreak(input string) string {
	return regexp.MustCompile("\\\\\\r?\\n").ReplaceAllString(input, "")
}

var spacePattern = regexp.MustCompile("\\s+")

type DockerfileItem struct {
	Kind    string
	Name    string
	Version string
}

var imagesPattern = regexp.MustCompile("FROM\\s+(\\S+)")

func imageChecker(line string) *DockerfileItem {
	if m := imagesPattern.FindStringSubmatch(line); m != nil {
		return &DockerfileItem{
			Kind: "image",
			Name: m[1],
		}
	}
	return nil
}

var aptInstallPattern = regexp.MustCompile("apt(?:-get)?\\s+install\\s([\\w\\s.-]+)")

func debChecker(line string) (r []DockerfileItem) {
	m := aptInstallPattern.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	var listStr = m[1]
	list := spacePattern.Split(listStr, -1)
	for _, s := range list {
		if strings.HasPrefix(s, "-") {
			continue
		}
		r = append(r, DockerfileItem{
			Kind: "deb",
			Name: s,
		})
	}
	return
}

var yumInstallPattern = regexp.MustCompile("yum\\s+install\\s([\\w\\s.-]+)")

func yumChecker(line string) (r []DockerfileItem) {
	m := yumInstallPattern.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	var listStr = m[1]
	list := spacePattern.Split(listStr, -1)
	for _, s := range list {
		if strings.HasPrefix(s, "-") {
			continue
		}
		r = append(r, DockerfileItem{
			Kind: "rpm",
			Name: s,
		})
	}
	return
}

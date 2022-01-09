package go_mod

import (
	"bufio"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Inspector struct{}

func (i *Inspector) String() string {
	return "GoModInspector@" + i.Version()
}

func (i *Inspector) Version() string {
	return "v0.0.1"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (i *Inspector) Inspect(dir string) ([]base.Module, error) {
	return ScanGoProject(dir)
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMGoMod
}

func New() base.Inspector {
	return &Inspector{}
}

func ScanGoProject(dir string) ([]base.Module, error) {
	version, e := execGoVersion()
	if e != nil {
		return nil, e
	}
	if e := execGoModTidy(dir); e != nil {
		logger.Err.Println("go mod tidy execute failed.", e.Error())
		return nil, e
	}
	deps, e := execGoModGraph(dir)
	if e != nil {
		logger.Err.Println("Scan go project failed, ", e.Error())
		return nil, e
	}

	module := base.Module{
		PackageManager: "Go",
		Language:       "Go",
		PackageFile:    "go.mod",
		Name:           deps.Name,
		Version:        deps.Version,
		RelativePath:   "go.mod",
		Dependencies:   deps.Dependencies,
		RuntimeInfo:    map[string]interface{}{"go_version": version},
	}
	return []base.Module{module}, nil
}

func execGoModGraph(dir string) (*base.Dependency, error) {
	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = dir
	r, w := io.Pipe()
	defer must.Close(w)
	cmd.Stdout = w
	if e := cmd.Start(); e != nil {
		logger.Err.Println("go mod graph execute failed.", e.Error())
		return nil, e
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var deps *base.Dependency
	go func() {
		deps = parseGoModGraph(r)
		wg.Done()
	}()
	if e := cmd.Wait(); e != nil {
		logger.Err.Println("go mod graph exit with err:", e.Error())
		return nil, e
	} else {
		logger.Info.Println("go mod graph exit with no error.")
	}
	must.Close(w)
	must.Close(w)
	wg.Wait()
	return deps, nil
}

func parseGoModGraph(reader io.Reader) *base.Dependency {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	scanner.Buffer(make([]byte, 1024*4), 1024*4)
	list := map[string][]string{}
	root := ""
	count := 0
	for scanner.Scan() {
		t := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(t) != 2 {
			logger.Debug.Println("Unrecognized line:", scanner.Text())
			continue
		}
		count++
		if root == "" {
			root = t[0]
		}
		list[t[0]] = append(list[t[0]], t[1])
	}
	logger.Debug.Println("go mod graph: Total process lines ", count)
	return _convDependency(list, root, nil)
}

func _convDependency(m map[string][]string, root string, visited []string) *base.Dependency {
	if len(visited) > 3 {
		return nil
	}
	for _, it := range visited {
		if it == root {
			logger.Warn.Println("Circular dependency:", strings.Join(visited, " -> "))
			return nil
		}
	}
	t := strings.Split(root, "@")
	if len(t) > 2 {
		logger.Debug.Println("Invalid id:", root)
		return nil
	}
	if len(t) == 1 {
		t = []string{t[0], ""}
	}
	d := base.Dependency{
		Name:    t[0],
		Version: t[1],
	}
	for _, it := range m[root] {
		if r := _convDependency(m, it, append(visited, root)); r != nil {
			d.Dependencies = append(d.Dependencies, *r)
		}
	}
	return &d
}

func execGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy", "-v")
	cmd.Dir = dir
	r, w := io.Pipe()
	defer must.Close(w)
	cmd.Stdout = w
	go func() {
		buf := bufio.NewScanner(r)
		buf.Split(bufio.ScanLines)
		buf.Buffer(make([]byte, 24*1024), 24*2014)
		for buf.Scan() {
			logger.Err.Println("go mod tidy:", buf.Text())
		}
	}()
	if e := cmd.Start(); e != nil {
		logger.Err.Println("Execute go mod tidy failed.", e.Error())
		return e
	}
	if e := cmd.Wait(); e != nil {
		logger.Err.Println("go mod tidy exit with errors.", e.Error())
	} else {
		logger.Info.Println("go mod tidy exit with no error.")
	}
	return nil
}

func execGoVersion() (string, error) {
	v, e := exec.Command("go", "version").Output()
	if e != nil {
		logger.Err.Println("go version execute failed.", e.Error())
		return "", e
	} else {
		logger.Info.Println("go version execute succeed.")
	}
	logger.Info.Println("go version:", string(v))
	return string(v), nil
}

package go_mod

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/simplejson"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	root, e := execGoListModule(dir)
	if e != nil {
		logger.Err.Println("execGoListModule:", e.Error())
		// language=json
		root, _ = simplejson.NewJSON([]byte(`{"Module": {"Path": "main"}}`))
	}

	deps, e := execGoList(dir)
	if e != nil {
		logger.Err.Println("Scan go project failed, ", e.Error())
		return nil, e
	}

	module := base.Module{
		PackageManager: "Go",
		Language:       "Go",
		PackageFile:    "go.mod",
		Name:           root.Get("Module", "Path").String("main"),
		Version:        "",
		RelativePath:   "go.mod",
		Dependencies:   deps,
		RuntimeInfo:    map[string]interface{}{"go_version": version},
	}
	return []base.Module{module}, nil
}

func execGoListModule(dir string) (*simplejson.JSON, error) {
	cmd := exec.Command("go", "list", "--json")
	cmd.Dir = dir
	data, e := cmd.Output()
	if e != nil {
		return nil, e
	}
	var d *simplejson.JSON
	if e := json.Unmarshal(data, &d); e != nil {
		return nil, e
	}
	if d == nil {
		return nil, errors.New("json is nil")
	}
	return d, nil
}

func execGoList(dir string) ([]base.Dependency, error) {
	cmd := exec.Command("go", "list", "--json", "-m", "all")
	cmd.Dir = dir
	data, e := cmd.Output()
	if e != nil {
		logger.Err.Println("go list execute failed.", e.Error())
		return nil, errors.New("Go list execute failed")
	}
	dep := make([]base.Dependency, 0)
	dec := json.NewDecoder(bytes.NewReader(data))
	for {
		var m *simplejson.JSON
		if e := dec.Decode(&m); e == io.EOF {
			break
		} else if e != nil {
			logger.Err.Println(e.Error())
			return nil, errors.Wrap(e, "parse go list failed")
		}
		if m == nil {
			continue
		}
		if m.Get("Version").String() == "" {
			continue
		}
		if !m.Get("Replace").IsNull() {
			replacePath := m.Get("Replace", "Path").String("")
			if replacePath == "" {
				continue
			}
			replaceVersion := m.Get("Replace", "Version").String()
			dep = append(dep, base.Dependency{
				Name:    replacePath,
				Version: replaceVersion,
			})
			continue
		}
		dep = append(dep, base.Dependency{
			Name:         m.Get("Path").String(),
			Version:      m.Get("Version").String(),
			Dependencies: []base.Dependency{},
		})
	}
	return dep, nil
}

func execGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy", "-v")
	cmd.Dir = dir
	envs := make([]string, 0)
	// Workaround for some environment that disables "go mod" default.
	for _, it := range os.Environ() {
		if strings.HasPrefix("GO111MODULE=", it) {
			continue
		}
		envs = append(envs, it)
	}
	envs = append(envs, "GO111MODULE=on")
	cmd.Env = envs
	stdoutReader, e := cmd.StdoutPipe()
	if e != nil {
		logger.Warn.Println("Get stdout pipe failed", e.Error())
		return e
	}
	go func() {
		buf := bufio.NewScanner(stdoutReader)
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

package go_mod

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/simplejson"
	"os/exec"
	"path/filepath"
)

var ErrGoEnv = model.NewInspectError(model.Go, "Check Go version failed. Please check you go environment.")

type Inspector struct{}

func (i *Inspector) String() string {
	return "GoModInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	return ScanGoProject(task)
}

func New() base.Inspector {
	return &Inspector{}
}

func ScanGoProject(task *model.InspectorTask) error {
	dir := task.ScanDir
	version, e := execGoVersion()
	if e != nil {
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]识别到您的环境中 Go 无法正常运行，可能会导致检测结果不完整，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
		return ErrGoEnv
	}
	if e := execGoModTidy(dir); e != nil {
		task.UI().Display(display.MsgWarn, fmt.Sprintf("[%s]通过 Go获取依赖信息失败，可能会导致检测结果不完整或失败，访问https://www.murphysec.com/docs/quick-start/language-support/ 了解详情", dir))
		logger.Err.Println("go mod tidy execute failed.", e.Error())
		return e
	}
	root, e := execGoListModule(dir)
	if e != nil {
		logger.Err.Println("execGoListModule:", e.Error())
		root = simplejson.New()
	}

	deps, e := execGoList(dir)
	if e != nil {
		logger.Err.Println("Scan go project failed, ", e.Error())
		return e
	}
	module := model.Module{
		PackageManager: model.PMGoMod,
		Language:       model.Go,
		PackageFile:    "go.mod",
		Name:           root.Get("Module", "Path").String(filepath.Base(dir)),
		Version:        "",
		FilePath:       filepath.Join(dir, "go.mod"),
		Dependencies:   deps,
		RuntimeInfo:    map[string]interface{}{"go_version": version},
	}
	task.AddModule(module)
	return nil
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

func execGoList(dir string) ([]model.Dependency, error) {
	cmd := exec.Command("go", "list", "--json", "-m", "all")
	cmd.Dir = dir
	data, e := cmd.Output()
	if e != nil {
		logger.Err.Println("go list execute failed.", e.Error())
		return nil, errors.New("Go list execute failed")
	}
	dep := make([]model.Dependency, 0)
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
			dep = append(dep, model.Dependency{
				Name:    replacePath,
				Version: replaceVersion,
			})
			continue
		}
		dep = append(dep, model.Dependency{
			Name:         m.Get("Path").String(),
			Version:      m.Get("Version").String(),
			Dependencies: []model.Dependency{},
		})
	}
	return dep, nil
}

func execGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy", "-v")
	cmd.Dir = dir
	logger.Debug.Println("Execute:", cmd.String(), cmd.Dir)
	output, e := cmd.CombinedOutput()
	if e == nil {
		logger.Info.Println("go mod tidy exit with no error.")
		return nil
	} else {
		logger.Err.Println("go mod tidy exit with errors.", e.Error())
		logger.Debug.Println("Output:", string(output))
		return errors.Wrap(e, "Go mod tidy execution failed.")
	}
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

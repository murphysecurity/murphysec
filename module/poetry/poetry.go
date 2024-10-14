package poetry

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrParsePoetry = poetryErr("ErrParsePoetry: Bad manifest")

type poetryErr string

func (p poetryErr) Error() string {
	return string(p)
}

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Poetry"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pyproject.toml"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	pyprojectFile := filepath.Join(task.Dir(), "pyproject.toml")
	data, e := utils.ReadFileLimited(pyprojectFile, 1024*1024*4)
	if e != nil {
		return errors.Wrap(e, "Read pyproject.toml fail")
	}
	manifest, e := parsePoetry(data)
	if e != nil {
		return errors.Wrap(e, "PoetryInspector")
	}
	cmap := map[string]string{}
	for _, it := range manifest.Dependencies {
		cmap[it.CompName] = it.CompVersion
	}
	poetryFile := filepath.Join(task.Dir(), "poetry.lock.py")
	if !utils.IsFile(poetryFile) {
		poetryFile = filepath.Join(task.Dir(), "poetry.lock")
	}
	if utils.IsFile(poetryFile) {
		if deps, e := parsePoetryLock(ctx, poetryFile); e == nil {
			for _, it := range deps {
				cmap[it.CompName] = it.CompVersion
			}
		}
	}
	module := model.Module{
		PackageManager: "poetry",
		ModuleName:     manifest.Name,
		ModulePath:     pyprojectFile,
	}
	for k, v := range cmap {
		var di model.DependencyItem
		di.CompName = k
		di.CompVersion = v
		di.EcoRepo = EcoRepo
		module.Dependencies = append(module.Dependencies, di)
	}

	task.AddModule(module)
	return nil
}

type Manifest struct {
	Name         string
	Dependencies []model.DependencyItem
}

var splitVersionPattern = regexp.MustCompile("[<>=~]?=")

func parsePoetry(input []byte) (*Manifest, error) {
	root := &tomlTree{}
	if e := toml.Unmarshal(input, &root.v); e != nil {
		return nil, errors.WithMessage(ErrParsePoetry, "Parse toml failed")
	}
	var m = make(map[string]any)
	if mm, ok := root.Get("tool", "poetry", "dependencies").v.(map[string]any); ok {
		for k, v := range mm {
			m[k] = v
		}
	}
	if mm, ok := root.Get("project", "dependencies").v.([]any); ok {
		for _, _s := range mm {
			var s, ok = _s.(string)
			if !ok {
				continue
			}
			var r = splitVersionPattern.Split(s, 2)
			if len(r) != 2 || strings.TrimSpace(r[0]) == "" {
				continue
			}
			m[strings.TrimSpace(r[0])] = strings.TrimSpace(r[1])
		}
	}
	var deps []model.DependencyItem
	for k, v := range m {
		v := strings.Trim(fmt.Sprint(v), "~^* ")
		if v == "" {
			continue
		}
		var di model.DependencyItem
		di.CompName = k
		di.CompVersion = v
		di.EcoRepo = EcoRepo
		deps = append(deps, di)
	}
	return &Manifest{
		Name:         root.Get("tool", "poetry", "name").String("<noname>"),
		Dependencies: deps,
	}, nil
}

type tomlTree struct {
	v any
}

func (t *tomlTree) AsArray() (rs []tomlTree) {
	arr, ok := t.v.([]any)
	if !ok {
		return
	}
	for _, it := range arr {
		rs = append(rs, tomlTree{v: it})
	}
	return
}

func (t *tomlTree) Get(path ...string) *tomlTree {
	cur := t
	for _, it := range path {
		m, ok := cur.v.(map[string]any)
		if !ok {
			return &tomlTree{}
		}
		cur = &tomlTree{m[it]}
	}
	return cur
}

func (t tomlTree) String(a ...string) string {
	if len(a) > 1 {
		panic("bad args")
	}
	s, ok := t.v.(string)
	if ok {
		return s
	} else {
		if len(a) == 1 {
			return a[0]
		} else {
			return ""
		}
	}
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "pip",
	Repository: "",
}

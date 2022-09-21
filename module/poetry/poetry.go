package poetry

import (
	"context"
	"github.com/google/uuid"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

var ErrParsePoetry = poetryErr("ErrParsePoetry: Bad manifest")

type poetryErr string

func (p poetryErr) Error() string {
	return string(p)
}

type Inspector struct{}

var Instance = &Inspector{}

func (i *Inspector) SupportFeature(feature base.Feature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Poetry"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pyproject.toml"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	pyprojectFile := filepath.Join(task.ScanDir, "pyproject.toml")
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
		cmap[it.Name] = it.Version
	}
	poetryFile := filepath.Join(task.ScanDir, "poetry.lock.py")
	if utils.IsFile(poetryFile) {
		if deps, e := parsePoetryLock(poetryFile); e == nil {
			for _, it := range deps {
				cmap[it.Name] = it.Version
			}
		}
	}
	module := model.Module{
		PackageManager: model.PMPoetry,
		Language:       model.Python,
		Name:           manifest.Name,
		Dependencies:   []model.Dependency{},
		UUID:           uuid.Must(uuid.NewRandom()),
		RelativePath:   task.ScanDir,
	}
	for k, v := range cmap {
		module.Dependencies = append(module.Dependencies, model.Dependency{Name: k, Version: v})
	}

	task.AddModule(module)
	return nil
}

type Manifest struct {
	Name         string
	Dependencies []model.Dependency
}

func parsePoetry(input []byte) (*Manifest, error) {
	root := &tomlTree{}
	if e := toml.Unmarshal(input, &root.v); e != nil {
		return nil, errors.WithMessage(ErrParsePoetry, "Parse toml failed")
	}
	m, ok := root.Get("tool", "poetry", "dependencies").v.(map[string]string)
	if !ok {
		return nil, errors.WithMessage(ErrParsePoetry, "bad toml")
	}
	var deps []model.Dependency
	for k, v := range m {
		v := strings.Trim(v, "~^* ")
		if v == "" {
			continue
		}
		deps = append(deps, model.Dependency{
			Name:    k,
			Version: v,
		})
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

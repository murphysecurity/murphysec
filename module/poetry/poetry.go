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

var ErrParsePoetry = errors.New("Bad manifest")

type Inspector struct{}

func (i *Inspector) String() string {
	return "PoetryInspector"
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
	task.AddModule(model.Module{
		PackageManager: model.PMPoetry,
		Language:       model.Python,
		PackageFile:    "pyprojject.toml",
		Name:           manifest.Name,
		Dependencies:   manifest.Dependencies,
		UUID:           uuid.Must(uuid.NewRandom()),
	})
	return nil
}

func New() base.Inspector {
	return &Inspector{}
}

type Manifest struct {
	Name         string
	Dependencies []model.Dependency
}

func parsePoetry(input []byte) (*Manifest, error) {
	root := &tomlTree{}
	if e := toml.Unmarshal(input, &root.v); e != nil {
		return nil, errors.Wrap(ErrParsePoetry, "Parse toml failed")
	}
	m, ok := root.Get("tool", "poetry", "dependencies").v.(map[string]string)
	if !ok {
		return nil, errors.Wrap(ErrParsePoetry, "bad toml")
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

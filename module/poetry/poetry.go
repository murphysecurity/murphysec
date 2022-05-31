package poetry

import (
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
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

func (i *Inspector) Inspect(task *base.ScanTask) ([]base.Module, error) {
	pyprojectFile := filepath.Join(task.ProjectDir, "pyproject.toml")
	data, e := utils.ReadFileLimited(pyprojectFile, 1024*1024*4)
	if e != nil {
		return nil, errors.Wrap(e, "Read pyproject.toml fail")
	}
	manifest, e := parsePoetry(data)
	if e != nil {
		return nil, errors.Wrap(e, "PoetryInspector")
	}

	return []base.Module{{
		PackageManager: "poetry",
		Language:       "python",
		PackageFile:    "pyprojject.toml",
		Name:           manifest.Name,
		Dependencies:   manifest.Dependencies,
		UUID:           uuid.Must(uuid.NewRandom()),
	}}, nil
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMPython
}

func New() base.Inspector {
	return &Inspector{}
}

type Manifest struct {
	Name         string
	Dependencies []base.Dependency
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
	var deps []base.Dependency
	for k, v := range m {
		v := strings.Trim(v, "~^* ")
		if v == "" {
			continue
		}
		deps = append(deps, base.Dependency{
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
			return nil
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

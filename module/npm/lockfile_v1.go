package npm

import (
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/samber/lo"
)

type v1Dep struct {
	Version      string            `json:"version"`
	Dependencies map[string]v1Dep  `json:"dependencies"`
	Dev          bool              `json:"dev"`
	Optional     bool              `json:"optional"`
	Requires     map[string]string `json:"requires"`
}

type v1Lockfile struct {
	Name         string           `json:"name"`
	Dependencies map[string]v1Dep `json:"dependencies"`
}

func processV1Lockfile(data []byte, requires []string) ([]model.DependencyItem, error) {
	var e error
	var lockfile v1Lockfile
	e = json.Unmarshal(data, &lockfile)
	if e != nil {
		return nil, fmt.Errorf("parsing v1 lockfile: bad format, %w", e)
	}
	requires = lo.Uniq(requires)
	var r []model.DependencyItem
	for _, depName := range requires {
		if dep, ok := lockfile.Dependencies[depName]; ok {
			if rr := v1ConvDepRecursive(depName, dep, []map[string]v1Dep{lockfile.Dependencies}, make(map[string]struct{})); rr != nil {
				r = append(r, *rr)
			}
		}
	}
	return r, nil
}

func v1ConvDepRecursive(name string, dep v1Dep, pp []map[string]v1Dep, visited map[string]struct{}) *model.DependencyItem {
	if _, ok := visited[name]; ok {
		return nil
	}
	visited[name] = struct{}{}
	defer func() { delete(visited, name) }()
	r := model.DependencyItem{
		Component: model.Component{
			CompName:    name,
			CompVersion: dep.Version,
			EcoRepo:     EcoRepo,
		},
	}
	if dep.Dev {
		r.IsOnline.SetOnline(false)
	}
	var queryMaps = make([]map[string]v1Dep, len(pp), len(pp)+1)
	copy(queryMaps, pp)
	if dep.Dependencies != nil {
		queryMaps = append(queryMaps, dep.Dependencies)
	}
o:
	for depName := range dep.Requires {
		queryMapOff := len(queryMaps) - 1
		for queryMapOff >= 0 {
			if child, ok := queryMaps[queryMapOff][depName]; ok {
				rr := v1ConvDepRecursive(depName, child, queryMaps[:queryMapOff+1], visited)
				if rr != nil {
					r.Dependencies = append(r.Dependencies, *rr)
				}
				continue o
			}
			queryMapOff--
		}
	}
	return &r
}

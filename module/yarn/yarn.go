package yarn

import "github.com/murphysecurity/murphysec/model"

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

type Dep struct {
	Name     string
	Version  string
	Children []Dep
}

func mapToModel(deps []Dep) []model.DependencyItem {
	var r = make([]model.DependencyItem, len(deps))
	for i := range deps {
		r[i] = model.DependencyItem{
			Component: model.Component{
				CompName:    deps[i].Name,
				CompVersion: deps[i].Version,
				EcoRepo:     EcoRepo,
			},
			Dependencies: mapToModel(deps[i].Children),
		}
	}
	return r
}

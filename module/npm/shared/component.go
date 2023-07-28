package shared

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
)

type Node struct {
	Name     string
	Version  string
	Children []*Node
	model.IsOnline
	Dev bool
}

type dependencyNotFoundError struct {
	name    string
	version string
}

func (c dependencyNotFoundError) Error() string {
	return fmt.Sprintf("dependency not found: %s@%s", c.name, c.version)
}

func CreateDependencyNotFoundError(name, version string) error {
	return &dependencyNotFoundError{name: name, version: version}
}

func ConvNodes(input []*Node) []model.DependencyItem {
	var r = _ConvNodes0(input)
	for i := range r {
		r[i].IsDirectDependency = true
	}
	return r
}

func _ConvNodes0(input []*Node) []model.DependencyItem {
	var r []model.DependencyItem
	for _, node := range input {
		d := model.DependencyItem{
			Component: model.Component{
				CompName:    node.Name,
				CompVersion: node.Version,
				EcoRepo:     EcoRepo,
			},
			Dependencies: _ConvNodes0(node.Children),
			IsOnline:     node.IsOnline,
		}
		if node.Dev {
			d.IsOnline.SetOnline(false)
		}
		r = append(r, d)
	}
	return r
}

var EcoRepo = model.EcoRepo{Ecosystem: "npm"}

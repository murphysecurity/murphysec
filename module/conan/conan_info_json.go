package conan

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"strings"
)

type _InfoJsonItem struct {
	RequiredBy  []string `json:"required_by"`
	DisplayName string   `json:"display_name"`
}

func (t *_ConanInfoJsonFile) ReadFromFile(path string) error {
	data, e := os.ReadFile(path)
	if e != nil {
		return errors.WithCause(ErrReadConanJsonFail, e)
	}
	if e := json.Unmarshal(data, &t); e != nil {
		return errors.WithCause(ErrReadConanJsonFail, e)
	}
	return nil
}

type _ConanInfoJsonFile []_InfoJsonItem

func (t _ConanInfoJsonFile) Tree() (*model.DependencyItem, error) {
	var rootName string
	for _, it := range t {
		if len(it.RequiredBy) == 0 && rootName == "" {
			rootName = it.DisplayName
			break
		}
	}
	if rootName == "" {
		return nil, ErrRootNodeNotFound
	}
	var depGraph = map[string][]string{}
	for _, it := range t {
		for _, requiredBy := range it.RequiredBy {
			depGraph[requiredBy] = append(depGraph[requiredBy], it.DisplayName)
		}
	}
	return _tree(rootName, depGraph, map[string]bool{}), nil
}

func _tree(name string, g map[string][]string, visitedName map[string]bool) *model.DependencyItem {
	if visitedName[name] {
		return nil
	}
	visitedName[name] = true
	defer delete(visitedName, name)
	r := strings.SplitN(name, "/", 2)
	d := &model.DependencyItem{
		Component: model.Component{
			CompName: r[0],
			EcoRepo:  EcoRepo,
		},
	}
	if len(r) > 1 {
		d.CompVersion = r[1]
	}
	for _, it := range g[name] {
		t := _tree(it, g, visitedName)
		if t == nil {
			continue
		}
		d.Dependencies = append(d.Dependencies, *t)
	}
	return d
}

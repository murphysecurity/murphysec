package cargo

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/simpletoml"
)

type cargoLock map[string]cargoItem
type cargoItem struct {
	Version      string
	Dependencies []string
}

func (c cargoLock) findRoot() string {
	var set = map[string]struct{}{}
	for s := range c {
		set[s] = struct{}{}
	}
	for _, i := range c {
		for _, it := range i.Dependencies {
			delete(set, it)
		}
	}
	for s := range set {
		return s
	}
	return ""
}

func parseCargoLock(input []byte) (cargoLock, error) {
	doc, e := simpletoml.UnmarshalTOML(input)
	if e != nil {
		return nil, fmt.Errorf("parseCargoLock: %w", e)
	}
	var rs = map[string]cargoItem{}
	for _, it := range doc.Get("package").TOMLArray() {
		var name = it.Get("name").String()
		if name == "" {
			continue
		}
		item := cargoItem{}
		item.Version = it.Get("version").String()
		var depsDistinctSet = map[string]struct{}{}
		for _, it := range it.Get("dependencies").TOMLArray() {
			if s := it.String(); s != "" {
				if _, ok := depsDistinctSet[s]; ok {
					continue
				}
				depsDistinctSet[s] = struct{}{}
				item.Dependencies = append(item.Dependencies, s)
			}
		}
		rs[name] = item
	}
	return rs, nil
}

func analyzeCargoLock(input []byte) (_ *model.DependencyItem, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("analyzeCargoLock: %w", err)
		}
	}()
	lock, e := parseCargoLock(input)
	if e != nil {
		return nil, e
	}
	rootName := lock.findRoot()
	if rootName == "" {
		return nil, fmt.Errorf("no root found")
	}
	return _buildTree(lock, rootName, map[string]struct{}{}), nil
}

func _buildTree(lock cargoLock, name string, visited map[string]struct{}) *model.DependencyItem {
	if _, ok := visited[name]; ok {
		return nil
	}
	visited[name] = struct{}{}
	defer delete(visited, name)
	item, ok := lock[name]
	if !ok {
		return nil
	}
	r := &model.DependencyItem{
		Component: model.Component{
			CompName:    name,
			CompVersion: item.Version,
			EcoRepo:     EcoRepo,
		},
	}
	for _, depName := range item.Dependencies {
		c := _buildTree(lock, depName, visited)
		if c == nil {
			continue
		}
		r.Dependencies = append(r.Dependencies, *c)
	}
	return r
}

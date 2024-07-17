package cargo

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/simpletoml"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"strings"
)

func splitNameVersionFromDepLine(line string) (name, version string) {
	var i = strings.Index(line, " ")
	if i == -1 {
		name = line
	} else {
		name = line[:i]
		version = line[i+1:]
	}
	return
}

func parseCargoLock2(input []byte) (rs map[[2]string][][2]string, e error) {
	doc, e := simpletoml.UnmarshalTOML(input)
	if e != nil {
		return nil, fmt.Errorf("parseCargoLock: %w", e)
	}
	rs = map[[2]string][][2]string{}
	var singleVersionMap = make(map[string]string)
	var depsMap = make(map[[2]string][]string)
	for _, it := range doc.Get("package").TOMLArray() {
		var name = it.Get("name").String()
		if name == "" {
			continue
		}
		var version = it.Get("version").String()
		if version == "" {
			continue
		}
		var key = [2]string{name, version}
		if _, ok := singleVersionMap[name]; ok {
			singleVersionMap[name] = ""
		} else {
			singleVersionMap[name] = version
		}
		depsMap[key] = lo.Uniq(fp.Map(func(t simpletoml.TOML) string { return t.String() })(it.Get("dependencies").TOMLArray()))
	}
	for key, deps := range depsMap {
		var rDeps [][2]string
		for _, dep := range deps {
			var dn, dv = splitNameVersionFromDepLine(dep)
			if dn == "" {
				continue
			}
			if dv == "" {
				dv = singleVersionMap[dn]
			}
			if dv == "" {
				continue
			}
			rDeps = append(rDeps, [2]string{dn, dv})
		}
		rs[key] = rDeps
	}
	return
}

func findRoots(input map[[2]string][][2]string) (roots [][2]string) {
	var set = make(map[[2]string]struct{})
	for i := range input {
		set[i] = struct{}{}
	}
	for _, deps := range input {
		for _, dep := range deps {
			delete(set, dep)
		}
	}
	return maps.Keys(set)
}

func analyzeCargoLock(input []byte) (rs []*model.DependencyItem, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("analyzeCargoLock: %w", err)
		}
	}()
	lock, e := parseCargoLock2(input)
	if e != nil {
		return nil, e
	}
	roots := findRoots(lock)
	if len(roots) == 0 {
		return nil, fmt.Errorf("no root found")
	}
	for _, root := range roots {
		var r = _buildTree(lock, root, map[[2]string]struct{}{})
		if r == nil {
			continue
		}
		rs = append(rs, r)
	}
	return
}

func _buildTree(lock map[[2]string][][2]string, key [2]string, visited map[[2]string]struct{}) *model.DependencyItem {
	if _, ok := visited[key]; ok {
		return nil
	}
	visited[key] = struct{}{}
	defer delete(visited, key)
	item, ok := lock[key]
	if !ok {
		return nil
	}
	r := &model.DependencyItem{
		Component: model.Component{
			CompName:    key[0],
			CompVersion: key[1],
			EcoRepo:     EcoRepo,
		},
	}
	for _, dep := range item {
		c := _buildTree(lock, dep, visited)
		if c == nil {
			continue
		}
		r.Dependencies = append(r.Dependencies, *c)
	}
	return r
}

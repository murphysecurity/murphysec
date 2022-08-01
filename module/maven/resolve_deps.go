package maven

import (
	"fmt"
	"sort"
	"strings"
)

type DepGraph map[Coordinate]map[Coordinate]struct{}

func (d DepGraph) DOT() string {
	var t []string
	for k, list := range d {
		for v := range list {
			t = append(t, fmt.Sprintf("  \"%s\" -> \"%s\"", k.String(), v.String()))
		}
	}
	sort.Strings(t)
	var s []string
	s = append(s, "digraph dep {")
	s = append(s, t...)
	s = append(s, "}")
	return strings.Join(s, "\n")
}

func (d DepGraph) Tree(root Coordinate) []Dependency {
	return d._tree(root, map[Coordinate]struct{}{})
}

func (d DepGraph) _tree(node Coordinate, visited map[Coordinate]struct{}) []Dependency {
	{
		if visited == nil {
			visited = map[Coordinate]struct{}{}
		}
		if _, ok := visited[node]; ok {
			return nil
		}
		visited[node] = struct{}{}
		defer delete(visited, node)
	}
	var rs []Dependency
	for it := range d[node] {
		rs = append(rs, Dependency{
			Coordinate: it,
			Children:   d._tree(it, visited),
		})
	}
	return rs
}

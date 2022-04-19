package maven

import (
	"fmt"
	"strings"
)

type DepGraph map[Coordinate]map[Coordinate]struct{}

func (d DepGraph) DOT() string {
	var s []string
	s = append(s, "digraph dep {")
	for k, list := range d {
		for v := range list {
			s = append(s, fmt.Sprintf("  \"%s\" -> \"%s\"", k.String(), v.String()))
		}
	}
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

type DepTreeAnalyzer struct {
	graph    map[Coordinate]map[Coordinate]struct{}
	resolver *Resolver
}

func NewDepTreeAnalyzer(resolver *Resolver) *DepTreeAnalyzer {
	return &DepTreeAnalyzer{
		graph:    map[Coordinate]map[Coordinate]struct{}{},
		resolver: resolver,
	}
}

func (d *DepTreeAnalyzer) Resolve(p *PomFile) DepGraph {
	d._resolve(p, nil, 5)
	return d.graph
}

func (d *DepTreeAnalyzer) _resolve(p *PomFile, visited map[Coordinate]struct{}, depth int) {
	if p == nil || !p.coordinate.Complete() || depth < 0 {
		return
	}

	{
		// circle detect
		if visited == nil {
			visited = map[Coordinate]struct{}{}
		}
		if _, ok := visited[p.coordinate]; ok {
			return
		}
		visited[p.coordinate] = struct{}{}
		defer delete(visited, p.coordinate)
	}

	// iterate all dependencies, fetch it, resolve it
	for _, dep := range p.dependencies {
		if d.graph[p.coordinate] == nil {
			d.graph[p.coordinate] = map[Coordinate]struct{}{}
		}
		d.graph[p.coordinate][dep.Coordinate] = struct{}{}
		pf := d.resolver.ResolveByCoordinate(dep.Coordinate)
		if pf == nil {
			continue
		}
		d.graph[p.coordinate][pf.coordinate] = struct{}{}
		d._resolve(pf, visited, depth-1)
	}

}

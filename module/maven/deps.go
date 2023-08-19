package maven

import (
	"sort"
)

type DepsMap struct {
	m map[Coordinate]depsElement
}

func newDepsMap() *DepsMap {
	return &DepsMap{
		m: map[Coordinate]depsElement{},
	}
}

func (d *DepsMap) ListAllEntries() []depsElement {
	var rs []depsElement
	for _, it := range d.m {
		rs = append(rs, it)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].coordinate.Compare(rs[j].coordinate) < 0
	})
	return rs
}

type depsElement struct {
	coordinate   Coordinate
	children     []Dependency
	relativePath string
}

func (d *DepsMap) put(coordinate Coordinate, children []Dependency, path string) {
	d.m[coordinate] = depsElement{
		coordinate:   coordinate,
		children:     children,
		relativePath: path,
	}
}

func (d *DepsMap) allEmpty() bool {
	for _, it := range d.m {
		if len(it.children) > 0 {
			return false
		}
	}
	return true
}

package shared

import (
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"strings"
)

type Visited struct {
	name    string
	version string
	parent  *Visited
}

func (v *Visited) Depth() int {
	var depth = 0
	var curr = v
	for curr != nil {
		curr = curr.parent
		depth++
	}
	return depth
}

func (v *Visited) Contains(name, version string) bool {
	var curr = v
	for curr != nil {
		if curr.name == name && curr.version == version {
			return true
		}
		curr = curr.parent
	}
	return false
}

func (v *Visited) CreateSub(name, version string) *Visited {
	if v.Contains(name, version) {
		return nil
	}
	return &Visited{name: name, version: version, parent: v}
}

func CreateVisited(name, version string) *Visited {
	return &Visited{
		name:    name,
		version: version,
	}
}

type revisitError struct {
	v *Visited
}

func (r revisitError) Error() string {
	var arr [][2]string
	var curr = r.v
	for curr != nil {
		arr = append(arr, [2]string{curr.name, curr.version})
		curr = curr.parent
	}
	lo.Reverse(arr)
	s := strings.Join(fp.Map(func(a [2]string) string { return a[0] + "@" + a[1] })(arr), " -> ")
	return "revisit: " + s
}

func CreateRevisitError(v *Visited) error {
	return revisitError{v: v}
}

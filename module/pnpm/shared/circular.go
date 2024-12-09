package shared

import (
	"github.com/repeale/fp-go"
	"strconv"
	"strings"
)

type CircularDetector struct {
	m    map[[2]string]struct{}
	path [][2]string
}

func (c *CircularDetector) Put(name, version string) {
	if len(c.m) != len(c.path) {
		panic("length inconsistent")
	}
	if c.m == nil {
		c.m = make(map[[2]string]struct{})
	}
	if _, ok := c.m[[2]string{name, version}]; ok {
		panic("circular dependency detected")
	}
	c.path = append(c.path, [2]string{name, version})
	c.m[[2]string{name, version}] = struct{}{}
}

func (c *CircularDetector) Leave(name, version string) {
	if len(c.path) == 0 {
		panic("length underflow")
	}
	if c.path[len(c.path)-1] != [2]string{name, version} {
		panic("name@version inconsistent")
	}
	c.path = c.path[:len(c.path)-1]
}

func (c *CircularDetector) Contains(name, version string) (ok bool) {
	_, ok = c.m[[2]string{name, version}]
	return
}

func (c *CircularDetector) Error(name, version string) error {
	if !c.Contains(name, version) {
		return nil
	}
	var r CircleDetected
	r.path = make([][2]string, len(c.path)+1)
	copy(r.path, c.path)
	r.path[len(c.path)] = [2]string{name, version}
	return r
}

type CircleDetected struct {
	path [][2]string
}

func (e CircleDetected) Error() string {
	var mapper = fp.Map(func(v [2]string) string { return v[0] + "@" + v[1] })
	return "circle detected[length:" + strconv.Itoa(len(e.path)-1) + "]: " + strings.Join(mapper(e.path), " -> ")
}

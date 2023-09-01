package shared

import (
	"github.com/repeale/fp-go"
	"strings"
)

type terror string

func (t terror) Error() string {
	return string(t)
}

const (
	ErrDependencyPath = terror("unprocessable dependency path")
)

type revisitError struct {
	v *Visited
}

func NewRevisitError(v *Visited) error {
	return &revisitError{v: v}
}

func (v revisitError) Error() string {
	var arr [][2]string
	var curr = v.v
	for curr != nil {
		arr = append(arr, [2]string{curr.name, curr.version})
	}
	var s = strings.Join(fp.Map(func(a [2]string) string { return a[0] + "@" + a[1] })(arr), " -> ")
	return "revisit key: " + s
}

type unknownPackageError struct {
	name    string
	version string
}

func (u unknownPackageError) Error() string {
	return "unknown package: " + u.name + "@" + u.version
}

func NewUnknownPackageError(name, version string) error {
	return &unknownPackageError{name, version}
}

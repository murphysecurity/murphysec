package conan

import (
	"github.com/murphysecurity/murphysec/errors"
	"strings"
)

//go:generate stringer -type _e -linecomment -output error_string.go
type _e int

const (
	_                      _e = iota
	ErrConanReport            // conan report error
	ErrConanNotFound          // conan not found
	ErrGetConanVersionFail    // get conan command version failed
	ErrRootNodeNotFound       // can't found root node in conan info graph
	ErrReadConanJsonFail      // read conan json failed
)

func (i _e) Error() string {
	return i.String()
}

type conanError string

func (e conanError) Error() string {
	var rs []string
	var i = -1
	for _, it := range strings.Split(string(e), "\n") {
		s := strings.Trim(it, "\n\r")
		if strings.TrimSpace(s) == "" {
			continue
		}
		if i == -1 && strings.HasPrefix(s, "ERROR:") {
			i = len(rs) // current row number in the result
		}
		rs = append(rs, s)
	}
	if i > -1 {
		rs = rs[i:]
	}
	return strings.Join(rs, "\n")
}

func (e conanError) ErrorMultiLine() []string {
	return strings.Split(e.Error(), "\n")
}

func (e conanError) Is(target error) bool {
	return e == target || errors.Is(target, ErrConanReport)
}

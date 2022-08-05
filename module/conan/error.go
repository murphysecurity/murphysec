package conan

import (
	"github.com/murphysecurity/murphysec/errors"
	"strings"
)

var ErrConanReport = errors.New("conan report error")
var ErrConanNotFound = errors.New("conan not found")
var ErrGetConanVersionFail = errors.New("get conan command version failed")
var ErrRootNodeNotFound = errors.New("can't found root node in conan info graph")
var ErrReadConanJsonFail = errors.New("read conan json failed")

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

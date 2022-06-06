package python

import (
	_ "embed"
	"strings"
)

//go:embed python_pkg_blacklist
var __pyPkgListT string
var pyPkgBlackList = func() map[string]bool {
	m := map[string]bool{}
	for _, s := range strings.Split(__pyPkgListT, "\n") {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "#") {
			continue
		}
		m[s] = true
	}
	return m
}()

package inspector

import (
	_ "embed"
	"path/filepath"
	"strings"
)

//go:embed auto_scan_ignore
var _dirIgnoreText string
var ignoredDirMap = map[string]struct{}{}

func init() {
	for _, s := range strings.Split(_dirIgnoreText, "\n") {
		s := strings.TrimSpace(s)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		ignoredDirMap[s] = struct{}{}
	}
}

func dirShouldIgnore(name string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}
	_, ok := ignoredDirMap[name]
	if ok {
		return true
	}
	_, ok = ignoredDirMap[filepath.Base(name)]
	return ok
}

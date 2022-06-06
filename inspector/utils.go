package inspector

import (
	_ "embed"
	"strings"
)

//go:embed auto_scan_ignore
var _dirIgnoreText string

var dirIgnored = func() func(name string) bool {
	m := map[string]struct{}{}
	for _, it := range strings.Split(_dirIgnoreText, "\n") {
		s := strings.TrimSpace(it)
		if strings.HasPrefix(s, "#") {
			continue
		}
		m[s] = struct{}{}
	}
	return func(name string) bool {
		if strings.HasPrefix(name, ".") {
			return true
		}
		_, ok := m[name]
		return ok
	}
}()

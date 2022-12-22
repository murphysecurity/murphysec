package python

import (
	"regexp"
	"strings"
)

func parseRequirements(data string) map[string]string {
	var pattern = regexp.MustCompile(`^([\w_.-]+)[>=<]+([\w.]+)$`)
	var deps = make(map[string]string)
	for _, s := range strings.Split(data, "\n") {
		s = strings.TrimSpace(s)
		m := pattern.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		k := strings.TrimSpace(m[1])
		if k == "" {
			continue
		}
		v := strings.TrimSpace(m[2])
		deps[k] = v
	}
	return deps
}

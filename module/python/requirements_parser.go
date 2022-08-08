package python

import (
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"strings"
)

func readRequirements(path string) ([]model.Dependency, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return nil, errors.Wrap(e, "read requirements failed")
	}
	return parseRequirements(string(data)), nil
}

func parseRequirements(data string) []model.Dependency {
	var deps []model.Dependency
	lines := strings.Split(data, "\n")
	for _, it := range lines {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		rs := strings.SplitN(it, "=", 2)
		rs[0] = strings.TrimSpace(strings.TrimRight(rs[0], ">=<"))
		if len(rs) > 1 {
			rs[1] = strings.TrimSpace(strings.TrimLeft(rs[1], ">=<"))
		}
		deps = append(deps, model.Dependency{
			Name:         rs[0],
			Version:      rs[1],
			Dependencies: nil,
		})
	}
	return deps
}

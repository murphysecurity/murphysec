package python

import (
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"regexp"
	"strings"
)

func readRequirements(path string) ([]model.DependencyItem, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return nil, errors.Wrap(e, "read requirements failed")
	}
	return parseRequirements(string(data)), nil
}

func parseRequirements(data string) []model.DependencyItem {
	var pattern = regexp.MustCompile(`^([\w_.-]+)[>=<]+([\w.]+)$`)
	var deps []model.DependencyItem
	for _, s := range strings.Split(data, "\n") {
		s = strings.TrimSpace(s)
		m := pattern.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		var di model.DependencyItem
		di.CompName = m[1]
		di.CompVersion = m[2]
		di.EcoRepo = EcoRepo
		deps = append(deps, di)
	}
	return deps
}

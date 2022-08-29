package python

import (
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"regexp"
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
	var pattern = regexp.MustCompile("^([\\w_-]+)[>=<]+([\\w.]+)$")
	var deps []model.Dependency
	for _, s := range strings.Split(data, "\n") {
		s = strings.TrimSpace(s)
		m := pattern.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		deps = append(deps, model.Dependency{
			Name:    m[1],
			Version: m[2],
		})
	}
	return deps
}

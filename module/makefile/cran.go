package makefile

import (
	"github.com/murphysecurity/murphysec/model"
	"regexp"
	"strings"
)

var makefileLineBreak = regexp.MustCompile("\\\r?\n")

func findingCRANItems(makefileData string) (rs []model.Dependency) {
	var pattern = regexp.MustCompile("R(?:script|\\s+).+-e.+[\"'].*install.packages\\(\\s*[\"'](\\w+)[\"']")
	for _, s := range strings.Split(makefileLineBreak.ReplaceAllString(makefileData, ""), "\n") {
		s = strings.TrimSpace(s)
		m := pattern.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		rs = append(rs, model.Dependency{
			Name: m[1],
		})
	}
	return
}

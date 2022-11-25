package gradle

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type GradleVersion struct {
	Version string            `json:"version"`
	Items   map[string]string `json:"items,omitempty"`
}

func (g *GradleVersion) String() string {
	var items []string
	for k, v := range g.Items {
		items = append(items, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(items)
	var rs = []string{fmt.Sprintf("Gradle: %s", g.Version)}
	rs = append(rs, items...)
	return strings.Join(rs, ", ")
}

func parseGradleVersionOutput(input string) (*GradleVersion, error) {
	var v = &GradleVersion{
		Version: "",
		Items:   map[string]string{},
	}
	var gvPattern = regexp.MustCompile(`(?m)^Gradle\s+([\w.-]+)`)
	var itemPattern = regexp.MustCompile(`(?m)^([\w -]+):\s+(.+)`)
	if m := gvPattern.FindStringSubmatch(input); m != nil {
		v.Version = strings.TrimSpace(m[1])
	}
	for _, g := range itemPattern.FindAllStringSubmatch(input, -1) {
		var name = g[1]
		var version = strings.TrimSpace(g[2])
		v.Items[name] = version
	}
	if v.Version == "" {
		return nil, fmt.Errorf("parse gradle version failed")
	}
	return v, nil
}

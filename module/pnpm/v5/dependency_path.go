package v5

import (
	"regexp"
	"strings"
)

var __pathParenthesesPattern = regexp.MustCompile(`(?:\([^)]+\)\s*)+$`)

func getNameVersionFromPath0(path string) (string, string) {
	path = __pathParenthesesPattern.ReplaceAllString(path, "") // remove all parentheses
	path = strings.TrimSuffix(path, "/")
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return "", ""
	}
	var version = strings.Trim(path[i:], "/")
	version = strings.SplitN(version, "_", 2)[0]
	var name = path[:i]
	name = strings.Trim(name, "/")
	return name, version
}

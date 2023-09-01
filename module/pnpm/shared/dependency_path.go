package shared

import (
	"regexp"
	"strings"
)

var __pathParenthesesPattern = regexp.MustCompile(`(?:\([^)]+\)\s*)+$`)

func trimParenthesesFromPath(path string) string {
	path = strings.TrimSpace(path)
	path = __pathParenthesesPattern.ReplaceAllString(path, "")
	return path
}

func GetNameFromPath(path string) (string, error) {
	path = trimParenthesesFromPath(path)
	if i := strings.LastIndex(path, "/"); i > -1 {
		if i+1 >= len(path) {
			return "", ErrDependencyPath
		}
		path = path[i+1:]
	}
	if i := strings.Index(path, "@"); i > -1 {
		path = path[:i]
	}
	return path, nil
}

func GetVersionFromPath(path string) (string, error) {
	path = trimParenthesesFromPath(path)
	if i := strings.LastIndex(path, "/"); i > -1 {
		if i+1 >= len(path) {
			return "", ErrDependencyPath
		}
		path = path[i+1:]
	}
	if i := strings.Index(path, "@"); i > -1 {
		if i+1 >= len(path) {
			return "", ErrDependencyPath
		}
		path = path[i+1:]
	}
	return path, nil
}

func ignoreError[T any](a T, _ error) T {
	return a
}

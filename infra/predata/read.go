package predata

import "strings"

func ParseString(s string) []string {
	var r []string
	for _, it := range strings.Split(s, "\n") {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		if it[0] == '#' {
			continue
		}
		r = append(r, it)
	}
	return r
}

func StringsToMapBool(s []string) map[string]bool {
	var r = map[string]bool{}
	for _, it := range s {
		r[it] = true
	}
	return r
}

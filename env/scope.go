package env

import "strings"

var Scope string

type ScopeSet map[string]struct{}

func (s ScopeSet) Has(scope string) bool {
	if scope == "" {
		return true
	}
	if s == nil {
		return false
	}
	if _, ok := s["all"]; ok {
		return true
	}
	_, ok := s[scope]
	return ok
}

func GetScanScopes() ScopeSet {
	var validScopes = map[string]bool{
		"compile":  true,
		"provided": true,
		"runtime":  true,
		"test":     true,
		"system":   true,
		"all":      true,
	}
	var rs = ScopeSet{}
	for _, s := range strings.Split(Scope, ",") {
		s = strings.TrimSpace(s)
		if _, ok := validScopes[s]; !ok {
			continue
		}
		rs[s] = struct{}{}
	}
	if len(rs) == 0 {
		rs["compile"] = struct{}{}
		rs["runtime"] = struct{}{}
	}

	return rs
}

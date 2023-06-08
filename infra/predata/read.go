package predata

import (
	"github.com/murphysecurity/murphysec/infra/sl"
	"github.com/repeale/fp-go"
	"strings"
)

func ParseString(s string) []string {
	return fp.Pipe2(fp.Map(strings.TrimSpace), fp.Filter(sl.NotF1(lineShouldSkip)))(strings.Split(s, "\n"))
}

func lineShouldSkip(s string) bool {
	return s == "" || s[0] == '#'
}

func StringsToMapBool(s []string) map[string]bool {
	var r = map[string]bool{}
	for _, it := range s {
		r[it] = true
	}
	return r
}

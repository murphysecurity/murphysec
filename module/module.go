package module

import (
	"github.com/murphysecurity/murphysec/module/base"
	"sort"
)

var Inspectors []base.Inspector

func GetSupportedModuleList() []string {
	var r []string
	for _, it := range Inspectors {
		r = append(r, it.String())
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

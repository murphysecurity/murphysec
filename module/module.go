package module

import (
	"github.com/murphysecurity/murphysec/model"
	"sort"
)

var Inspectors []model.Inspector

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

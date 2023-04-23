package inspector

import (
	"context"
	"fmt"
	"github.com/murphysecurity/fix-tools/fix"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
)

func scanFragment(ctx context.Context, dir string, components []model.Component) ([]model.ComponentCodeFragment, error) {
	components = utils.DistinctSlice(components)
	if len(components) == 0 {
		return make([]model.ComponentCodeFragment, 0), nil
	}
	var result = make([]model.ComponentCodeFragment, 0)
	for _, component := range components {
		var param = fix.FixParams{
			ShowOnly: true,
			CompList: []fix.Comp{{
				CompName:    component.CompName,
				CompVersion: component.CompVersion,
			}},
			PackageManager: "maven",
			Dir:            dir,
		}
		previews, e := param.Fix()
		if e != nil {
			return nil, fmt.Errorf("scan fragment: %w", e)
		}
		var r = model.ComponentCodeFragment{
			Component:     component,
			CodeFragments: utils.NoNilSlice(previews),
		}
		result = append(result, r)
	}

	return result, nil
}

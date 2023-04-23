package inspector

import (
	"context"
	"fmt"
	"github.com/murphysecurity/fix-tools/fix"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
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
			CodeFragments: make([]model.CodeFragment, 0),
		}
		for _, it := range previews {
			if len(it.Content) == 0 {
				continue
			}
			if filepath.IsAbs(it.Path) {
				panic("is abs path")
			}
			var fp = filepath.ToSlash(it.Path)

			var t string
			var lineBegin = it.Content[0].Line
			for _, content := range it.Content {
				t = t + "\n" + content.Text
			}
			r.CodeFragments = append(r.CodeFragments, model.CodeFragment{
				Text:         t,
				LineBegin:    lineBegin,
				RelativePath: fp,
			})
		}
	}

	return result, nil
}

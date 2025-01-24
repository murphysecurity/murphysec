package inspector

import (
	"context"
	"github.com/murphysecurity/fix-tools/fix"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
)

func scanFragment(ctx context.Context, dir string, components []model.Component) ([]model.ComponentCodeFragment, error) {
	logger := logctx.Use(ctx)
	components = utils.DistinctSlice(components)
	var mapper = fp.Map(func(it model.Component) fix.Comp { return fix.Comp{CompName: it.CompName, CompVersion: it.CompVersion} })
	var filter = func(eco string) func(it model.Component) bool {
		return func(it model.Component) bool { return it.Ecosystem == eco }
	}
	var mavenComponents = fp.Pipe2(fp.Filter(filter("maven")), mapper)(components)
	if len(mavenComponents) == 0 {
		return make([]model.ComponentCodeFragment, 0), nil
	}
	var param = fix.FixParams{
		ShowOnly:       true,
		CompList:       mavenComponents,
		PackageManager: "maven_cli",
		RepoType:       "local",
		Dir:            dir,
	}
	logger.Sugar().Debugf("scan fragment: %d", len(components))
	xdResult := param.Fix()
	logger.Sugar().Debugf("scan fragment: completed")
	var groupedXResult = lo.GroupBy(xdResult.Preview, func(it fix.Preview) [2]string { return [2]string{it.CompName, it.CompVersion} })
	var result = make([]model.ComponentCodeFragment, 0)
	for _, it := range components {
		var key = [2]string{it.CompName, it.CompVersion}
		var codeFragment = groupedXResult[key]
		var response = fix.Response{Err: nil, Preview: codeFragment}
		result = append(result, model.ComponentCodeFragment{Component: it, CodeFragmentResult: response})
	}
	return result, nil
}

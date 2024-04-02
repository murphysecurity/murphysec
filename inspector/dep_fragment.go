package inspector

import (
	"context"
	"encoding/json"
	"github.com/murphysecurity/fix-tools/fix"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"runtime"
	"sync"
)

func scanFragment(ctx context.Context, dir string, components []model.Component) ([]model.ComponentCodeFragment, error) {
	logger := logctx.Use(ctx)
	components = utils.DistinctSlice(components)
	if len(components) == 0 {
		return make([]model.ComponentCodeFragment, 0), nil
	}
	var result = make([]model.ComponentCodeFragment, 0)
	var wg sync.WaitGroup
	var componentCh = make(chan model.Component, len(components))
	for _, it := range components {
		if it.Ecosystem != "maven" {
			continue
		}
		componentCh <- it
	}
	close(componentCh)
	var resultCh = make(chan model.ComponentCodeFragment, len(components))
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer func() { wg.Done() }()
			for component := range componentCh {
				var param = fix.FixParams{
					ShowOnly: true,
					CompList: []fix.Comp{{
						CompName:    component.CompName,
						CompVersion: component.CompVersion,
					}},
					PackageManager: "maven",
					RepoType:       "local",
					Dir:            dir,
				}
				logger.Sugar().Debugf("fix: %s", string(must.A(json.Marshal(param))))
				xdResult := param.Fix()
				if xdResult.Err != nil {
					logger.Sugar().Errorf("scan fragment: %s", xdResult.Err)
				}
				resultCh <- model.ComponentCodeFragment{
					Component:          component,
					CodeFragmentResult: xdResult,
				}
			}
		}()
	}

	wg.Wait()
	close(resultCh)
	for ch := range resultCh {
		result = append(result, ch)
	}
	return result, nil
}

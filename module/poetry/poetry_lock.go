package poetry

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pelletier/go-toml/v2"
)

func parsePoetryLock(ctx context.Context, f string) (rs []model.DependencyItem, e error) {
	var data []byte
	var logger = logctx.Use(ctx).Sugar()
	data, e = utils.ReadFileLimited(f, 4*1024*1024)
	if e != nil {
		logger.Warnf("Read file failed. %v %v", e, f)
		return nil, e
	}
	root := &tomlTree{}
	if e := toml.Unmarshal(data, &root); e != nil {
		logger.Warnf("Parse toml failed. %v %v", e.Error(), f)
		return nil, e
	}
	for _, it := range root.Get("package").AsArray() {
		rs = append(rs, model.DependencyItem{
			Component: model.Component{
				CompName:    it.Get("name").String(),
				CompVersion: it.Get("version").String(),
				EcoRepo:     EcoRepo,
			},
		})
	}
	logger.Infof("Parse poetry.lock, found %d", len(rs))
	return
}

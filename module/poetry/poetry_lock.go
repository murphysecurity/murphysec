package poetry

import (
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pelletier/go-toml/v2"
)

func parsePoetryLock(f string) (rs []model.Dependency, e error) {
	var data []byte
	data, e = utils.ReadFileLimited(f, 4*1024*1024)
	if e != nil {
		logger.Warn.Println("Read file failed.", e.Error(), f)
		return nil, e
	}
	root := &tomlTree{}
	if e := toml.Unmarshal(data, &root); e != nil {
		logger.Warn.Println("Parse toml failed.", e.Error(), f)
		return nil, e
	}
	for _, it := range root.Get("package").AsArray() {
		rs = append(rs, model.Dependency{
			Name:    it.Get("name").String(),
			Version: it.Get("version").String(),
		})
	}
	logger.Info.Println("Parse poetry.lock, found", len(rs))
	return
}

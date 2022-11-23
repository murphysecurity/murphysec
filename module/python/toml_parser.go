package python

import (
	"context"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/simpletoml"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.uber.org/zap"
	"golang.org/x/mod/semver"
	"regexp"
)

var ErrParseToml = errors.New("Parse toml failed")

func tomlBuildSysFile(ctx context.Context, path string) ([]model.DependencyItem, error) {
	logger := utils.UseLogger(ctx)
	logger.Debug("Process toml buildSys file", zap.String("path", path))
	data, e := utils.ReadFileLimited(path, 4*1024*1024)
	if e != nil {
		return nil, e
	}
	return tomlBuildSys(data)
}

func tomlBuildSys(data []byte) ([]model.DependencyItem, error) {
	// numpy==1.13.3 Cython>=0.29.13 wheel
	pa := regexp.MustCompile("([\\w.-]+)(?:[>=]?=([\\w.-]+))?")
	t, e := simpletoml.UnmarshalTOML(data)
	if e != nil {
		return nil, errors.WithCause(ErrParseToml, e)
	}
	rsm := orderedmap.New[string, string]()
	for _, it := range t.Get("build-system", "requires").TOMLArray() {
		m := pa.FindStringSubmatch(it.String(""))
		if m == nil {
			continue
		}
		name := m[1]
		version := m[2]
		if oldVer, ok := rsm.Get(name); ok && oldVer != "" {
			if version == "" {
				continue
			}
			if semver.IsValid(version) && semver.Compare(oldVer, version) < 0 {
				rsm.Set(name, version)
			}
		} else {
			rsm.Set(name, version)
		}
	}

	r := make([]model.DependencyItem, 0)
	for pair := rsm.Oldest(); pair != nil; pair = pair.Next() {
		var di model.DependencyItem
		di.CompName = pair.Key
		di.CompVersion = pair.Value
		di.EcoRepo = EcoRepo
		r = append(r, di)
	}
	return r, nil
}

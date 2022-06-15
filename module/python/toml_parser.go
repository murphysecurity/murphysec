package python

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/simpletoml"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"regexp"
)

var ErrParseToml = errors.New("Parse toml failed")

func tomlBuildSysFile(path string) ([]model.Dependency, error) {
	data, e := utils.ReadFileLimited(path, 4*1024*1024)
	if e != nil {
		return nil, e
	}
	return tomlBuildSys(data)
}

func tomlBuildSys(data []byte) ([]model.Dependency, error) {
	// numpy==1.13.3 Cython>=0.29.13 wheel
	pa := regexp.MustCompile("([\\w.-]+)(?:[>=]?=([\\w.-]+))?")
	t, e := simpletoml.UnmarshalTOML(data)
	if e != nil {
		return nil, errors.Wrap(ErrParseToml, e.Error())
	}
	rsm := map[string]string{}
	for _, it := range t.Get("build-system", "requires").TOMLArray() {
		m := pa.FindStringSubmatch(it.String(""))
		if m == nil {
			continue
		}
		name := m[1]
		version := m[2]
		if oldVer, ok := rsm[name]; ok && oldVer != "" {
			if version == "" {
				continue
			}
			if semver.IsValid(version) && semver.Compare(oldVer, version) < 0 {
				rsm[name] = version
			}
		} else {
			rsm[name] = version
		}
	}

	r := make([]model.Dependency, 0)
	for name, version := range rsm {
		r = append(r, model.Dependency{
			Name:    name,
			Version: version,
		})
	}
	return r, nil
}

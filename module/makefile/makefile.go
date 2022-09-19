package makefile

import (
	"context"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "MakefileInspector"
}

func (i Inspector) CheckDir(dir string) bool {
	return utils.IsFile("Makefile")
}

func (i Inspector) InspectProject(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (i Inspector) SupportFeature(feature base.Feature) bool {
	switch feature {
	case base.FeatureAllowNested:
		return true
	}
	return false
}

var Instance base.Inspector = &Inspector{}

package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/url"
	"testing"
)

func TestReadLocalProject(t *testing.T) {
	modules, e := ReadLocalProject(context.TODO(), "./__test/multi_module")
	assert.NoError(t, e)
	assert.EqualValues(t, "1.0.0-SNAPSHOT", modules[1].ParentCoordinate().Version)
	assert.EqualValues(t, "1.0.0-SNAPSHOT", modules[0].Project.Version)
	assert.EqualValues(t, 2, len(modules))
}

func TestResolve(t *testing.T) {
	logger := must.A(zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel)))
	ctx := utils.WithLogger(context.TODO(), logger)

	modules := must.A(ReadLocalProject(ctx, "./__test/multi_module"))
	resolver := NewPomResolver(ctx)
	//resolver.AddRepo(NewHttpRepo(ctx, *must.A(url.Parse("https://repo1.maven.org/maven2/"))))
	resolver.AddRepo(NewHttpRepo(ctx, *must.A(url.Parse("https://maven.aliyun.com/repository/public"))))
	for _, module := range modules {
		resolver.pomCache.add(module)
	}
	for _, module := range modules {
		_ = must.A(resolver.ResolvePom(ctx, module.Coordinate()))
	}
	BuildDepTree(ctx, resolver, modules[1].Coordinate())
}

package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var _prepared = false

func prepareTest() {
	if _prepared {
		return
	}
	e := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".noext" {
			_ = os.Link(path, strings.TrimSuffix(path, ".noext"))
		}
		return nil
	})
	if e != nil {
		panic(e)
	}
	_prepared = true
}

func TestReadLocalProject(t *testing.T) {
	prepareTest()
	modules, e := ReadLocalProject(context.TODO(), "./__test/multi_module")
	assert.NoError(t, e)
	assert.EqualValues(t, "1.0.0-SNAPSHOT", modules[1].ParentCoordinate().Version)
	assert.EqualValues(t, "1.0.0-SNAPSHOT", modules[0].Project.Version)
	assert.EqualValues(t, 2, len(modules))
}

func TestResolve(t *testing.T) {
	prepareTest()
	mavenRepo := os.Getenv("DEFAULT_MAVEN_REPO")
	if mavenRepo == "" && os.Getenv("CI") != "" {
		t.Skip("Currently in CI environment, the environment variable DEFAULT_MAVEN_REPO not set, skip test")
		return
	}
	if mavenRepo == "" {
		mavenRepo = "https://maven.aliyun.com/repository/public"
	}
	logger := must.A(zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel)))
	ctx := logctx.With(context.TODO(), logger)

	modules := must.A(ReadLocalProject(ctx, "./__test/multi_module"))
	resolver := NewPomResolver(ctx, []M2Remote{newHttpRemote(*must.A(url.Parse(mavenRepo)))})
	for _, module := range modules {
		resolver.addPom(module)
	}
	//for _, module := range modules {
	//	_ = must.A(resolver.ResolvePom(ctx, module.Coordinate()))
	//}
	rp := must.A(resolver.ResolvePom(ctx, modules[1].Coordinate()))
	r := BuildDepTree(ctx, resolver, modules[1].Coordinate())
	fmt.Println(r, rp)
}

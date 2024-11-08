package v9

import (
	"context"
	"embed"
	_ "embed"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/repeale/fp-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"testing"
)

//go:embed testdata/*
var data embed.FS

func doTest0(reader io.Reader, t *testing.T) struct{} {
	tree, e := Parse(logctx.With(context.Background(), must.A(zap.NewDevelopment())), reader)
	assert.NoError(t, e)
	assert.NotNil(t, tree)
	return struct{}{}
}

func void[T, R any](f func(T) R) func(T) {
	return func(t T) { f(t) }
}

func TestParse(t *testing.T) {
	entries, e := data.ReadDir("testdata")
	assert.NoError(t, e)
	assert.NotEmpty(t, entries)
	e = fs.WalkDir(data, ".", func(path string, d fs.DirEntry, e error) error {
		assert.NoError(t, e)
		if d.IsDir() {
			return nil
		}
		t.Log(path)
		f, e := data.Open(path)
		assert.NoError(t, e)
		t.Run(path, void(fp.Curry2(doTest0)(f)))
		e = f.Close()
		assert.NoError(t, e)
		return e
	})
	assert.NoError(t, e)
}

func Test_splitNameVersion(t *testing.T) {
	n, v, _ := splitNameVersion("typescript-eslint@8.3.0(eslint@9.9.1(jiti@1.21.6))(typescript@5.5.4)")
	assert.Equal(t, "typescript-eslint", n)
	assert.Equal(t, "8.3.0", v)
}

func Test_splitRealVersion(t *testing.T) {
	assert.Equal(t, "8.3.0", splitRealVersionInVersionString("8.3.0(eslint@9.9.1(jiti@1.21.6))(typescript@5.5.4)"))
}

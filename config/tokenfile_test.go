package config

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestTokenFile(t *testing.T) {
	p, e := homedir.Expand(DefaultTokenFile)
	if e != nil {
		panic(e)
	}
	if utils.IsFile(p) {
		t.SkipNow()
	}
	var logger, _ = zap.NewDevelopment()
	var ctx = logctx.With(context.TODO(), logger)
	var (
		token string
	)
	_, e = ReadTokenFile(ctx)
	assert.ErrorIs(t, e, ErrNoToken)
	assert.NoError(t, WriteLocalTokenFile(ctx, "foo"))
	token, e = ReadTokenFile(ctx)
	assert.NoError(t, e)
	assert.Equal(t, "foo", token)
	assert.NoError(t, RemoveTokenFile(ctx))
}

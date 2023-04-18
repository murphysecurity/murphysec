package config

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestRepoConfig(t *testing.T) {
	var ctx = logctx.With(context.TODO(), must.A(zap.NewDevelopment()))
	assert.NoError(t, WriteRepoConfig(ctx, ".", model.AccessTypeCli, RepoConfig{TaskId: "111"}))
	co, e := ReadRepoConfig(ctx, ".", model.AccessTypeCli)
	assert.NoError(t, e)
	assert.Equal(t, "111", co.TaskId)
}

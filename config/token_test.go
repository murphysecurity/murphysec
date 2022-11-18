package config

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestReadToken(t *testing.T) {
	s, e := GetToken(logctx.With(context.TODO(), must.M1(zap.NewDevelopment())))
	if e == nil {
		assert.NotEmpty(t, s)
	} else {
		assert.ErrorIs(t, e, ErrNoToken)
	}
}

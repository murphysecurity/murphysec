package cpphasher

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"path/filepath"
	"testing"
)

func TestMD5HashingCppFiles(t *testing.T) {
	// work around
	logger := must.M1(zap.NewDevelopment())
	ctx := logctx.With(context.TODO(), logger)
	cppFileExtSet[".go"] = true
	defer delete(cppFileExtSet, ".go")
	h, e := MD5HashingCppFiles(ctx, must.M1(filepath.Abs(".")))
	assert.NoError(t, e)
	assert.True(t, len(h) > 0)
}

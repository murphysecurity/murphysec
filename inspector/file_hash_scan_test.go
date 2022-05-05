package inspector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileHashScan(t *testing.T) {
	ctx := &ScanContext{
		ProjectDir: ".",
	}
	CxxExtSet[".go"] = true
	assert.NoError(t, FileHashScan(ctx))
	assert.True(t, len(ctx.FileHashes) > 0)
}

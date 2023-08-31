package v5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLockfile(t *testing.T) {
	for i, s := range testDataList {
		lockfile, e := ParseLockfile([]byte(s))
		assert.NoError(t, e, i)
		assert.NotNil(t, lockfile, i)
	}
}

func TestAnalyzeDepTree(t *testing.T) {
	for _, s := range testDataList {
		lockfile, e := ParseLockfile([]byte(s))
		assert.NoError(t, e)
		trees := AnalyzeDepTree(lockfile)
		t.Log(len(trees))
		//t.Log(trees)
	}
}

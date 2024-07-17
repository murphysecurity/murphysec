package cargo

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed cargo-lock-test.toml
var __cargo_lock_test []byte

func Test_analyzeCargoLock(t *testing.T) {
	tree, e := analyzeCargoLock(__cargo_lock_test)
	assert.Equal(t, 1, len(tree))
	assert.NoError(t, e)
	for i := 0; i < 20; i++ {
		t2, e := analyzeCargoLock(__cargo_lock_test)
		assert.NoError(t, e)
		assert.Equal(t, tree, t2)
	}
}

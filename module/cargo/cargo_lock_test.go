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
	assert.NoError(t, e)
	assert.NotNil(t, tree)
	tree2, _ := analyzeCargoLock(__cargo_lock_test)
	assert.Equal(t, tree, tree2)
}

func Test_buildTree(t *testing.T) {
	m := cargoLock{
		"a": cargoItem{
			Version:      "",
			Dependencies: []string{"a"},
		},
	}
	assert.NotNil(t, _buildTree(m, "a", map[string]struct{}{}))
}

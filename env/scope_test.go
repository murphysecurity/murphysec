package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetScanScopes(t *testing.T) {
	GetScanScopes()
	Scope = "s, test, s"
	scope := GetScanScopes()
	assert.True(t, scope.Has("test"))
	assert.True(t, scope.Has(""))
	assert.False(t, scope.Has("foo"))
	Scope = "s, ,all"
	scope = GetScanScopes()
	assert.True(t, scope.Has("foo"))
}

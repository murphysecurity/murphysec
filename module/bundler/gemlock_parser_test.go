package bundler

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test_gemlock
var testGemLock string

func TestParseGemLock(t *testing.T) {
	tree, e := parseGemLock(testGemLock)
	assert.NoError(t, e)
	assert.EqualValues(t, 3, len(tree.get("GIT").children))
	assert.EqualValues(t, "GIT", tree.get("GIT").line)
}

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

func Test_getDepGraph(t *testing.T) {
	var data = `
GEM
  remote: http://rubygems.org/
  specs:
    rake (12.3.3)
    test-unit (2.5.5)

PLATFORMS
  ruby

DEPENDENCIES
  rake
  test-unit (~> 2.4)
`
	tree, e := getDepGraph(data)
	assert.NoError(t, e)
	t.Log(tree)
}

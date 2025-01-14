package bundler

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
func TestParseGem(t *testing.T) {
	var m gemfile
	file, _ := os.Open("E:\\Desktop\\kubernetes-1.22.1\\cluster\\addons\\fluentd-elasticsearch\\fluentd-es-image\\Gemfile")
	m.Parse(file)
	for _, m := range m.Gems {
		fmt.Printf("%s :%s \n", m.Name, m.Version)
	}
}

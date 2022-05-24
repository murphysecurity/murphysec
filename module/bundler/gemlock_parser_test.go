package bundler

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test_gemlock
var testGemLock string

func TestParseGemLock(t *testing.T) {
	tree, e := parseGemLock(testGemLock)
	assert.NoError(t, e)
	fmt.Println(tree.get("GIT"))
}

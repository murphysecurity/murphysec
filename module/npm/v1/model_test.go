package v1

import (
	_ "embed"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed v1_testdata.json
var testdata1 []byte

func TestUnmarshal(t *testing.T) {
	var root lockRoot
	assert.NoError(t, json.Unmarshal(testdata1, &root))
	postprocessPkg(&root.lockPkg, nil)
	t.Log(1)
}

func TestParseLockfile(t *testing.T) {
	var requires = [][2]string{
		{"firebase-admin", "^8.10.0"},
		{"firebase-functions", "^3.6.1"},
		{"firebase-functions-test", "^0.2.0"},
	}
	r, e := ParseLockfile(testdata1)
	assert.NoError(t, e)
	n, e := r.Build(requires, false)
	assert.NoError(t, e)
	t.Log(n)
}
func TestParseLockfileStrictly(t *testing.T) {
	t.Skip("skip strict parse due a circular dependency, need more investigation")
	var requires = [][2]string{
		{"firebase-admin", "^8.10.0"},
		{"firebase-functions", "^3.6.1"},
		{"firebase-functions-test", "^0.2.0"},
	}
	r, e := ParseLockfile(testdata1)
	assert.NoError(t, e)
	n, e := r.Build(requires, true)
	assert.NoError(t, e)
	t.Log(n)
}

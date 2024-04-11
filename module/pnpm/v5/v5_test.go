package v5

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
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

func TestBuildDepTree(t *testing.T) {
	l, _ := ParseLockfile([]byte(testDataList[4]))
	tree := BuildDepTree(l, nil, "")
	assert.NotNil(t, tree)
	t.Log(string(must.A(json.MarshalIndent(tree, "| ", " "))))
}

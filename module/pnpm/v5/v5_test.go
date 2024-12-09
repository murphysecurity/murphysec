package v5

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestLockfile(t *testing.T) {
	files, _ := testFiles.ReadDir("testdata")
	assert.NotEmpty(t, files)
	for _, s := range files {
		t.Run(s.Name(), func(t *testing.T) {
			f, e := testFiles.Open("testdata/" + s.Name())
			assert.NoError(t, e)
			defer func() { assert.NoError(t, f.Close()) }()
			data, e := io.ReadAll(f)
			assert.NoError(t, e)
			lockfile, e := ParseLockfile(data)
			assert.NoError(t, e)
			assert.NotNil(t, lockfile)

		})
	}
}

func TestBuildDepTree(t *testing.T) {

	l, _ := ParseLockfile([]byte(testData5))
	tree := BuildDepTree(l, nil, "")
	assert.NotNil(t, tree)
	t.Log(string(must.A(json.MarshalIndent(tree, "| ", " "))))
}

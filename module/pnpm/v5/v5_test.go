package v5

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
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

func TestParseLockfile_ConsistentAcrossMultipleRuns(t *testing.T) {
	files, _ := testFiles.ReadDir("testdata")
	assert.NotEmpty(t, files)
	for _, s := range files {
		t.Run(s.Name(), func(t *testing.T) {
			f, e := testFiles.Open("testdata/" + s.Name())
			assert.NoError(t, e)
			defer func() { assert.NoError(t, f.Close()) }()
			data, e := io.ReadAll(f)
			assert.NoError(t, e)

			var baseline []byte
			for i := 0; i < 10; i++ {
				lockfile, e := ParseLockfile(data)
				assert.NoError(t, e)
				assert.NotNil(t, lockfile)
				got, e := json.Marshal(lockfile)
				assert.NoError(t, e)
				if i == 0 {
					baseline = got
					continue
				}
				assert.Equal(t, string(baseline), string(got), "parse result changed at run %d", i+1)
			}
		})
	}
}

func TestBuildDepTree_AllTestdata(t *testing.T) {
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

			tree := BuildDepTree(lockfile, nil, "")
			if len(lockfile.Dependencies) > 0 || len(lockfile.DevDependencies) > 0 {
				assert.NotNil(t, tree)
			}

			for importerName, importer := range lockfile.Importers {
				importerTree := BuildDepTree(lockfile, importer, importerName)
				if importerTree != nil {
					assert.Equal(t, importerName, importerTree.Name)
				}
			}
		})
	}
}

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

func TestParseNameVersionFromPackagePath(t *testing.T) {
	testCases := []struct {
		path    string
		name    string
		version string
		ok      bool
	}{
		{
			path:    "/@babel/helper-compilation-targets/7.21.4_@babel+core@7.12.10",
			name:    "@babel/helper-compilation-targets",
			version: "7.21.4",
			ok:      true,
		},
		{
			path:    "/foo/1.2.3_a@1.0.0_b@2.0.0",
			name:    "foo",
			version: "1.2.3",
			ok:      true,
		},
		{
			path:    "/@types/babel__core/7.20.1",
			name:    "@types/babel__core",
			version: "7.20.1",
			ok:      true,
		},
		{
			path: "/bad",
			ok:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			name, version, ok := parseNameVersionFromPackagePath(tc.path)
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.name, name)
			assert.Equal(t, tc.version, version)
		})
	}
}

func TestBuildIndexes_MergeMultipleContexts(t *testing.T) {
	l := &Lockfile{
		Packages: map[string]*Pkg{
			"/foo/1.2.3_peer-a@1.0.0": {
				Dependencies: map[string]string{
					"a": "1.0.0",
				},
				Dev: true,
			},
			"/foo/1.2.3_peer-b@2.0.0": {
				Dependencies: map[string]string{
					"b": "2.0.0",
				},
				Dev: false,
			},
		},
	}

	l.buildIndexes()

	pkg := l.pkgIndexes[[2]string{"foo", "1.2.3"}]
	assert.NotNil(t, pkg)
	assert.Equal(t, "foo", pkg.Name)
	assert.Equal(t, "1.2.3", pkg.Version)
	assert.Equal(t, map[string]string{
		"a": "1.0.0",
		"b": "2.0.0",
	}, pkg.Dependencies)
	assert.False(t, pkg.Dev)
}

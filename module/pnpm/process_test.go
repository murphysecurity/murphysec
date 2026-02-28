package pnpm

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed v5/testdata/*.yaml v6/testdata/*.yaml v9/testdata/*.yaml
var processTestdata embed.FS

func TestProcessLockfile_ConsistentAcrossMultipleRuns(t *testing.T) {
	var testCases []string
	e := fs.WalkDir(processTestdata, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		testCases = append(testCases, path)
		return nil
	})
	assert.NoError(t, e)
	sort.Strings(testCases)

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			data, e := processTestdata.ReadFile(testCase)
			assert.NoError(t, e)

			var baseline []byte
			for i := 0; i < 10; i++ {
				var result processDirResult
				processLockfile(context.Background(), data, &result)
				assert.NoError(t, result.e)

				got, e := json.Marshal(result.trees)
				assert.NoError(t, e)
				if i == 0 {
					baseline = got
					continue
				}
				assert.Equal(t, string(baseline), string(got), "process result changed at run %d", i+1)
			}
		})
	}
}

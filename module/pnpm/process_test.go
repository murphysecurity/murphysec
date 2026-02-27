package pnpm

import (
	"context"
	"embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed v5/testdata/*.yaml v6/testdata/*.yaml v9/testdata/*.yaml
var processTestdata embed.FS

func TestProcessLockfile_ConsistentAcrossMultipleRuns(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{name: "v5-1", path: "v5/testdata/1.yaml"},
		{name: "v5-5", path: "v5/testdata/5.yaml"},
		{name: "v6-1", path: "v6/testdata/1.yaml"},
		{name: "v6-3", path: "v6/testdata/3.yaml"},
		{name: "v9-1", path: "v9/testdata/1.yaml"},
		{name: "v9-9", path: "v9/testdata/9.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, e := processTestdata.ReadFile(tc.path)
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

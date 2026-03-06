package npm

import (
	"encoding/json"
	"testing"

	"github.com/murphysecurity/murphysec/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessLockfileV3_PeerDepsAndAliasNames(t *testing.T) {
	lock := map[string]any{
		"name":            "renderer",
		"version":         "0.0.1",
		"lockfileVersion": 3,
		"packages": map[string]any{
			"": map[string]any{
				"dependencies": map[string]any{
					"vite":          "^5.2.14",
					"@isaacs/cliui": "^8.0.2",
				},
			},
			"node_modules/vite": map[string]any{
				"name":    "vite",
				"version": "5.2.14",
				"peerDependencies": map[string]any{
					"terser":      "^5.4.0",
					"@types/node": "^18.0.0 || >=20.0.0",
				},
			},
			"node_modules/terser": map[string]any{
				"name":    "terser",
				"version": "5.31.0",
				"dependencies": map[string]any{
					"@jridgewell/source-map": "^0.3.3",
					"commander":              "^2.20.0",
					"source-map-support":     "~0.5.20",
				},
			},
			"node_modules/terser/node_modules/commander": map[string]any{
				"name":    "commander",
				"version": "2.20.3",
			},
			"node_modules/@jridgewell/source-map": map[string]any{
				"name":    "@jridgewell/source-map",
				"version": "0.3.6",
			},
			"node_modules/source-map-support": map[string]any{
				"name":    "source-map-support",
				"version": "0.5.21",
				"dependencies": map[string]any{
					"buffer-from": "^1.0.0",
					"source-map":  "^0.6.0",
				},
			},
			"node_modules/buffer-from": map[string]any{
				"name":    "buffer-from",
				"version": "1.1.2",
			},
			"node_modules/source-map": map[string]any{
				"name":    "source-map",
				"version": "0.6.1",
			},
			"node_modules/@types/node": map[string]any{
				"name":    "@types/node",
				"version": "20.12.7",
				"dependencies": map[string]any{
					"undici-types": "~5.26.4",
				},
			},
			"node_modules/undici-types": map[string]any{
				"name":    "undici-types",
				"version": "5.26.5",
			},
			"node_modules/@isaacs/cliui": map[string]any{
				"name":    "@isaacs/cliui",
				"version": "8.0.2",
				"dependencies": map[string]any{
					"string-width-cjs": "npm:string-width@^4.2.0",
					"strip-ansi-cjs":   "npm:strip-ansi@^6.0.1",
					"wrap-ansi-cjs":    "npm:wrap-ansi@^7.0.0",
				},
			},
			"node_modules/string-width-cjs": map[string]any{
				"name":    "string-width",
				"version": "4.2.3",
			},
			"node_modules/strip-ansi-cjs": map[string]any{
				"name":    "strip-ansi",
				"version": "6.0.1",
			},
			"node_modules/wrap-ansi-cjs": map[string]any{
				"name":    "wrap-ansi",
				"version": "7.0.0",
			},
		},
	}
	data, err := json.Marshal(lock)
	require.NoError(t, err)
	parsed, err := processLockfileV3(data)
	require.NoError(t, err)

	got := map[string]string{}
	var walk func(items []model.DependencyItem)
	walk = func(items []model.DependencyItem) {
		for _, it := range items {
			got[it.CompName] = it.CompVersion
			walk(it.Dependencies)
		}
	}
	walk(parsed.Dependencies)

	assert.Equal(t, "0.3.6", got["@jridgewell/source-map"])
	assert.Equal(t, "20.12.7", got["@types/node"])
	assert.Equal(t, "1.1.2", got["buffer-from"])
	assert.Equal(t, "2.20.3", got["commander"])
	assert.Equal(t, "0.6.1", got["source-map"])
	assert.Equal(t, "0.5.21", got["source-map-support"])
	assert.Equal(t, "4.2.3", got["string-width-cjs"])
	assert.Equal(t, "6.0.1", got["strip-ansi-cjs"])
	assert.Equal(t, "5.31.0", got["terser"])
	assert.Equal(t, "5.26.5", got["undici-types"])
	assert.Equal(t, "7.0.0", got["wrap-ansi-cjs"])
}

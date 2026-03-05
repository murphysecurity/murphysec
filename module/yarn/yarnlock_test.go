package yarn

import (
	"context"
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/fixtures/*
var yarnLockFixtures embed.FS

func TestAnalyzeYarnDep_WithPreparedFixtures(t *testing.T) {
	testCases := []struct {
		name          string
		fixturePath   string
		dependency    string
		rangeVersion  string
		expectVersion string
		expectChild   []string
	}{
		{
			name:          "v1 lockfile",
			fixturePath:   "testdata/fixtures/101",
			dependency:    "accepts",
			rangeVersion:  "~1.3.3",
			expectVersion: "1.3.3",
			expectChild:   []string{"mime-types", "negotiator"},
		},
		{
			name:          "metadata version 4",
			fixturePath:   "testdata/fixtures/401",
			dependency:    "chalk",
			rangeVersion:  "^2.4.1",
			expectVersion: "2.4.2",
		},
		{
			name:          "metadata version 7",
			fixturePath:   "testdata/fixtures/701",
			dependency:    "@actions/core",
			rangeVersion:  "^1.2.6",
			expectVersion: "1.2.6",
		},
		{
			name:          "metadata version 8",
			fixturePath:   "testdata/fixtures/801",
			dependency:    "@antfu/install-pkg",
			rangeVersion:  "^0.1.1",
			expectVersion: "0.1.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lockData, e := yarnLockFixtures.ReadFile(tc.fixturePath)
			assert.NoError(t, e)

			dir := t.TempDir()
			e = os.WriteFile(filepath.Join(dir, "yarn.lock"), lockData, 0o600)
			assert.NoError(t, e)

			pkg := map[string]any{
				"name":    "ut-yarn",
				"version": "1.0.0",
				"dependencies": map[string]string{
					tc.dependency: tc.rangeVersion,
				},
			}
			pkgData, e := json.Marshal(pkg)
			assert.NoError(t, e)
			e = os.WriteFile(filepath.Join(dir, "package.json"), pkgData, 0o600)
			assert.NoError(t, e)

			deps, e := analyzeYarnDep(context.Background(), dir)
			assert.NoError(t, e)
			assert.NotEmpty(t, deps)

			root := findDepByName(deps, tc.dependency)
			if assert.NotNil(t, root) {
				assert.Equal(t, tc.expectVersion, root.Version)
				for _, child := range tc.expectChild {
					assert.Contains(t, depNames(root.Children), child)
				}
			}
		})
	}
}

func TestParseYarnLockYaml_SplitsCombinedSelectors(t *testing.T) {
	data, e := yarnLockFixtures.ReadFile("testdata/fixtures/801")
	assert.NoError(t, e)

	lockfile, e := parseYarnLockYaml(string(data))
	assert.NoError(t, e)
	assert.NotEmpty(t, lockfile)
	assert.Contains(t, lockfile, "@ampproject/remapping@npm:^2.2.0")
	assert.Contains(t, lockfile, "@ampproject/remapping@npm:^2.2.1")
}

func TestAnalyzeYarnDep_RoundTripStable(t *testing.T) {
	testCases := []struct {
		name         string
		fixturePath  string
		dependencies map[string]string
	}{
		{
			name:        "v1",
			fixturePath: "testdata/fixtures/101",
			dependencies: map[string]string{
				"accepts": "~1.3.3",
			},
		},
		{
			name:        "metadata-v4",
			fixturePath: "testdata/fixtures/401",
			dependencies: map[string]string{
				"chalk": "^2.4.1",
			},
		},
		{
			name:        "metadata-v7",
			fixturePath: "testdata/fixtures/701",
			dependencies: map[string]string{
				"@actions/core": "^1.2.6",
			},
		},
		{
			name:        "metadata-v8",
			fixturePath: "testdata/fixtures/801",
			dependencies: map[string]string{
				"@antfu/install-pkg": "^0.1.1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lockData, e := yarnLockFixtures.ReadFile(tc.fixturePath)
			assert.NoError(t, e)

			dir := t.TempDir()
			e = os.WriteFile(filepath.Join(dir, "yarn.lock"), lockData, 0o600)
			assert.NoError(t, e)

			pkg := map[string]any{
				"name":         "ut-yarn-roundtrip",
				"version":      "1.0.0",
				"dependencies": tc.dependencies,
			}
			pkgData, e := json.Marshal(pkg)
			assert.NoError(t, e)
			e = os.WriteFile(filepath.Join(dir, "package.json"), pkgData, 0o600)
			assert.NoError(t, e)

			var baseline string
			for i := 0; i < 10; i++ {
				deps, e := analyzeYarnDep(context.Background(), dir)
				assert.NoError(t, e)
				assert.NotEmpty(t, deps)

				normalized := normalizeDeps(deps)
				raw, e := json.Marshal(normalized)
				assert.NoError(t, e)
				if i == 0 {
					baseline = string(raw)
					continue
				}
				assert.Equal(t, baseline, string(raw), "roundtrip changed at run %d", i+1)
			}
		})
	}
}

func TestAnalyzeYarnDep_DedupAcrossDependenciesAndDevDependencies(t *testing.T) {
	lockData, e := yarnLockFixtures.ReadFile("testdata/fixtures/401")
	assert.NoError(t, e)

	dir := t.TempDir()
	e = os.WriteFile(filepath.Join(dir, "yarn.lock"), lockData, 0o600)
	assert.NoError(t, e)

	pkg := map[string]any{
		"name":    "ut-yarn-dedup",
		"version": "1.0.0",
		"dependencies": map[string]string{
			"chalk": "^2.4.1",
		},
		"devDependencies": map[string]string{
			"chalk": "^2.4.1",
		},
	}
	pkgData, e := json.Marshal(pkg)
	assert.NoError(t, e)
	e = os.WriteFile(filepath.Join(dir, "package.json"), pkgData, 0o600)
	assert.NoError(t, e)

	deps, e := analyzeYarnDep(context.Background(), dir)
	assert.NoError(t, e)

	count := 0
	for _, dep := range deps {
		if dep.Name == "chalk" && dep.Version == "2.4.2" {
			count++
		}
	}
	assert.Equal(t, 1, count, "chalk should be deduplicated across dependencies and devDependencies")
}

func findDepByName(deps []Dep, name string) *Dep {
	for i := range deps {
		if deps[i].Name == name {
			return &deps[i]
		}
	}
	return nil
}

func depNames(deps []Dep) []string {
	var r []string
	for _, dep := range deps {
		r = append(r, dep.Name)
	}
	return r
}

func normalizeDeps(deps []Dep) []Dep {
	out := make([]Dep, len(deps))
	for i := range deps {
		out[i] = Dep{
			Name:    deps[i].Name,
			Version: deps[i].Version,
		}
		out[i].Children = normalizeDeps(deps[i].Children)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].Version < out[j].Version
		}
		return out[i].Name < out[j].Name
	})
	return out
}

package pnpm

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed v6_test_file.yaml
var v6TestData []byte

func TestDependencyPath_getVersionFromPath(t *testing.T) {
	var s string
	var e error
	s, e = getVersionFromPath("/acorn-private-class-elements@1.0.0(acorn@8.8.2)")
	assert.NoError(t, e)
	assert.Equal(t, "1.0.0", s)
	s, e = getVersionFromPath("/@esbuild/android-arm64@0.17.19")
	assert.NoError(t, e)
	assert.Equal(t, "0.17.19", s)
}

func TestDependencyPath_getNameFromPath(t *testing.T) {
	var s string
	var e error
	s, e = getNameFromPath("/acorn-private-class-elements@1.0.0(acorn@8.8.2)")
	assert.NoError(t, e)
	assert.Equal(t, "acorn-private-class-elements", s)
	s, e = getNameFromPath("/@esbuild/android-arm64@0.17.19")
	assert.NoError(t, e)
	assert.Equal(t, "android-arm64", s)
}

func TestVersionIndicator(t *testing.T) {
	s, e := parseLockfileVersion(v6TestData)
	assert.NoError(t, e)
	assert.True(t, strings.HasPrefix(s, "6."))
}

func TestV6Postprocess(t *testing.T) {
	s, e := parseV6Lockfile(v6TestData, true)
	assert.NoError(t, e)
	assert.NotEmptyf(t, s.Packages, "s.Packages empty")
	assert.NotEmptyf(t, s.Dependencies, "s.Dependencies empty")
	assert.NotEmptyf(t, s.DevDependencies, "s.DevDependencies empty")
	assert.NotEmptyf(t, s.pathMapping, "s.pathMapping empty")
}

func TestV6BuildDependencyTree(t *testing.T) {
	s, e := parseV6Lockfile(v6TestData, true)
	assert.NoError(t, e)
	list, e := s.buildDependencyTree(true)
	assert.NoError(t, e)
	t.Log(list)
}

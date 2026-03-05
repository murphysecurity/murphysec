package poetry

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

//go:embed pyproject.toml
var data []byte

func TestParseToml(t *testing.T) {
	root := &tomlTree{}
	assert.NoError(t, toml.Unmarshal(data, &root.v))
	assert.Equal(t, "map[python:*]", fmt.Sprint(root.Get("tool", "poetry", "dependencies").v))
	assert.Equal(t, "poetry-demo", root.Get("tool", "poetry", "name").v)
}

//go:embed poetry.lock.py
var __lockData []byte

func TestParsePoetryLock(t *testing.T) {
	root := &tomlTree{}
	assert.NoError(t, toml.Unmarshal(__lockData, &root.v))
	assert.Equal(t, 13, len(root.Get("package").AsArray()))
	assert.Equal(t, "main", root.Get("package").AsArray()[0].Get("category").v)
}

//go:embed uv.lock
var __uvLockData []byte

func TestParseUvLock(t *testing.T) {
	dir := t.TempDir()
	uvLockPath := filepath.Join(dir, "uv.lock")
	assert.NoError(t, os.WriteFile(uvLockPath, __uvLockData, 0o600))

	deps, e := parsePoetryLock(context.Background(), uvLockPath)
	assert.NoError(t, e)
	assert.Equal(t, 3, len(deps))
	assert.Equal(t, "annotated-types", deps[0].CompName)
	assert.Equal(t, "0.7.0", deps[0].CompVersion)
	assert.Equal(t, "anyio", deps[1].CompName)
	assert.Equal(t, "4.10.0", deps[1].CompVersion)
	assert.Equal(t, "idna", deps[2].CompName)
	assert.Equal(t, "3.10", deps[2].CompVersion)
}

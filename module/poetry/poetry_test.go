package poetry

import (
	_ "embed"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
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

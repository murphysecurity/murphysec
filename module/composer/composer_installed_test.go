package composer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComposerInstalled_ObjectSchema(t *testing.T) {
	data := []byte(`{
  "packages": [
    {
      "name": "symfony/polyfill-mbstring",
      "version": "v1.29.0",
      "require": {
        "php": ">=7.1"
      }
    },
    {
      "name": "guzzlehttp/guzzle",
      "version": "7.8.1",
      "require": {
        "php": "^7.2.5 || ^8.0",
        "psr/http-client": "^1.0"
      }
    }
  ],
  "dev": false
}`)

	pkgs, err := parseComposerInstalled(data)
	require.NoError(t, err)
	require.Len(t, pkgs, 2)
	assert.Equal(t, "symfony/polyfill-mbstring", pkgs[0].Name)
	assert.Equal(t, "v1.29.0", pkgs[0].Version)
	assert.Contains(t, pkgs[0].Require, "php")
	assert.Equal(t, "guzzlehttp/guzzle", pkgs[1].Name)
}

func TestParseComposerInstalled_ArraySchema(t *testing.T) {
	data := []byte(`[
  {
    "name": "monolog/monolog",
    "version": "3.5.0",
    "require": {
      "php": ">=8.1",
      "psr/log": "^2.0 || ^3.0"
    }
  }
]`)

	pkgs, err := parseComposerInstalled(data)
	require.NoError(t, err)
	require.Len(t, pkgs, 1)
	assert.Equal(t, "monolog/monolog", pkgs[0].Name)
	assert.Equal(t, "3.5.0", pkgs[0].Version)
	assert.Contains(t, pkgs[0].Require, "php")
	assert.Contains(t, pkgs[0].Require, "psr/log")
}

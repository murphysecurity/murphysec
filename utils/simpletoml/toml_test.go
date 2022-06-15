package simpletoml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	// language=toml
	var data = `
[build-system]
requires = ["setuptools>=40.8.0", "wheel"]
build-backend = "setuptools.build_meta:__legacy__"
`
	to, e := UnmarshalTOML([]byte(data))
	assert.NoError(t, e)
	t.Log(to)
}

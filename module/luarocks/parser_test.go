package luarocks

import (
	"context"
	_ "embed"
	"github.com/murphysecurity/murphysec/module/luarocks/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testcase/1.lua
var tc1 string

//go:embed testcase/2.lua
var tc2 string

func Test_parse(t *testing.T) {
	var r parser.IStart_Context
	var e error
	r, e = parse(tc1)
	assert.NoError(t, e, "tc1")
	t.Log(r)
	r, e = parse(tc2)
	assert.NoError(t, e, "tc2")
	t.Log(r)
	r, e = parse("{")
	assert.Error(t, e, "aaa")
	t.Log(r)
}

func Test_analyze(t *testing.T) {
	var p, e = parse(tc1)
	assert.NoError(t, e)
	t.Log(analyze(context.TODO(), p))
}

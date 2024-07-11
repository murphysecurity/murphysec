package luarocks

import (
	_ "embed"
	"testing"
)

//go:embed testcase/1.lua
var tc1 string

//go:embed testcase/2.lua
var tc2 string

func Test_parse(t *testing.T) {
	parse(tc1)
	parse(tc2)
}

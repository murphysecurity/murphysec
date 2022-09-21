package sbt

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test_output
var __test_output []byte

func TestSbtDependencyTreeOutputParser_Result(t *testing.T) {
	writer := newSbtDependencyTreeOutputParser()
	_, e := writer.Write(__test_output)
	assert.NoError(t, e)
	writer.Close()
	root, e := writer.Result()
	assert.NoError(t, e)
	t.Log(root)
}

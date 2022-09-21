package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIntStack(t *testing.T) {
	var stack = NewIntStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	assert.Equal(t, 3, stack.Peek())
	assert.Equal(t, 3, stack.Pop())
	assert.Equal(t, 2, stack.Pop())
	assert.Equal(t, 1, stack.Pop())
	assert.Equal(t, 0, stack.Len())
	assert.Panics(t, func() {
		var a *IntStack
		a.Len()
	})
}

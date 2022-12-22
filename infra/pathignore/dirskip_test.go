package pathignore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirName(t *testing.T) {
	assert.True(t, DirName("node_modules"))
	assert.False(t, DirName("aaaa"))
}

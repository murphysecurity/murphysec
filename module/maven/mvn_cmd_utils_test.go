package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckMvnCommand(t *testing.T) {
	info, e := CheckMvnCommand()
	assert.NoError(t, e)
	t.Log(info.String())
}

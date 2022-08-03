package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckMvnCommand(t *testing.T) {
	info, e := CheckMvnCommand()
	assert.NoError(t, e)
	if e == nil {
		t.Log(info.String())
	}
}

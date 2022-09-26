package maven

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckMvnCommand(t *testing.T) {
	info, e := CheckMvnCommand(context.TODO())
	assert.NoError(t, e)
	if e == nil {
		t.Log(info.String())
	}
}

package semerr

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemErr(t *testing.T) {
	var ErrT = New("TE")
	e := ErrT.Decorate(errors.New("awsl"))
	assert.True(t, errors.Is(e, ErrT))
	assert.True(t, errors.Is(ErrT, e))
	assert.True(t, errors.Is(e, e))
	assert.False(t, errors.Is(errors.New("awsl"), errors.New("awsl")))
}

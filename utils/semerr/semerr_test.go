package semerr

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSemErr(t *testing.T) {
	var ErrT = New("TE")
	e := ErrT.Decorate(errors.Wrap(errors.New("awsl"), "awsl*2"))
	assert.True(t, errors.Is(e, ErrT))
	assert.True(t, errors.Is(ErrT, e))
	assert.True(t, errors.Is(e, e))
	assert.False(t, errors.Is(errors.New("awsl"), errors.New("awsl")))
	var ErrE = New("EE")
	e = ErrE.Decorate(e)
	assert.True(t, errors.Is(e, ErrE))
	assert.True(t, errors.Is(e, ErrT))
	var ErrQ = New("QE")
	assert.False(t, errors.Is(e, ErrQ))
	assert.False(t, errors.Is(e, ErrQ.Decorate(io.EOF)))
}

func TestSemErrPrint(t *testing.T) {
	var ErrT = New("TE")
	e1 := errors.New("foo")
	e2 := errors.Wrap(e1, "bar")
	e := ErrT.Decorate(e2)
	t.Logf("%+v\n", e)
	t.Logf("============\n")
	t.Logf("%v\n", e)
}

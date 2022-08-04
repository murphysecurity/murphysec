package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	var err = New("test")
	assert.Error(t, err)
	assert.Equal(t, "test", err.Error())
	var we = Wrap(err, "wrapped")
	assert.Equal(t, "wrapped: test", we.Error())
	assert.True(t, Is(we, err))
	assert.Equal(t, err, Unwrap(we))

	var de = WithDetail(err, "detail")
	assert.Equal(t, "test: detail", de.Error())
	assert.ErrorIs(t, de, err)

	var cause = New("cause")
	var wc = WithCause(err, cause)
	assert.Equal(t, "test: cause", wc.Error())
	assert.ErrorIs(t, wc, cause)
	assert.NotErrorIs(t, wc, New("test"))
}

package logpipe

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLogPipe(t *testing.T) {
	l, _ := zap.NewDevelopment()
	p := New(l, "aa")
	var e error
	_, e = p.Write([]byte(`aaaaa1
ssssss
ssssss
ssssss`))
	assert.NoError(t, e)
	_, e = p.Write([]byte(`aaaaa2
`))
	assert.NoError(t, e)
	assert.NoError(t, p.Close())
}

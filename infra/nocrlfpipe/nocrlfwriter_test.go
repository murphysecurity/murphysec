package nocrlfpipe

import (
	"bytes"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoCrLfWriter_Write(t *testing.T) {
	var a = "foo\r\nbar"
	var buf = &bytes.Buffer{}
	must.A(NewNoCrlfWriter(buf).Write([]byte(a)))
	assert.Equal(t, "foobar", buf.String())
	assert.NotEqual(t, "foobar", a)
}

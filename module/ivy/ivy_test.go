package ivy

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed ivy-1.xml.test
var sample1 []byte

func TestParseIvy(t *testing.T) {
	data, e := readIvyXml(nil, bytes.NewReader(sample1))
	assert.NoError(t, e)
	assert.NotEmpty(t, data.Dependencies)
	for _, d := range data.Dependencies {
		t.Logf("%s:%s\n", d.CompName, d.CompVersion)
	}
}

package maven

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_ResolveProperty(t *testing.T) {
	m := map[string]string{
		"a": "1",
		"b": "2${a}",
		"c": "foo${b}${d}",
	}
	assert.Equal(t, _resolveProperty(m, nil, "c"), "foo21${d}")
	t.Log(_resolveProperty(m, nil, "c"))
}

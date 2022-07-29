package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProperty_Resolve(t *testing.T) {
	m := map[string]string{
		"a": "1",
		"b": "2${a}",
		"c": "foo${b}${d}${}aaa${}",
	}
	resolver := newProperties()
	resolver.PutMap(m)
	assert.Equal(t, "foo21${d}${}aaa${}", resolver.Resolve(m["c"]))
}

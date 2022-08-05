package conan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConanError_Error(t *testing.T) {
	var input = `poco/1.9.4aaaa: Not found in local cache, looking in remotes...
poco/1.9.4aaaa: Trying with 'conancenter'...
ERROR: Unable to find 'poco/1.9.4aaaa' in remotes
`
	var ce = error(conanError(input))
	assert.ErrorIs(t, ce, ErrConanReport)
	assert.Equal(t, ce.Error(), "ERROR: Unable to find 'poco/1.9.4aaaa' in remotes")
}

package rebar3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var input = `===> Verifying dependencies...
└─ myproj─0.1.0 (project app)
   ├─ eleveldb─2.0.2+build.565.ref061405f (git repo)
   ├─ lager─3.9.2 (hex package)
   │  └─ goldrush─0.1.9 (hex package)
   ├─ riak_dt─2.1.1 (hex package)
   ├─ sext─1.5.0 (hex package)
   └─ swc─1.0.0 (git repo)
`
	var nodes = parseRebar3TreeOutput(input)
	assert.Equal(t, 1, len(nodes))
}

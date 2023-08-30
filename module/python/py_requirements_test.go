package python

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestParseRequirements(t *testing.T) {
	var data = `pbr!=2.1.0,>=2.0.0 # Apache-2.0
docker>=2.4.2 # Apache-2.0
Jinja2!=2.9.0,!=2.9.1,!=2.9.2,!=2.9.3,!=2.9.4,>=2.8 # BSD License (3 clause)
gitdb>=0.6.4 # BSD License (3 clause)
GitPython>=1.0.1 # BSD License (3 clause)
six>=1.10.0 # MIT
oslo.config>=5.1.0 # Apache-2.0
oslo.utils>=3.33.0 # Apache-2.0
setuptools!=24.0.0,!=34.0.0,!=34.0.1,!=34.0.2,!=34.0.3,!=34.1.0,!=34.1.1,!=34.2.0,!=34.3.0,!=34.3.1,!=34.3.2,!=36.2.0,>=16.0.0 # PSF/ZPL
netaddr>=0.7.18 # BSD

`
	var expect = `GitPython1.0.1 Jinja22.8 docker2.4.2 gitdb0.6.4 netaddr0.7.18 oslo.config5.1.0 oslo.utils3.33.0 pbr2.0.0 setuptools16.0.0 six1.10.0`
	var r []string
	for k, v := range parseRequirements(data) {
		r = append(r, k+v)
	}
	sort.Strings(r)
	var actual = strings.Join(r, " ")
	t.Log(actual)
	assert.Equal(t, expect, actual)
}

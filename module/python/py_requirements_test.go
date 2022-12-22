package python

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestParseRequirements(t *testing.T) {
	var data = `# development tools
Werkzeug>=0.6.2

# Testing dependencies
mock>=1.0.1
WebTest>=1.3.4
django-webtest>=1.5.3
factory-boy==2.1.1
httpretty==0.6.3
ruamel.yaml=1.0.0

# documentation
Sphinx==1.2b3

# Code style and coverage
flake8>=0.8
coveralls>=0.1.1,<0.2

# we need this to make sure that we test against Oscar 0.6
-e git+https://github.com/tangentlabs/django-oscar.git@89d12c8701d293f23afa19c6efac17b249ae1b6d#egg=django-oscar

# Others
Whoosh>=2.4.1
`
	var expect = `Sphinx1.2b3 WebTest1.3.4 Werkzeug0.6.2 Whoosh2.4.1 django-webtest1.5.3 factory-boy2.1.1 flake80.8 httpretty0.6.3 mock1.0.1 ruamel.yaml1.0.0`
	var r []string
	for k, v := range parseRequirements(data) {
		r = append(r, k+v)
	}
	sort.Strings(r)
	var actual = strings.Join(r, " ")
	t.Log(actual)
	assert.Equal(t, expect, actual)
}

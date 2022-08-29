package python

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/model"
	"github.com/stretchr/testify/assert"
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
	//language=json
	var rs = `[
  {
  "name": "Werkzeug",
  "version": "0.6.2"
  },
  {
  "name": "mock",
  "version": "1.0.1"
  },
  {
  "name": "WebTest",
  "version": "1.3.4"
  },
  {
  "name": "django-webtest",
  "version": "1.5.3"
  },
  {
  "name": "factory-boy",
  "version": "2.1.1"
  },
  {
  "name": "httpretty",
  "version": "0.6.3"
  },
  {
  "name": "Sphinx",
  "version": "1.2b3"
  },
  {
  "name": "flake8",
  "version": "0.8"
  },
  {
  "name": "Whoosh",
  "version": "2.4.1"
  }
  ]
`
	var r []model.Dependency
	assert.NoError(t, json.Unmarshal([]byte(rs), &r))
	assert.Equal(t, r, parseRequirements(data))
}

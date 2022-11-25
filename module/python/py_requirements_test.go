package python

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
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

	var rs =
	// language=json
	`[
          {
            "comp_name": "Werkzeug",
            "comp_version": "0.6.2",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "mock",
            "comp_version": "1.0.1",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "WebTest",
            "comp_version": "1.3.4",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "django-webtest",
            "comp_version": "1.5.3",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "factory-boy",
            "comp_version": "2.1.1",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "httpretty",
            "comp_version": "0.6.3",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "ruamel.yaml",
            "comp_version": "1.0.0",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "Sphinx",
            "comp_version": "1.2b3",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "flake8",
            "comp_version": "0.8",
            "ecosystem": "pip",
            "repository": ""
          },
          {
            "comp_name": "Whoosh",
            "comp_version": "2.4.1",
            "ecosystem": "pip",
            "repository": ""
          }
        ]`
	t.Log(string(must.A(json.MarshalIndent(parseRequirements(data), "", "  "))))
	var r []model.DependencyItem
	assert.NoError(t, json.Unmarshal([]byte(rs), &r))
	assert.Equal(t, r, parseRequirements(data))
}

package python

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tomlBuildSys(t *testing.T) {
	//language=toml
	var _t = `
[build-system]
requires = [
    "wheel",
    "setuptools",
    "Cython>=0.29.13",
    "numpy==1.13.3; python_version=='3.5' and platform_system!='AIX'",
    "numpy==1.17.3; python_version>='3.8' and platform_system!='AIX'",
    "pybind11>=2.2.4",
]
`
	//language=json
	var target = `[{"name":"wheel","version":""},{"name":"setuptools","version":""},{"name":"Cython","version":"0.29.13"},{"name":"numpy","version":"1.13.3"},{"name":"pybind11","version":"2.2.4"}]`
	r, e := tomlBuildSys([]byte(_t))
	assert.NoError(t, e)
	var a, b any
	must.Must(json.Unmarshal([]byte(target), &a))
	must.Must(json.Unmarshal(must.A(json.Marshal(r)), &b))
	assert.EqualValues(t, a, b)
}

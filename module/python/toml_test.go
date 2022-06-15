package python

import (
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
    "numpy==1.13.3; python_version=='3.6' and platform_system!='AIX'",
    "numpy==1.14.5; python_version=='3.7' and platform_system!='AIX'",
    "numpy==1.17.3; python_version>='3.8' and platform_system!='AIX'",
    "numpy==1.16.0; python_version=='3.5' and platform_system=='AIX'",
    "numpy==1.16.0; python_version=='3.6' and platform_system=='AIX'",
    "numpy==1.16.0; python_version=='3.7' and platform_system=='AIX'",
    "numpy==1.17.3; python_version>='3.8' and platform_system=='AIX'",
    "pybind11>=2.2.4",
]
`
	r, e := tomlBuildSys([]byte(_t))
	assert.NoError(t, e)
	t.Log(r)
}

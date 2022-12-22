package python

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parsePipListJson(t *testing.T) {
	var r = `[{"name": "certifi", "version": "2022.9.24"}, {"name": "charset-normalizer", "version": "2.0.12"}, {"name": "cos-python-sdk-v5", "version": "1.9.20"}, {"name": "crcmod", "version": "1.7"}, {"name": "dicttoxml", "version": "1.7.4"}, {"name": "idna", "version": "3.4"}, {"name": "kafka-python", "version": "2.0.2"}, {"name": "pip", "version": "9.0.3"}, {"name": "pycryptodome", "version": "3.15.0"}, {"name": "requests", "version": "2.27.1"}, {"name": "setuptools", "version": "39.2.0"}, {"name": "six", "version": "1.16.0"}, {"name": "urllib3", "version": "1.26.12"}]`
	m, e := parsePipListJson([]byte(r))
	assert.NoError(t, e)
	assert.NotNil(t, m)
}

func Test_parsePipListDefaultFormat(t *testing.T) {
	var r = `
DEPRECATION: The default format will switch to columns in the future. You can use --format=(legacy|columns) (or define a format=(legacy|columns) in your pip.conf under the [list] section) to disable this warning.
certifi (2022.9.24)
urllib3 (1.26.12)
`
	m, e := parsePipListDefaultFormat([]byte(r))
	assert.NoError(t, e)
	assert.NotNil(t, m)
	assert.Len(t, m, 2)
}

package gradle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseGradleVersionOutput(t *testing.T) {
	var data = `------------------------------------------------------------
Gradle 6.7.1
------------------------------------------------------------

Build time:   2020-11-16 17:09:24 UTC
Revision:     2972ff02f3210d2ceed2f1ea880f026acfbab5c0

Kotlin:       1.3.72
Groovy:       2.5.12
Ant:          Apache Ant(TM) version 1.10.8 compiled on May 10 2020
JVM:          1.8.0_345 (Azul Systems, Inc. 25.345-b01)
OS:           Windows 10 10.0 amd64
`
	v, e := parseGradleVersionOutput(data)
	assert.NoError(t, e)
	t.Log(v)

}

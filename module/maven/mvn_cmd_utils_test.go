package maven

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckMvnCommand(t *testing.T) {
	info, e := CheckMvnCommand(context.TODO())
	assert.NoError(t, e)
	if e == nil {
		t.Log(info.String())
	}
}

func Test_ParseMvnVersion(t *testing.T) {
	s := `Apache Maven 3.8.6 (84538c9988a25aec085021c365c560670ad80f63)
Maven home: C:\Users\iseki\scoop\apps\maven\current
Java version: 1.8.0_345, vendor: Azul Systems, Inc., runtime: C:\Program Files\Zulu\zulu-8\jre
Default locale: zh_CN, platform encoding: GBK
OS name: "windows 10", version: "10.0", arch: "amd64", family: "windows"

`
	v := parseMvnVersion(s)
	assert.Equal(t, "3.8.6", v)
}

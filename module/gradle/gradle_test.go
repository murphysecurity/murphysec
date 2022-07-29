package gradle

import (
	_ "embed"
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"strings"
	"testing"
)

//go:embed parse_gradle_dep_testcase
var __parse_gradle_dep_testcase0 string

func TestGradleDep(t *testing.T) {
	lines := strings.Split(__parse_gradle_dep_testcase0, "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	must.A(json.Marshal(parseGradleDependencies(lines)))
}

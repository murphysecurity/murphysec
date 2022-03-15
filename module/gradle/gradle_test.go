package gradle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/utils/must"
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
	fmt.Println(string(must.Byte(json.Marshal(parseGradleDependencies(lines)))))
}

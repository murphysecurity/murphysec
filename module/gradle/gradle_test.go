package gradle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/utils/must"
	"strings"
	"testing"
)

//go:embed parse_gradle_dep_testcase
var __parse_gradle_dep_testcase0 string

//go:embed parse_gradle_dep_testcase2
var __parse_gradle_dep_testcase2 string

func TestGradleDep(t *testing.T) {
	lines := strings.Split(__parse_gradle_dep_testcase0, "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	fmt.Println(string(must.A(json.Marshal(parseGradleDependencies(lines)))))
}

func TestGradleDep2(t *testing.T) {
	lines := strings.Split(__parse_gradle_dep_testcase2, "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	fmt.Println(string(must.A(json.Marshal(parseGradleDependencies(lines)))))
}

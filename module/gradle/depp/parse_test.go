package depp

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//go:embed parse_gradle_dep_testcase
var parse_gradle_dep_testcase []byte

//go:embed parse_gradle_dep_testcase2
var parse_gradle_dep_testcase2 []byte

//go:embed parse_gradle_dep_testcase3
var parse_gradle_dep_testcase3 []byte

//go:embed parse_gradle_dep_testcase4
var parse_gradle_dep_testcase4 []byte

func Test_parse_1(t *testing.T) {
	e := Parse(bytes.NewReader(parse_gradle_dep_testcase), func(project string, task string, data []TreeNode) {
		fmt.Println(project, " -> ", task)
		printTree(data, 0, os.Stdout)
	})
	assert.NoError(t, e)
}

func Test_parse_2(t *testing.T) {
	e := Parse(bytes.NewReader(parse_gradle_dep_testcase2), func(project string, task string, data []TreeNode) {
		fmt.Println(project, " -> ", task)
		printTree(data, 0, os.Stdout)
	})
	assert.NoError(t, e)
}
func Test_parse_3(t *testing.T) {
	e := Parse(bytes.NewReader(parse_gradle_dep_testcase3), func(project string, task string, data []TreeNode) {
		fmt.Println(project, " -> ", task)
		printTree(data, 0, os.Stdout)
	})
	assert.NoError(t, e)
}
func Test_parse_4(t *testing.T) {
	e := Parse(bytes.NewReader(parse_gradle_dep_testcase4), func(project string, task string, data []TreeNode) {
		fmt.Println(project, " -> ", task)
		printTree(data, 0, os.Stdout)
	})
	assert.NoError(t, e)
}

package depp

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
func Test_file(t *testing.T) {

	filename := "C:\\Users\\陈浩轩\\Desktop\\mofei\\murphysec\\module\\gradle\\depp\\parse_gradle_dep_testcase4"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	newContent := make([]byte, 0, len(content))
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 180 {
			newContent = append(newContent, line[180:]...)
		} else {
			newContent = append(newContent, line...)
		}
		newContent = append(newContent, '\n')
	}
	scanner.Err()
	err = ioutil.WriteFile(filename, newContent, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Content has been updated.")
}

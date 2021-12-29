package scanner

import (
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/util/must"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(string(must.Byte(json.MarshalIndent(ScanDir("C:\\Users\\iseki\\GolandProjects"), "", "  "))))
}
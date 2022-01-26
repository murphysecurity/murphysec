package findfile

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	l := []string{
		"a",
		"a/b/c",
		"b/c",
	}
	tempBase := must.String(os.MkdirTemp("", "test-testfind-*"))
	defer must.Must(os.RemoveAll(tempBase))
	for _, s := range l {
		fp := filepath.Join(tempBase, s)
		must.Must(os.MkdirAll(fp, 0755))
	}

	// test
	iter := Find(tempBase, Option{
		MaxDepth:    2,
		ExcludeFile: false,
		ExcludeDir:  false,
		Predication: func(name string, dir string) bool {
			return true
		},
	})
	var rs []string
	for iter.Next() {
		must.Must(iter.Err())
		rs = append(rs, strings.ReplaceAll(must.String(filepath.Rel(tempBase, iter.Path())), "\\", "/"))
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i] < rs[j]
	})
	//fmt.Println(string(must.Byte(json.Marshal(rs))))
	fmt.Println(string(must.Byte(json.Marshal(rs))))
	assert.Equal(t, string(must.Byte(json.Marshal(rs))), "[\"a\",\"a/b\",\"a/b/c\",\"b\",\"b/c\"]")
}

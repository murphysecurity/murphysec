package buildout

import (
	"bytes"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"unsafe"
)
import _ "embed"

//go:embed test_metadata_1
var testMetadata1 string

//go:embed test_metadata_2
var testMetadata2 string

//go:embed test_metadata_3
var testMetadata3 string

func TestParseMetadata(t *testing.T) {
	testParseMetadata(t, "test_metadata_1", testMetadata1)
	testParseMetadata(t, "test_metadata_2", testMetadata2)
	testParseMetadata(t, "test_metadata_3", testMetadata3)
}

func testParseMetadata(t *testing.T, name string, data string) {
	t.Run(name, func(t *testing.T) {
		var inputBytes = unsafe.Slice(unsafe.StringData(data), len(data))
		var result, err = ParseMetadata(bytes.NewReader(inputBytes))
		assert.NoError(t, err)
		var entries = lo.Entries(result)
		sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
		for _, entry := range entries {
			for _, value := range entry.Value {
				t.Logf("%s: %s", entry.Key, value)
			}
		}
		assert.NotEmpty(t, result["Name"])
		assert.NotEmpty(t, result["Version"])
	})

}

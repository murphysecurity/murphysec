package cocoapods

import (
	_ "embed"
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed Podfile.lock
var testData string

func TestPodParser(t *testing.T) {
	root, e := parse(testData)
	assert.NoError(t, e)
	t.Log(string(must.A(json.MarshalIndent(root, "", "  "))))
	//t.Log(string(must.Byte(json.MarshalIndent(root.get("DEPENDENCIES:"), "", "  "))))
}

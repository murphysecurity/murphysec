package cocoapods

import (
	_ "embed"
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test_Podfile_lock
var testData string

//go:embed test_podlock.json
var testResult string

func TestPodParser(t *testing.T) {
	root, e := parse(testData)
	assert.NoError(t, e)
	var a map[string]any
	var b map[string]any
	must.Must(json.Unmarshal([]byte(testResult), &a))
	must.Must(json.Unmarshal(must.A(json.Marshal(root)), &b))
	assert.EqualValues(t, a, b)
}

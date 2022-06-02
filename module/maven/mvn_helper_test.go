package maven

import (
	_ "embed"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"murphysec-cli-simple/utils/must"
	"testing"
)

//go:embed test-graph.json
var testGraph []byte

func TestParseGraph(t *testing.T) {
	var d dependencyGraph
	assert.NoError(t, json.Unmarshal(testGraph, &d))

	t.Log(string(must.Byte(json.MarshalIndent(d.Tree(), "", "  "))))
}

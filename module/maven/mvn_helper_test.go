package maven

import (
	_ "embed"
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test-graph.json
var testGraph []byte

func TestParseGraph(t *testing.T) {
	var d dependencyGraph
	assert.NoError(t, json.Unmarshal(testGraph, &d))

	t.Log(string(must.A(json.MarshalIndent(d.Tree(), "", "  "))))
}

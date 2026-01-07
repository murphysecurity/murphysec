package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzippedBody_Read(t *testing.T) {
	var data = map[string]any{}
	for i := 0; i < 10000; i++ {
		data[fmt.Sprint(i)] = fmt.Sprint(i)
	}
	var body = NewGzippedBody(NewJsonRequestBody(data))
	gr, _ := gzip.NewReader(body)
	var dec = json.NewDecoder(gr)
	var data2 map[string]any
	assert.NoError(t, dec.Decode(&data2))
	assert.EqualValues(t, data, data2)
}

package pubspec

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata-1
var testdata1 []byte

func TestName(t *testing.T) {
	var d, e = parseFile(context.TODO(), bytes.NewReader(testdata1))
	assert.NoError(t, e)
	assert.Equal(t, 23, len(d))
}

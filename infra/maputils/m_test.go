package maputils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var __mrefMap = map[string]int64{
	"a": 3,
	"b": 2,
	"c": 1,
}

func TestKeysSortedByValue(t *testing.T) {
	assert.EqualValues(t, []string{"c", "b", "a"}, KeysSortedByValue(__mrefMap))
}

func TestValuesSortedByKey(t *testing.T) {
	assert.EqualValues(t, []int64{3, 2, 1}, ValuesSortedByKey(__mrefMap))
}

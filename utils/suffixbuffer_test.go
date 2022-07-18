package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuffixBuffer1(t *testing.T) {
	var d = [][2]string{
		{"123", "123"},
		{"1234", "234"},
		{"12345678", "678"},
		{"1", "1"},
	}
	for _, it := range d {
		r := MkSuffixBuffer(3)
		r.write([]byte(it[0]))
		assert.Equal(t, it[1], string(r.Bytes()))
	}
}
func TestSuffixBuffer2(t *testing.T) {
	r := MkSuffixBuffer(3)
	r.write([]byte("12"))
	r.write([]byte("34"))
	r.write([]byte("567"))
	assert.Equal(t, "567", string(r.Bytes()))
}

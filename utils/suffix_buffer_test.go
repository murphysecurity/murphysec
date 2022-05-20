package utils

import (
	"github.com/stretchr/testify/assert"
	"murphysec-cli-simple/utils/must"
	"testing"
)

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func TestSuffixBuffer1(t *testing.T) {
	b := NewSuffixBuffer(3)
	must.Int(b.Write([]byte("1")))
	must.Int(b.Write([]byte("23")))
	assert.False(t, b.Truncated())
	assert.Equal(t, "123", b.String())
}

func TestSuffixBuffer2(t *testing.T) {
	b := NewSuffixBuffer(3)
	must.Int(b.Write([]byte("1")))
	must.Int(b.Write([]byte("234")))
	must.Int(b.Write([]byte("567")))
	must.Int(b.Write([]byte("89")))
	assert.True(t, b.Truncated())
	assert.Equal(t, "789", b.String())
}

func TestSuffixBuffer3(t *testing.T) {
	b := NewSuffixBuffer(3)
	must.Int(b.Write([]byte("1")))
	must.Int(b.Write([]byte("23")))
	must.Int(b.Write([]byte("4")))
	assert.True(t, b.Truncated())
	assert.Equal(t, "234", b.String())
}

func TestSuffixBuffer4(t *testing.T) {
	b := NewSuffixBuffer(3)
	must.Int(b.Write([]byte("12345678")))
	assert.True(t, b.Truncated())
	assert.Equal(t, "678", b.String())
}

func TestSuffixBuffer5(t *testing.T) {
	b := NewSuffixBuffer(3)
	must.Int(b.Write([]byte("1")))
	must.Int(b.Write([]byte("2")))
	assert.False(t, b.Truncated())
	assert.Equal(t, "12", b.String())
}

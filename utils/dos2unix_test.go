package utils

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"murphysec-cli-simple/utils/must"
	"testing"
)

func TestDos2UnixWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := Dos2UnixWriter(b)
	must.Int(w.Write([]byte("aa\r\nbb\ncc\r\r\na")))
	w.Close()
	assert.Equal(t, b.Bytes(), []byte("aa\nbb\ncc\r\na"))
}

func TestUnix2DosWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := Unix2DosWriter(b)
	must.Int(w.Write([]byte("aa\r\nbb\n\r\ncc\n")))
	w.Close()
	assert.Equal(t, b.Bytes(), []byte("aa\r\nbb\r\n\r\ncc\r\n"))
}

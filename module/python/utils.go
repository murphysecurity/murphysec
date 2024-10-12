package python

import (
	"bufio"
	"golang.org/x/net/html/charset"
	"io"
	"os"
)

func readTextFile(path string, maxLength int) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer func() { _ = f.Close() }()
	r, e := charset.NewReader(bufio.NewReader(f), "")
	if e != nil {
		_ = f.Close()
		f, e = os.Open(path)
		if e != nil {
			return nil, e
		}
		r = f
	}
	return io.ReadAll(io.LimitReader(r, int64(maxLength)))
}

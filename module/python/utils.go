package python

import (
	"io"
	"os"
)

func readFile(path string, maxLength int) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer func() { _ = f.Close() }()
	return io.ReadAll(io.LimitReader(f, int64(maxLength)))
}

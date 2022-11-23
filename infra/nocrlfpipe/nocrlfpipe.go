package nocrlfpipe

import (
	"io"
)

type noCrLfWriter struct {
	w io.Writer
}

func (w noCrLfWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func NewNoCrlfWriter(w io.Writer) io.WriteCloser {
	return noCrLfWriter{w: w}
}

func (w noCrLfWriter) Write(input []byte) (int, error) {
	buf := make([]byte, 0, len(input))
	for _, it := range input {
		if it == '\n' || it == '\r' {
			continue
		}
		buf = append(buf, it)
	}
	if len(buf) == 0 {
		return len(input), nil
	}
	if _, e := w.w.Write(buf); e != nil {
		return 0, e
	} else {
		return len(input), nil
	}
}

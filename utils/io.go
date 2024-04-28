package utils

import (
	"errors"
	"io"
)

func MultiWriteCloser(writers ...io.Writer) io.WriteCloser {
	return mwc{writers: append([]io.Writer{}, writers...)}
}

type mwc struct {
	writers []io.Writer
}

func (m mwc) Write(p []byte) (n int, err error) {
	var es []error
	for _, it := range m.writers {
		var e error
		n, e = it.Write(p)
		if e != nil {
			es = append(es, e)
		}
	}
	if len(es) > 0 {
		return n, errors.Join(es...)
	}
	return n, nil
}

func (m mwc) Close() error {
	var es []error
	for _, it := range m.writers {
		if c, ok := it.(io.Closer); ok {
			if e := c.Close(); e != nil {
				es = append(es, e)
			}
		}
	}
	if len(es) > 0 {
		return errors.Join(es...)
	}
	return nil
}

var _ io.WriteCloser = (*mwc)(nil)

package utils

import "io"

func Dos2UnixWriter(w io.Writer) io.WriteCloser {
	return &dos2unixW{w: w}
}

type dos2unixW struct {
	w      io.Writer
	lastCR bool
}

func (d *dos2unixW) Close() error {
	if d.lastCR {
		if _, e := d.w.Write([]byte{'\r'}); e != nil {
			return e
		}
	}
	if c, ok := d.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (d *dos2unixW) Write(p []byte) (n int, err error) {
	b := make([]byte, 0, len(p))
	for _, ch := range p {
		if d.lastCR {
			switch ch {
			case '\n':
				b = append(b, '\n')
			case '\r':
				b = append(b, '\r')
			default:
				b = append(b, '\r', ch)
			}
		} else {
			if ch != '\r' {
				b = append(b, ch)
			}
		}
		d.lastCR = ch == '\r'
	}
	_, e := d.w.Write(b)
	if e != nil {
		return 0, e
	} else {
		return len(p), e
	}

}

func Unix2DosWriter(w io.Writer) io.WriteCloser {
	return &unix2dosW{w: w}
}

type unix2dosW struct {
	w      io.Writer
	lastCR bool
}

func (u *unix2dosW) Write(p []byte) (n int, err error) {
	b := make([]byte, 0, len(p))
	for _, ch := range p {
		if u.lastCR {
			u.lastCR = false
			b = append(b, ch)
		} else {
			if ch == '\n' {
				b = append(b, '\r')
			}
			if ch == '\r' {
				u.lastCR = true
			}
			b = append(b, ch)
		}
	}
	_, e := u.w.Write(b)
	if e != nil {
		return 0, e
	} else {
		return len(p), e
	}
}

func (u *unix2dosW) Close() error {
	if c, ok := u.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

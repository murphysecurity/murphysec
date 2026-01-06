package api

import (
	"bytes"
	"compress/gzip"
	"io"
)

func NewGzippedBody(body interface {
	GetBody() (io.ReadCloser, error)
	io.ReadCloser
}) interface {
	GetBody() (io.ReadCloser, error)
	io.ReadCloser
} {
	return &gzippedBody{body: body, reader: body}
}

type gzippedBody struct {
	buf        bytes.Buffer
	end        bool
	gzipWriter *gzip.Writer
	body       interface {
		GetBody() (io.ReadCloser, error)
		io.ReadCloser
	}
	reader io.ReadCloser
}

func (g *gzippedBody) GetBody() (io.ReadCloser, error) {
	reader, e := g.body.GetBody()
	if e != nil {
		return nil, e
	}
	return &gzippedBody{body: g.body, reader: reader}, nil
}

func (g *gzippedBody) Read(p []byte) (n int, err error) {
	if g.gzipWriter == nil {
		g.gzipWriter = gzip.NewWriter(&g.buf)
	}
	for g.buf.Available() == 0 {
		if g.end {
			return 0, io.EOF
		}
		_, e := io.CopyN(g.gzipWriter, g.reader, 16*1024)
		if e == io.EOF {
			g.end = true
			e = g.gzipWriter.Close()
			if e != nil {
				return 0, e
			}
		}
		if e != nil {
			return 0, e
		}
	}

	return g.reader.Read(p)
}

func (g *gzippedBody) Close() error {
	return g.reader.Close()
}

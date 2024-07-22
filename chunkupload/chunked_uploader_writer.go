package chunkupload

import (
	"bytes"
	"context"
	"golang.org/x/sync/errgroup"
	"io"
)

type uploadWriter struct {
	MinChunkSize int
	UploadFn     func(chunkId int, data []byte) error
	buf          *bytes.Buffer
	chunkId      int
	eg           *errgroup.Group
	ec           context.Context
	closed       bool
}

func NewUploadWriter(ctx context.Context, minChunkSize int, concurrency int, uploadFn func(chunkId int, data []byte) error) io.WriteCloser {
	if ctx == nil {
		panic("ctx == nil")
	}
	if minChunkSize < 0 {
		panic("minChunkSize < 0")
	}
	if concurrency < 1 {
		panic("concurrency < 1")
	}
	if uploadFn == nil {
		panic("uploadFn == nil")
	}
	eg, ec := errgroup.WithContext(ctx)
	eg.SetLimit(concurrency)
	return &uploadWriter{
		MinChunkSize: minChunkSize,
		UploadFn:     uploadFn,
		eg:           eg,
		ec:           ec,
	}
}

func (u *uploadWriter) ensureInit() {
	if u.buf == nil {
		u.buf = new(bytes.Buffer)
	}
}

func (u *uploadWriter) Write(p []byte) (n int, err error) {
	if u.closed {
		return 0, io.ErrClosedPipe
	}
	u.ensureInit()
	if u.ec.Err() != nil {
		return 0, u.ec.Err()
	}
	var dataToSend []byte
	if u.buf.Len() == 0 && len(p) > u.MinChunkSize {
		// fast-path, avoid memory copy
		dataToSend = p
	} else {
		n, err = u.buf.Write(p)
		if u.buf.Len() < u.MinChunkSize {
			return
		}
		dataToSend = u.buf.Bytes()
		u.buf = nil
	}
	if dataToSend != nil {
		u.chunkId++
		u.eg.Go(func() error {
			return u.UploadFn(u.chunkId, dataToSend)
		})
	}
	return
}

func (u *uploadWriter) Close() error {
	u.ensureInit()
	if u.buf.Len() > 0 {
		var data = u.buf.Bytes()
		u.buf = nil
		u.chunkId++
		u.eg.Go(func() error {
			return u.UploadFn(u.chunkId, data)
		})
	}
	u.closed = true
	return u.eg.Wait()
}

var _ io.WriteCloser = (*uploadWriter)(nil)

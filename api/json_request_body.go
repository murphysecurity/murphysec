package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

type JsonRequestBody interface {
	GetBody() (io.ReadCloser, error)
	io.ReadCloser
}

func NewJsonRequestBody(V any) JsonRequestBody {
	return &jsonReqBody{v: V}
}

type jsonReqBody struct {
	v      any
	pr     *io.PipeReader
	closed bool
}

func (j *jsonReqBody) GetBody() (io.ReadCloser, error) {
	return &jsonReqBody{v: j.v}, nil
}

func (j *jsonReqBody) Read(p []byte) (n int, err error) {
	if j.closed {
		return 0, io.ErrClosedPipe
	}
	if j.pr == nil {
		var pw *io.PipeWriter
		j.pr, pw = io.Pipe()
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					switch e := r.(type) {
					case error:
						_ = pw.CloseWithError(e)
					default:
						_ = pw.CloseWithError(errors.New("panic in json encoding routing"))
					}
				}
			}()
			var bw = bufio.NewWriter(pw)
			defer func() { _ = pw.Close() }()
			var e = json.NewEncoder(pw).Encode(j.v)
			if e != nil {
				_ = pw.CloseWithError(e)
				return
			}
			e = bw.Flush()
			if e != nil {
				_ = pw.CloseWithError(e)
				return
			}
			_ = pw.Close()
		}()
	}
	return j.pr.Read(p)
}

func (j *jsonReqBody) Close() error {
	if j.closed {
		return io.ErrClosedPipe
	}
	j.closed = true
	return j.pr.Close()
}

var _ io.ReadCloser = &jsonReqBody{}

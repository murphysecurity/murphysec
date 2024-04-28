package maven

import (
	"bufio"
	"bytes"
	list "github.com/bahlo/generic-list-go"
	"io"
	"sync"
	"unsafe"
)

type eRecorder struct {
	s  string
	wg sync.WaitGroup
	pw io.WriteCloser
}

func (r *eRecorder) Write(p []byte) (n int, err error) {
	return r.pw.Write(p)
}

func (r *eRecorder) Close() error {
	return r.pw.Close()
}

func (r *eRecorder) String() string {
	r.wg.Wait()
	return r.s
}

const prefixErrorS = "[ERROR"

var prefixError = unsafe.Slice(unsafe.StringData(prefixErrorS), len(prefixErrorS))

func launchERecorder() (r *eRecorder) {
	pr, pw := io.Pipe()
	r = &eRecorder{pw: pw}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		var scanner = bufio.NewScanner(pr)
		scanner.Split(bufio.ScanLines)
		scanner.Buffer(nil, 4096)
		var q = list.New[[]byte]()
		for scanner.Scan() {
			var b = scanner.Bytes()
			if !bytes.HasPrefix(b, prefixError) {
				continue
			}
			q.PushBack(b)
			if q.Len() > 10 {
				q.Remove(q.Front())
			}
		}
		var byteList [][]byte
		for q.Len() > 0 {
			var t = q.Front()
			byteList = append(byteList, t.Value)
			q.Remove(t)
		}
		var bs = bytes.Join(byteList, []byte("\n"))
		if len(bs) > 0 {
			r.s = unsafe.String(&bs[0], len(bs))
		}
	}()
	return
}

var _ io.WriteCloser = (*eRecorder)(nil)

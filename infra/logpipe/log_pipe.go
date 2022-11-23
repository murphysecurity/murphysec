package logpipe

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
)

type Pipe struct {
	w *io.PipeWriter
}

func (l *Pipe) Write(data []byte) (int, error) {
	w := l.w
	_, _ = w.Write(data)
	return len(data), nil
}

func (l *Pipe) Close() error {
	return l.w.Close()
}

func New(logger *zap.Logger, prefix string) *Pipe {
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.WithOptions(zap.WithCaller(false))
	r, w := io.Pipe()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	lp := &Pipe{w: w}
	go func() {
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			logger.Debug(fmt.Sprintf("%s: %s", prefix, scanner.Text()))
		}
		// drain
		for {
			if _, e := r.Read(make([]byte, 128)); e != nil {
				break
			}
		}
	}()
	return lp
}

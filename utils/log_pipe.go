package utils

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
)

type LogPipe struct {
	w *io.PipeWriter
}

func (l *LogPipe) Write(data []byte) (int, error) {
	w := l.w
	_, _ = w.Write(data)
	return len(data), nil
}

func (l *LogPipe) Close() error {
	return l.w.Close()
}

func NewLogPipe(logger *zap.Logger, prefix string) *LogPipe {
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.WithOptions(zap.WithCaller(false))
	r, w := io.Pipe()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	lp := &LogPipe{w: w}
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

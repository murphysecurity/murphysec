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

func NewLogPipe(logger *zap.Logger) *LogPipe {
	r, w := io.Pipe()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	lp := &LogPipe{w: w}
	go func() {
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			logger.Info(fmt.Sprintf("Maven output: %s", scanner.Text()))
		}
		// drain
		_, _ = r.Read(make([]byte, 128))
	}()
	return lp
}

package utils

import (
	"go.uber.org/zap"
	"io"
)

func LogCloseErr(logger *zap.SugaredLogger, pipeName string, closer io.Closer) {
	if closer == nil {
		panic("closer == nil")
	}
	if e := closer.Close(); e != nil {
		logger.Warnf("close pipe %s failed: %v", pipeName, e)
	}
}

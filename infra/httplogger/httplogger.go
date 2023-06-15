package httplogger

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"strings"
)

type middlewareWithZapLogger struct {
	Tripper http.RoundTripper
	Logger  *zap.Logger
}

func (m *middlewareWithZapLogger) RoundTrip(request *http.Request) (resp *http.Response, e error) {
	var dump []byte
	dumpBody := true
	if strings.Contains(request.Header.Get("content-type"), "octet-stream") {
		dumpBody = false
	}
	dump, e = httputil.DumpRequestOut(request, dumpBody)
	if e != nil {
		m.Logger.Error("dump request out failed", zap.Error(e))
	} else {
		m.Logger.Debug("Http request", zap.ByteString("dump", dump))
	}
	resp, e = m.Tripper.RoundTrip(request)
	if e != nil {
		return
	}
	dump, e = httputil.DumpResponse(resp, true)
	if e != nil {
		m.Logger.Error("dump response out failed", zap.Error(e))
	} else {
		m.Logger.Debug("http response", zap.ByteString("dump", dump))
	}
	return
}

var _ http.RoundTripper = (*middlewareWithZapLogger)(nil)

func New(tripper http.RoundTripper, logger *zap.Logger) http.RoundTripper {
	if tripper == nil {
		panic("tripper == nil")
	}
	if logger == nil {
		panic("logger == nil")
	}
	return &middlewareWithZapLogger{
		Tripper: tripper,
		Logger:  logger,
	}
}

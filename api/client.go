package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/tlserr"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
)

//go:generate stringer -type apiError -output apierror_string.go -linecomment
type apiError int

const (
	_                        apiError = iota
	ErrTLSError                       // api: tls error
	ErrTimeout                        // api: timeout
	ErrGeneral                        // api: general error
	ErrUnprocessableResponse          // api: cannot process server response
	ErrTokenInvalid                   // api: token invalid
	ErrBadURL                         // api: bad URL
)

func (i apiError) Error() string {
	return i.String()
}

var _DefaultClient *Client

func DefaultClient() *Client {
	return _DefaultClient
}

func InitDefaultClient(config *Config) error {
	must.NotNil(config)
	c, e := config.Build()
	if e != nil {
		return e
	}
	_DefaultClient = c
	return nil
}

type Client struct {
	client  *http.Client
	baseUrl *url.URL
	token   string
	logger  *zap.Logger
}

func (c *Client) DoJson(req *http.Request, resBody interface{}) (e error) {
	if t := reflect.TypeOf(resBody); t != nil && t.Kind() != reflect.Ptr {
		panic("resBody must be a pointer or nil")
	}
	var (
		httpResponse *http.Response
		logger       = c.logger.Sugar()
	)

	defer func() {
		if e != nil {
			logger.Errorf("Request error: %v", e)
		}
	}()

	req.Header.Set("User-Agent", version.UserAgent())
	req.Header.Set("Authorization", "Bearer "+c.token)
	logger.Debugf("Request: %v", req.URL)
	httpResponse, e = c.client.Do(req)
	if isHttpTimeout(e) {
		return ErrTimeout
	}
	if tlserr.IsTLSError(e) {
		return errors.WithCause(ErrTLSError, e)
	}
	if e != nil {
		return errors.WithCause(ErrGeneral, e)
	}
	logger.Infof("API response - %d", httpResponse.StatusCode)
	var data []byte
	data, e = io.ReadAll(httpResponse.Body)
	if e != nil {
		return errors.WithCause(ErrGeneral, e)
	}
	defer httpResponse.Body.Close()

	var statusCode = httpResponse.StatusCode

	// Normal code
	if statusCode >= 200 && statusCode < 300 {
		if resBody == nil {
			return nil
		}
		if e = json.Unmarshal(data, resBody); e != nil {
			return errors.WithCause(ErrUnprocessableResponse, e)
		}
		return nil
	}

	// Error
	if statusCode == 401 {
		return ErrTokenInvalid
	}
	var m GeneralError
	if e = json.Unmarshal(data, &m); e != nil {
		logger.Error("Server error response can't be parsed, suppressed", zap.Error(e))
		return fmt.Errorf("%w ([%d]%s)", ErrUnprocessableResponse, statusCode, httpResponse.Status)
	} else {
		return &m
	}
}

func (c *Client) GET(url *url.URL) *http.Request {
	return must.A(http.NewRequest(http.MethodGet, url.String(), nil))
}

func (c *Client) POST(url *url.URL, body io.Reader) *http.Request {
	return must.A(http.NewRequest(http.MethodPost, url.String(), body))
}

func (c *Client) PostJson(url *url.URL, data any) *http.Request {
	u := c.POST(url, bytes.NewReader(must.A(json.Marshal(data))))
	u.Header.Set("Content-Type", "application/json")
	return u
}

func isHttpTimeout(e error) bool {
	r, ok := e.(*url.Error)
	return ok && r.Timeout()
}

type GeneralError struct {
	Code  int    `json:"code"`
	MsgZh string `json:"msg_zh"`
}

func (c *GeneralError) Error() string {
	return fmt.Sprintf("[%d]%s", c.Code, c.MsgZh)
}

func joinURL(base *url.URL, relPath string) *url.URL {
	var u = *base // copy
	u.Path = path.Join(u.Path, relPath)
	return &u
}

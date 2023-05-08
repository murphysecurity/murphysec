package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
)

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
	ctx     context.Context
	client  *http.Client
	baseUrl *url.URL
	token   string
	logger  *zap.Logger
}

func (c Client) BaseURLText() string {
	return strings.TrimSuffix(c.baseUrl.String(), "/")
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
	if e != nil {
		return errorOf(e)
	}
	logger.Infof("API response - %d", httpResponse.StatusCode)
	var data []byte
	defer httpResponse.Body.Close()
	data, e = io.ReadAll(httpResponse.Body)
	if e != nil {
		return errorOf(e)
	}

	var statusCode = httpResponse.StatusCode
	// Normal code
	if statusCode >= 200 && statusCode < 300 {
		if resBody == nil {
			return nil
		}
		if e = json.Unmarshal(data, resBody); e != nil {
			return &Error{
				HTTPStatus:            statusCode,
				UnprocessableResponse: true,
				Cause:                 e,
			}
		}
		return nil
	}

	// Error
	var m generalErrorResponse
	if e = json.Unmarshal(data, &m); e != nil {
		return &Error{
			Cause:                 e,
			HTTPStatus:            statusCode,
			UnprocessableResponse: true,
		}
	} else {
		return &Error{
			HTTPStatus: statusCode,
			Code:       m.Code,
			Message:    m.MsgZh,
		}
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

type generalErrorResponse struct {
	Code  int    `json:"code"`
	MsgZh string `json:"msg_zh"`
}

func joinURL(base *url.URL, relPath string) *url.URL {
	var u = *base // copy
	u.Path = path.Join(u.Path, relPath)
	return &u
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Logger = zap.NewNop()
var NetworkLogger = zap.NewNop()

const HeaderMachineId = "machine-id"

var machineId = version.MachineId()

var C *Client

type _LoggingMiddleware struct {
	Transport http.RoundTripper
}

func (t *_LoggingMiddleware) RoundTrip(request *http.Request) (resp *http.Response, e error) {
	var dump []byte
	dump, e = httputil.DumpRequestOut(request, true)
	if e != nil {
		Logger.Error("dump request out failed", zap.Error(e))
	} else {
		NetworkLogger.Debug("Http request", zap.ByteString("dump", dump))
	}
	resp, e = t.Transport.RoundTrip(request)
	if e != nil {
		return
	}
	dump, e = httputil.DumpResponse(resp, true)
	if e != nil {
		Logger.Error("dump response out failed", zap.Error(e))
	} else {
		NetworkLogger.Debug("http response", zap.ByteString("dump", dump))
	}
	return
}

type Client struct {
	client  *http.Client
	baseUrl string
	Token   string
}

func (c *Client) BaseURL() string {
	return strings.TrimRight(c.baseUrl, "/")
}

func NewClient() *Client {
	cl := &Client{
		client: &http.Client{
			Transport: &_LoggingMiddleware{Transport: http.DefaultTransport},
			Timeout:   time.Second * 300,
		},
		baseUrl: env.ServerBaseUrl(),
	}
	if i, e := strconv.Atoi(os.Getenv("API_TIMEOUT")); e != nil {
		cl.client.Timeout = time.Duration(int64(time.Second) * int64(i))
	}
	Logger.Info("Http client created", zap.String("baseUrl", env.ServerBaseUrl()), zap.Duration("timeout", cl.client.Timeout))
	return cl
}

func (c *Client) POST(relUri string, body io.Reader) *http.Request {
	u, e := http.NewRequest(http.MethodPost, c.baseUrl+relUri, body)
	u.Header.Set(HeaderMachineId, machineId)
	if e != nil {
		panic(e)
	}
	return u
}

func (c *Client) PostJson(relUri string, a interface{}) *http.Request {
	u := c.POST(relUri, bytes.NewReader(must.A(json.Marshal(a))))
	u.Header.Set("Content-Type", "application/json")
	u.Header.Set(HeaderMachineId, machineId)
	return u
}

func (c *Client) GET(relUri string) *http.Request {
	u, e := http.NewRequest(http.MethodGet, c.baseUrl+relUri, nil)
	u.Header.Set(HeaderMachineId, machineId)
	if e != nil {
		panic(e)
	}
	return u
}

func (c *Client) DoJson(req *http.Request, resBody interface{}) (err error) {
	if t := reflect.TypeOf(resBody); t != nil && t.Kind() != reflect.Ptr {
		panic("resBody must be a pointer or nil")
	}

	req.Header.Set("User-Agent", version.UserAgent())
	Logger.Debug("Send request", zap.String("uri", req.URL.RequestURI()))

	var httpResponse *http.Response
	httpResponse, err = c.client.Do(req)
	if utils.IsHttpTimeout(err) {
		return ErrTimeout
	}
	if utils.IsTlsCertError(err) {
		return errors.WithCause(ErrTlsRequest, err)
	}
	if err != nil {
		return errors.WithCause(ErrServerRequest, err)
	}
	Logger.Info("API response", zap.Any("status", httpResponse.StatusCode))
	var data []byte
	data, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return errors.WithCause(ErrServerRequest, err)
	}
	defer httpResponse.Body.Close()

	var statusCode = httpResponse.StatusCode

	// Normal code
	if statusCode >= 200 && statusCode < 300 {
		if resBody == nil {
			return nil
		}
		if e := json.Unmarshal(data, resBody); e != nil {
			return errors.WithCause(UnprocessableResponse, e)
		}
		return nil
	}

	// Error
	httpMsg := fmt.Sprintf("HTTP status %d - %s", statusCode, httpResponse.Status)
	if statusCode == 401 {
		return ErrTokenInvalid
	}
	var m CommonApiErr
	if e := json.Unmarshal(data, &m); e != nil {
		Logger.Error("Server error response can't be parsed, suppressed", zap.Error(e))
	} else {
		return &m
	}
	return errors.WithDetail(ErrServerRequest, httpMsg)
}

type CommonApiErr struct {
	EError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

var BaseCommonApiError = &CommonApiErr{}

func (c *CommonApiErr) Error() string {
	return fmt.Sprintf("[%d]%s", c.EError.Code, c.EError.Details)
}

func (c *CommonApiErr) Is(e error) bool {
	return e == BaseCommonApiError
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func (c *Client) DoJson(req *http.Request, resBody interface{}) error {
	var noBody bool
	if t := reflect.TypeOf(resBody); t == nil {
		noBody = true
	} else {
		if t.Kind() != reflect.Ptr {
			panic("resBody must be a pointer")
		}
	}
	Logger.Info("Send request", zap.String("uri", req.URL.RequestURI()))
	res, e := c.client.Do(req)
	if e != nil {
		e := e.(*url.Error)
		Logger.Info("Request failed", zap.Error(e))
		if e.Timeout() {
			Logger.Error("Request timeout")
			return ErrTimeout
		}
		return errors.Wrap(ErrServerRequest, e.Error())
	}
	Logger.Info("API response", zap.Any("status", res.StatusCode))
	data, e := io.ReadAll(res.Body)
	if e != nil {
		return errors.Wrap(ErrServerRequest, "read response body failed:"+e.Error())
	}
	defer res.Body.Close()
	var mimeType string
	contentType := res.Header.Get("content-type")
	if contentType != "" {
		var err error
		mimeType, _, err = mime.ParseMediaType(contentType)
		if err != nil {
			return errors.Wrap(ErrServerRequest, "parse content-type failed: "+e.Error())
		}
	}
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if noBody {
			return nil
		}
		if mimeType != "application/json" {
			return errors.Wrap(UnprocessableResponse, "MIME-type: "+mimeType)
		}
		if e := json.Unmarshal(data, resBody); e != nil {
			Logger.Error("Parse server response json failed", zap.Error(e))
			return errors.Wrap(ErrServerRequest, "parse response body as json failed")
		}
		return nil
	}
	if res.StatusCode >= 400 {
		baseErr := ErrServerRequest
		httpMsg := fmt.Sprintf("http status %d - %s", res.StatusCode, res.Status)
		if res.StatusCode == 401 {
			return ErrTokenInvalid
		}
		if mimeType == "" {
			return errors.Wrap(baseErr, httpMsg)
		}
		if mimeType == "application/json" {
			var m CommonApiErr
			if e := json.Unmarshal(data, &m); e != nil {
				Logger.Error("Parse server response json failed", zap.Error(e))
				return errors.Wrap(baseErr, httpMsg)
			}
			return &m
		}
	}
	return errors.Wrap(ErrServerRequest, fmt.Sprintf("http code %d - %s", res.StatusCode, res.Status))
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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrTokenInvalid = model.WrapIdeaErr(errors.New("Token invalid"), model.IdeaApiTimeout)
var ErrServerRequest = model.WrapIdeaErr(errors.New("Send request failed"), model.IdeaServerRequestFailed)
var UnprocessableResponse = model.WrapIdeaErr(errors.New("Unprocessable response"), model.IdeaServerRequestFailed)
var ErrTimeout = model.WrapIdeaErr(errors.New("API request timeout"), model.IdeaApiTimeout)

var C *Client

type Client struct {
	client  *http.Client
	baseUrl string
	Token   string
}

func NewClient(baseUrl string) *Client {
	c := new(http.Client)
	p := regexp.MustCompile("/*$")
	baseUrl = p.ReplaceAllString(strings.TrimSpace(baseUrl), "")
	c.Timeout = time.Second * 300
	i, e := strconv.Atoi(os.Getenv("API_TIMEOUT"))
	if e == nil && i > 0 {
		c.Timeout = time.Duration(int64(time.Second) * int64(i))
	}
	cl := &Client{client: c, baseUrl: baseUrl}
	return cl
}

func (c *Client) POST(relUri string, body io.Reader) *http.Request {
	u, e := http.NewRequest(http.MethodPost, c.baseUrl+relUri, body)
	if e != nil {
		panic(e)
	}
	return u
}

func (c *Client) PostJson(relUri string, a interface{}) *http.Request {
	u := c.POST(relUri, bytes.NewReader(must.A(json.Marshal(a))))
	u.Header.Set("Content-Type", "application/json")
	return u
}

func (c *Client) GET(relUri string) *http.Request {
	u, e := http.NewRequest(http.MethodGet, c.baseUrl+relUri, nil)
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
	logger.Logger.Info("Send request: ", req.URL.RequestURI())
	res, e := c.client.Do(req)
	if e != nil {
		e := e.(*url.Error)
		logger.Logger.Error("Request failed: ", e.Error())
		if e.Timeout() {
			logger.Logger.Error("Request timeout")
			return ErrTimeout
		}
		return errors.Wrap(ErrServerRequest, e.Error())
	}
	logger.Logger.Info("API response:", res.StatusCode)
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
			logger.Debug.Println("Server data:", string(data))
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
				logger.Debug.Println("Server data:", string(data))
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

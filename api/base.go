package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/logger"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var CliServerAddressOverride string

// serverAddress returns API logger URL
func serverAddress() string {
	var envServer string
	envServer = strings.Trim(strings.TrimSpace(CliServerAddressOverride), "/")
	if envServer == "" {
		envServer = strings.Trim(strings.TrimSpace(os.Getenv("MPS_CLI_SERVER")), "/")
	}
	if len(envServer) == 0 {
		return "https://www.murphysec.com"
	}
	return envServer
}

var client *http.Client

func init() {
	c := new(http.Client)
	c.Timeout = time.Second * 300
	i, e := strconv.Atoi(os.Getenv("API_TIMEOUT"))
	if e == nil && i > 0 {
		c.Timeout = time.Duration(int64(time.Second) * int64(i))
	}
	client = c
}

var ErrTokenInvalid = errors.New("Token invalid")
var ErrSendRequest = errors.New("Send request failed")

type CommonApiErr struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

func readCommonErr(data []byte, statusCode int) error {
	var c CommonApiErr
	if e := json.Unmarshal(data, &c); e != nil {
		return errors.Wrap(e, "read error json failed")
	}
	return errors.New(fmt.Sprintf("API err[%d]: %s", statusCode, c.Error.Message))
}

func readHttpBody(res *http.Response) ([]byte, error) {
	data, e := io.ReadAll(res.Body)
	if e != nil {
		logger.Warn.Println("read body failed.", e.Error())
		return nil, e
	}
	logger.Debug.Println("body size", len(data), "bytes")
	_ = res.Body.Close()
	return data, e
}

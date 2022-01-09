package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/utils/simplejson"
	"murphysec-cli-simple/version"
	"net/http"
	"time"
)

// CheckAPIToken returns a boolean indicating the token availability.
func CheckAPIToken(token string) (bool, error) {
	if len(token) == 0 {
		return false, errors.New("Token not set")
	}
	url := serverAddress() + "/v1/token/check"
	client := http.Client{Timeout: 10 * time.Second}
	body := map[string]interface{}{
		"token": token,
	}
	data := must.Byte(json.Marshal(body))
	req, e := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	must.Must(e)
	req.Header.Set("user-agent", version.UserAgent())
	req.Header.Set("content-type", "application/json")
	post, err := client.Do(req)
	if err != nil {
		return false, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer post.Body.Close()
	if post.StatusCode != 200 {
		return false, errors.New(fmt.Sprintf("API request failed: %d", post.StatusCode))
	}
	j, err := simplejson.NewFromReader(post.Body)
	if err != nil {
		return false, err
	}
	if checkBool, b := j.Get("data", "available").CheckBool(); !b {
		return false, errors.New("API error: response can't be processed.")
	} else {
		return checkBool, nil
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/simplejson"
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
	post, err := client.Post(url, "application/json", bytes.NewBuffer(data))
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

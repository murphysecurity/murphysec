package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/simplejson"
	"net/http"
	"time"
)

func Report(body *ScanRequestBody) (*ScanResult, error) {
	if defaultToken == "" {
		return nil, errors.New("API token not set")
	}
	url := serverAddress() + "/v1/cli/report2"
	client := http.Client{Timeout: 300 * time.Second}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(must.Byte(json.Marshal(body))))
	must.Must(err)
	request.Header.Add("Authorization", fmt.Sprintf("token %s", defaultToken))
	output.Debug(fmt.Sprintf("Request: %s", request.RequestURI))
	do, err := client.Do(request)
	if err != nil {
		output.Error(fmt.Sprintf("err: %v", err.Error()))
		return nil, err
	}
	output.Debug(fmt.Sprintf("Response: [%d]%s", do.StatusCode, do.Status))
	//goland:noinspection GoUnhandledErrorResult
	defer do.Body.Close()
	if err != nil {
		return nil, err
	}
	j, err := simplejson.NewFromReader(do.Body)
	if do.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed, statusCode: %d", do.StatusCode)
	}
	if ec := j.Get("code").Int(); ec != 0 {
		return nil, fmt.Errorf("API request failed: %d - %s", ec, j.Get("info").String())
	}
	var r ScanResult
	if e := json.Unmarshal(must.Byte(json.Marshal(j.Get("data"))), &r); e != nil {
		return nil, errors.Wrap(e, "API result unmarshal failed")
	}
	return &r, nil
}

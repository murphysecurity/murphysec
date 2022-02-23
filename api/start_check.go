package api

import (
	"bytes"
	"encoding/json"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

func StartCheck(taskId string) error {
	must.True(taskId != "")
	body := must.Byte(json.Marshal(map[string]interface{}{"task_info": taskId}))
	resp, e := http.Post(serverAddress()+"/message/v2/access/client/start_check", "application/json", bytes.NewReader(body))
	if e != nil {
		return ErrSendRequest
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	data, e := readHttpBody(resp)
	if e != nil {
		return e
	}
	return readCommonErr(data, resp.StatusCode)
}

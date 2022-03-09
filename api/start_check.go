package api

import (
	"bytes"
	"encoding/json"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

type TaskKind string

const (
	TaskKindNormal TaskKind = "Normal"
	TaskKindBinary TaskKind = "Binary"
)

func StartCheck(taskId string, taskKind TaskKind) error {
	must.True(taskId != "")
	body := must.Byte(json.Marshal(map[string]interface{}{"task_info": taskId, "task_kind": taskKind}))
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

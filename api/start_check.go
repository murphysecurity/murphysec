package api

import (
	"murphysec-cli-simple/utils/must"
)

func StartCheck(taskId string) error {
	must.True(taskId != "")
	httpReq := C.PostJson("/message/v2/access/client/start_check", map[string]interface{}{"task_info": taskId})
	return C.DoJson(httpReq, nil)
}

func StartCheckTaskType(taskId string, kind TaskKind) error {
	must.True(taskId != "")
	httpReq := C.PostJson("/message/v2/access/client/start_check", map[string]interface{}{"task_info": taskId, "task_kind": kind})
	return C.DoJson(httpReq, nil)
}

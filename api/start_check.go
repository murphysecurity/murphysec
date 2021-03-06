package api

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
)

func StartCheckTaskType(taskId string, kind model.TaskKind) error {
	must.True(taskId != "")
	httpReq := C.PostJson("/message/v2/access/client/start_check", map[string]interface{}{"task_info": taskId, "task_kind": kind})
	return C.DoJson(httpReq, nil)
}

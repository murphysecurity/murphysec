package api

import (
	"fmt"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/utils/must"
	"time"
)

func QueryResult(taskId string) (*model.TaskScanResponse, error) {
	must.True(taskId != "")
	for {
		var r = struct {
			Data model.TaskScanResponse `json:"data"`
		}{}
		httpReq := C.GET(fmt.Sprintf("/message/v2/access/detect/task_scan?scan_id=%s", taskId))
		if e := C.DoJson(httpReq, &r); e != nil {
			return nil, e
		}
		if !r.Data.Complete {
			time.Sleep(time.Second * 2)
			continue
		}
		return &r.Data, nil
	}
}

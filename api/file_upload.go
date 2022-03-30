package api

import (
	"fmt"
	"io"
)

func UploadChunk(taskId string, chunkId int, reader io.Reader) error {
	req := C.POST("/message/v2/access/client/upload_check_files", reader)
	v := req.URL.Query()
	v.Set("task_info", taskId)
	v.Set("chunk_id", fmt.Sprintf("%04d", chunkId))
	req.URL.RawQuery = v.Encode()
	req.Header.Set("Content-Type", "application/gzip")
	return C.DoJson(req, nil)
}

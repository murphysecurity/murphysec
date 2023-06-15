package api

import (
	"io"
	"strconv"
)

func UploadCheckFiles(client *Client, taskId string, subTaskId string, chunkId int, reader io.Reader) error {
	checkNotNull(client)
	checkNotZeroInt(chunkId)
	checkNotNull(reader)
	u := joinURL(client.baseUrl, "/platform3/v3/client/upload_check_file")
	q := u.Query()
	q.Add("chunk_id", strconv.Itoa(chunkId))
	q.Add("subtask_id", subTaskId)
	u.RawQuery = q.Encode()
	req := client.POST(u, reader)
	req.Header.Set("content-type", "application/octet-stream")
	return client.DoJson(req, nil)
}

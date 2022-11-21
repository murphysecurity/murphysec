package api

import "io"

func UploadCheckFiles(client *Client, taskId string, subTaskId string, chunkId int, reader io.Reader) error {
	checkNotNull(client)
	checkNotZeroInt(chunkId)
	checkNotNull(reader)
	return client.DoJson(client.POST("/v3/client/upload_check_files", reader), nil)
}

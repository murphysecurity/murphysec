package api

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"time"
)

func QueryResult(ctx context.Context, client *Client, subtaskId string) (r *model.ScanResultResponse, e error) {
	for {
		r, e := QueryResultImmediately(client, subtaskId)
		if e != nil {
			return nil, e
		}
		if r.Complete {
			return r, nil
		}
		time.Sleep(time.Second * 2)
	}
}
func QueryResultImmediately(client *Client, subtaskId string) (r *model.ScanResultResponse, e error) {
	checkNotNull(client)
	e = client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/result"), map[string]any{"subtask_id": subtaskId}), &r)
	return
}

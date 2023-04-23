package api

func StartCheck(client *Client, subtaskId string) error {
	checkNotNull(client)
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/start_check"), map[string]any{"subtask_id": subtaskId}), nil)
}

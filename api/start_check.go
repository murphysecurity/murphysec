package api

import "github.com/murphysecurity/murphysec/model"

func StartCheck(client *Client, task *model.ScanTask) error {
	checkNotNull(client)
	// 后端一定要我把 maven 参数再传一遍
	var data = map[string]any{
		"subtask_id":           task.SubtaskId,
		"package_private_name": task.MavenSourceName,
	}
	if task.MavenSourceId != "" {
		data["package_private_id"] = task.MavenSourceId
	}
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/start_check"), data), nil)
}

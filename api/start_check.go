package api

type Subtype string

const (
	SubtypeSourcecode Subtype = "sourcecode"
	SubtypeSBOM       Subtype = "sbom"
)

func StartCheck(client *Client, subtaskId string, subtype Subtype) error {
	checkNotNull(client)
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/start_check"), map[string]any{"subtask_id": subtaskId, "subtype": subtype}), nil)
}

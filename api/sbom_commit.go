package api

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
)

func SubmitSBOM(client *Client, subtaskId string, modules []model.Module, codeFragments []model.ComponentCodeFragment) error {
	checkNotNull(client)
	must.NotZero(subtaskId)
	var req = map[string]any{
		"subtask_id":     subtaskId,
		"modules":        utils.NoNilSlice(modules),
		"code_fragments": utils.NoNilSlice(codeFragments),
	}
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/upload_data"), req), nil)
}

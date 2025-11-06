package api

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
)

func SubmitSBOM(ctx context.Context, client *Client, subtaskId string, modules []model.Module, codeFragments []model.ComponentCodeFragment) error {
	checkNotNull(client)
	must.NotZero(subtaskId)
	var req = map[string]any{
		"subtask_id":     subtaskId,
		"modules":        utils.NoNilSlice(modules),
		"code_fragments": utils.NoNilSlice(codeFragments),
		"scan_warnings":  scanerr.GetAll(ctx),
	}
	if env.StorageUploadSBom != "" {
		var f = must.A(os.OpenFile(env.StorageUploadSBom, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666))
		var bf = bufio.NewWriter(f)
		var encoder = json.NewEncoder(bf)
		must.M(encoder.Encode(req))
		must.M(bf.Flush())
		must.M(f.Close())
	}
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/upload_data"), req), nil)
}

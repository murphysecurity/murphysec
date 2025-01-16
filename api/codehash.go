package api

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"io"
	"net/http"
	"os"
)

func SubmitCodeHash(ctx context.Context, client *Client, path string) error {
	checkNotNull(client)
	var u = joinURL(client.baseUrl, "/platform3/v3/client/upload_feature").String()
	var fBody = func() (io.ReadCloser, error) {
		f, e := os.Open(path)
		if e != nil {
			return nil, e
		}
		return utils.NewBufferedReader(f), nil
	}
	firstTimeBody, e := fBody()
	if e != nil {
		return e
	}
	var req = must.A(http.NewRequestWithContext(ctx, http.MethodPost, u, firstTimeBody))
	req.GetBody = fBody
	req.Header.Set("Content-Type", "application/json")
	return client.DoJson(req, nil)
}

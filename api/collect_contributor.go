package api

import (
	"context"
	"github.com/murphysecurity/murphysec/collect"
)

func ReportCollectedContributors(ctx context.Context, client *Client, data *collect.ContributorUpload) {
	_ = client.DoJson(client.PostJson(joinURL(client.baseUrl, "/committer/save"), data), nil)
}

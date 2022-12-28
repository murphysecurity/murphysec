package maven

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

const httpRemoteMaxBody = 4 * 1024 * 1024

func newHttpRemote(baseURL url.URL) *httpRemote {
	return &httpRemote{baseURL: baseURL}
}

type httpRemote struct {
	baseURL url.URL
}

func (h *httpRemote) GetPath(ctx context.Context, _path string) ([]byte, error) {
	u := h.baseURL
	u.Path = path.Join(u.Path, _path)
	request, e := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if e != nil {
		return nil, e
	}
	resp, e := http.DefaultClient.Do(request)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if e != nil {
		return nil, e
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrRemoteNoResource
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d - %s", resp.StatusCode, resp.Status)
	}
	if resp.ContentLength > httpRemoteMaxBody {
		return nil, fmt.Errorf("response body too large")
	}
	return io.ReadAll(resp.Body)
}

func (h *httpRemote) String() string {
	var u = h.baseURL
	// avoid credentials leak
	if u.User != nil {
		un := u.User.Username()
		if un != "" {
			un = "***"
		}
		_, po := u.User.Password()
		if po {
			u.User = url.UserPassword(un, "***")
		} else {
			u.User = url.User(un)
		}
	}
	u.User = nil // avoid credentials leak
	return "HttpFetcher[baseURL=" + u.String() + "]"
}

var _ fmt.Stringer = (*httpRemote)(nil)
var _ M2Remote = (*httpRemote)(nil)

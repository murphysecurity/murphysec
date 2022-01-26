package maven

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func NewHttpRepo(baseUrl string) (*HttpRepo, error) {
	u, e := url.Parse(baseUrl)
	if e != nil {
		return nil, e
	}
	this := &HttpRepo{
		baseUrl:            u,
		httpClient:         http.DefaultClient,
		concurrencyLimiter: make(chan struct{}, 8),
	}
	logger.Debug.Println("New http repo", this)
	return this, nil
}

func MustNewHttpRepo(baseUrl string) *HttpRepo {
	h, e := NewHttpRepo(baseUrl)
	if e != nil {
		panic(e)
	}
	return h
}

type HttpRepo struct {
	baseUrl            *url.URL
	httpClient         *http.Client
	concurrencyLimiter chan struct{}
}

func (r *HttpRepo) String() string {
	return fmt.Sprintf("HttpRepo[%s]", r.baseUrl.String())
}

func (r *HttpRepo) FetchPomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error) {
	logger.Debug.Println("Repo fetch pom:", r, coordinate)
	select {
	case <-ctx.Done():
		logger.Info.Println("HttpRepo fetch cancelled", coordinate)
		return nil, ctx.Err()
	case r.concurrencyLimiter <- struct{}{}:
	}
	must.True(len(r.concurrencyLimiter) > 0)
	defer func() { <-r.concurrencyLimiter }()
	requestUrl := convCoordinateToRepoUrl(*r.baseUrl, coordinate).String()
	logger.Debug.Println("Send request:", requestUrl)
	request, e := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if e != nil {
		return nil, errors.Wrap(e, "create http request failed")
	}
	response, e := r.httpClient.Do(request)
	if e != nil {
		return nil, ErrShouldRetry.Decorate(e)
	}
	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, ErrArtifactNotFoundInRepo.Decorate(errors.Errorf("pom not found: %s", coordinate.String()))
	case http.StatusUnauthorized:
		return nil, ErrRepoAuthRequired.Decorate(errors.Errorf("can't fetch pom: %s", coordinate.String()))
	case http.StatusGatewayTimeout, http.StatusBadGateway, http.StatusTooManyRequests:
		return nil, ErrShouldRetry.Decorate(errors.Errorf("http status: %s", response.Status))
	case http.StatusOK:
	default:
		return nil, errors.New(fmt.Sprintf("http status: %s", response.Status))
	}
	// read body
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, ErrShouldRetry.Decorate(errors.Wrap(e, "read http body failed"))
	}
	pom, e := NewPomFileFromData(body)
	if e != nil {
		return nil, e
	}
	return pom, nil
}

func convCoordinateToRepoUrl(base url.URL, coordinate Coordinate) *url.URL {
	must.True(coordinate.GroupId != "" && coordinate.Version != "" && coordinate.ArtifactId != "")
	var seg []string
	seg = append(seg, base.Path)
	seg = append(seg, strings.Split(coordinate.GroupId, ".")...)
	seg = append(seg, coordinate.ArtifactId, coordinate.Version, coordinate.ArtifactId+"-"+coordinate.Version+".pom")
	base.Path = path.Join(seg...)
	return &base
}

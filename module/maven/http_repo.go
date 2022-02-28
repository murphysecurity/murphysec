package maven

import (
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"io"
	"murphysec-cli-simple/logger"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

type HttpRepo struct {
	baseUrl  url.URL
	m        map[Coordinate]chan struct{}
	l        sync.Mutex
	cache    map[Coordinate]*gopom.Project
	cacheErr map[Coordinate]error
}

func (r *HttpRepo) String() string {
	return fmt.Sprintf("HttpRepo[%s]", r.baseUrl.String())
}

func NewHttpRepo(baseUrl url.URL) *HttpRepo {
	return &HttpRepo{
		m:        map[Coordinate]chan struct{}{},
		l:        sync.Mutex{},
		cache:    map[Coordinate]*gopom.Project{},
		cacheErr: map[Coordinate]error{},
		baseUrl:  baseUrl,
	}
}

func (r *HttpRepo) Fetch(coordinate Coordinate) (*gopom.Project, error) {
	if !coordinate.Complete() {
		return nil, ErrInvalidCoordinate
	}
	r.l.Lock()
	if ch, ok := r.m[coordinate]; ok {
		r.l.Unlock()
		<-ch
		r.l.Lock()
		cp, ce := r.cache[coordinate], r.cacheErr[coordinate]
		r.l.Unlock()
		return cp, ce
	}
	ch := make(chan struct{})
	r.m[coordinate] = ch
	r.l.Unlock()
	defer func() { close(ch) }()
	//url := r.baseDir

	var u url.URL
	{
		u = r.baseUrl
		for _, it := range strings.Split(coordinate.GroupId, ".") {
			u.Path = path.Join(u.Path, it)
		}
		u.Path = path.Join(u.Path, coordinate.ArtifactId, coordinate.Version, fmt.Sprintf("%s-%s.pom", coordinate.ArtifactId, coordinate.Version))
	}
	logger.Info.Println("Request pom:", u.String())
	pom, e := fetchPom(u.String())
	r.l.Lock()
	r.cache[coordinate] = pom
	r.cacheErr[coordinate] = e
	r.l.Unlock()
	return pom, e
}

func fetchPom(url string) (*gopom.Project, error) {
	resp, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrArtifactNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http %d - %s", resp.StatusCode, resp.Status))
	}
	data, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, errors.New("read body failed")
	}
	var p gopom.Project
	if e := xml.Unmarshal(data, &p); e != nil {
		return nil, errors.Wrap(e, "parse pom failed")
	}
	return &p, nil
}

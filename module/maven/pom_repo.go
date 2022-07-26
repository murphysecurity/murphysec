package maven

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/vifraa/gopom"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type PomRepo interface {
	Fetch(coordinate Coordinate) (*UnresolvedPom, error)
}
type HttpRepo struct {
	baseUrl  url.URL
	m        map[Coordinate]chan struct{}
	l        sync.Mutex
	cache    map[Coordinate]*gopom.Project
	cacheErr map[Coordinate]error
	logger   *zap.Logger
}

func (r *HttpRepo) String() string {
	return fmt.Sprintf("HttpRepo[%s]", r.baseUrl.String())
}

func NewHttpRepo(ctx context.Context, baseUrl url.URL) *HttpRepo {
	return &HttpRepo{
		m:        map[Coordinate]chan struct{}{},
		l:        sync.Mutex{},
		cache:    map[Coordinate]*gopom.Project{},
		cacheErr: map[Coordinate]error{},
		baseUrl:  baseUrl,
		logger:   utils.UseLogger(ctx),
	}
}

func (r *HttpRepo) Fetch(coordinate Coordinate) (*UnresolvedPom, error) {
	logger := r.logger
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
		return &UnresolvedPom{cp}, ce
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
	logger.Sugar().Debugf("Request pom: %s", u.String())
	pom, e := fetchPom(u.String())
	r.l.Lock()
	r.cache[coordinate] = pom
	r.cacheErr[coordinate] = e
	r.l.Unlock()
	return &UnresolvedPom{pom}, e
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
		return nil, ErrArtifactNotFound.Detailed(fmt.Sprintf("HTTP %d - %s", resp.StatusCode, resp.Status))
	}
	data, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, ErrParsePomFailed.DetailedWrap("read body", e)
	}
	return parsePom(bytes.NewReader(data))
}

type LocalRepo struct {
	baseDir string
}

func NewLocalRepo(s string) *LocalRepo {
	return &LocalRepo{baseDir: s}
}

func (l *LocalRepo) String() string {
	return fmt.Sprintf("LocalRepo[%s]", l.baseDir)
}

func (l *LocalRepo) Fetch(coordinate Coordinate) (*UnresolvedPom, error) {
	if !coordinate.Complete() {
		return nil, ErrInvalidCoordinate
	}
	p := filepath.Join(l.baseDir)
	for _, s := range strings.Split(coordinate.GroupId, ".") {
		p = filepath.Join(p, s)
	}
	p = filepath.Join(p, coordinate.ArtifactId, coordinate.Version, coordinate.ArtifactId+"-"+coordinate.Version+".pom")
	if !utils.IsFile(p) {
		return nil, ErrArtifactNotFound
	}
	data, e := ioutil.ReadFile(p)
	if e != nil {
		return nil, ErrParsePomFailed.DetailedWrap("open pom", e)
	}
	var proj gopom.Project
	if e := xml.Unmarshal(data, &proj); e != nil {
		return nil, ErrParsePomFailed.Wrap(e)
	}
	return &UnresolvedPom{Project: &proj}, nil
}

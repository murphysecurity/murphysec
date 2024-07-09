package maven

import (
	"bytes"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"path"
	"strings"
	"sync"
)

type fetcher struct {
	remotes []M2Remote
	cache   fetcherCache
}

func newFetcher(remotes ...M2Remote) *fetcher {
	return &fetcher{
		remotes: remotes,
		cache:   fetcherCache{},
	}
}

func (f *fetcher) FetchPom(ctx context.Context, coordinate Coordinate) (*UnresolvedPom, error) {
	if a, b := f.cache.Get(coordinate); a != nil || b != nil {
		return a, b
	}
	a, b := f._fetchPom(ctx, coordinate)
	f.cache.Put(coordinate, a, b)
	return a, b
}

func (f *fetcher) _fetchPom(ctx context.Context, coordinate Coordinate) (*UnresolvedPom, error) {
	logger := logctx.Use(ctx).Sugar()
	basePath := path.Join(strings.Split(coordinate.GroupId, ".")...)
	basePath = path.Join(basePath, coordinate.ArtifactId, coordinate.Version)
	version := coordinate.Version
	for _, remote := range f.remotes {
		if coordinate.IsSnapshotVersion() {
			meta, e := fetchMetadata(ctx, coordinate, remote)
			if errors.Is(e, ErrRemoteNoResource) {
				continue
			}
			if e != nil {
				return nil, e
			}
			if v := meta.getPomVersionSnapshotSuffix(); v != "" {
				version = v
			}
		}
		target := path.Join(basePath, coordinate.ArtifactId+"-"+version+".pom")
		logger.Debugf("[%v]fetching pom %v", remote, coordinate)
		data, e := remote.GetPath(ctx, target)
		if e != nil {
			logger.Debugf("fetch %s %v", target, e)
		}
		if errors.Is(e, ErrRemoteNoResource) {
			continue
		}
		if e != nil {
			return nil, fmt.Errorf("fetch ")
		}
		pom, e := parsePom(bytes.NewReader(data))
		if e != nil {
			return nil, fmt.Errorf("parse pom: %s, %v", coordinate, e)
		}
		return &UnresolvedPom{Project: pom}, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrArtifactNotFound, coordinate)
}

func fetchMetadata(ctx context.Context, coordinate Coordinate, remote M2Remote) (*Metadata, error) {
	logger := logctx.Use(ctx).Sugar()
	basePath := path.Join(strings.Split(coordinate.GroupId, ".")...)
	basePath = path.Join(basePath, coordinate.ArtifactId, coordinate.Version)
	target := path.Join(basePath, "maven-metadata.xml")
	logger.Debugf("[%v]fetching metadata %v", remote, coordinate)
	data, e := remote.GetPath(ctx, target)
	if e != nil {
		logger.Debugf("fetch %s %v", coordinate, e)
		return nil, e
	}
	meta, e := parsePomVersionMeta(bytes.NewReader(data))
	if e != nil {
		return nil, fmt.Errorf("parse metadata failed, %v %v ", coordinate, e)
	}
	return meta, nil
}

type M2Remote interface {
	// GetPath might return any error is: ErrRemoteNoResource
	GetPath(ctx context.Context, path string) ([]byte, error)
}

type fetcherCacheItem struct {
	up *UnresolvedPom
	e  error
}

type fetcherCache struct {
	m     map[Coordinate]fetcherCacheItem
	mutex sync.Mutex
}

func (f *fetcherCache) Get(coordinate Coordinate) (*UnresolvedPom, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.m == nil {
		return nil, nil
	}
	r := f.m[coordinate]
	return r.up, r.e
}

func (f *fetcherCache) Put(coordinate Coordinate, pom *UnresolvedPom, e error) {
	if pom == nil && e == nil {
		panic("pom == nil && e == nil")
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.m == nil {
		f.m = map[Coordinate]fetcherCacheItem{}
	}
	f.m[coordinate] = fetcherCacheItem{pom, e}
}

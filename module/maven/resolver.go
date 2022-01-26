package maven

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"sync"
)

var ErrArtifactResolveFailed = errors.New("resolve artifact failed")

type Resolver struct {
	l               sync.Locker
	repoList        []Repo
	pomCache        map[Coordinate]*PomFile
	resolveErrCache map[Coordinate]error
	resolvingTask   map[Coordinate]bool
	cond            sync.Cond
}

func NewResolver(repos ...Repo) *Resolver {
	if len(repos) == 0 {
		repos = DefaultMavenRepo()
	}
	return &Resolver{
		repoList:        repos,
		pomCache:        map[Coordinate]*PomFile{},
		resolveErrCache: map[Coordinate]error{},
		resolvingTask:   map[Coordinate]bool{},
		cond:            *sync.NewCond(&sync.Mutex{}),
		l:               &sync.Mutex{},
	}
}

func (r *Resolver) AddRepo(repos ...Repo) {
	for _, it := range repos {
		r.repoList = append(r.repoList, it)
	}
}

func (r *Resolver) getPomFromCache(coordinate Coordinate) (*PomFile, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if e := r.resolveErrCache[coordinate]; e != nil {
		return nil, errors.Wrap(e, "cached error.")
	}
	if cachedPom := r.pomCache[coordinate]; cachedPom != nil {
		return cachedPom, nil
	}
	return nil, nil
}

func (r *Resolver) writePomToCache(pom *PomFile) {
	r.l.Lock()
	defer r.l.Unlock()
	r.pomCache[pom.Coordinate()] = pom
}

func (r *Resolver) writeErrToCache(coordinate Coordinate, e error) {
	r.l.Lock()
	defer r.l.Unlock()
	r.resolveErrCache[coordinate] = e
}

// 从各 repo 中获取 pom
func (r *Resolver) getPomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error) {
	must.True(coordinate.HasVersion())
	func() {
		r.cond.L.Lock()
		defer r.cond.L.Unlock()
		for r.resolvingTask[coordinate] {
			r.cond.Wait()
		}
		r.resolvingTask[coordinate] = true
	}()
	defer r.cond.Broadcast()
	defer func() {
		r.cond.L.Lock()
		defer r.cond.L.Unlock()
		delete(r.resolvingTask, coordinate)
	}()
	if p, e := r.getPomFromCache(coordinate); p != nil || e != nil {
		return p, e
	}
	for _, repo := range r.repoList {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		pom, e := repo.FetchPomFile(ctx, coordinate)
		if e != nil {
			if errors.Is(e, ErrArtifactNotFoundInRepo) {
				logger.Info.Printf("Couldn't found %v in %v\n", coordinate, repo)
				continue
			}
			r.writeErrToCache(coordinate, e)
			return nil, errors.Wrap(e, "unexpected error during artifact resolving")
		}
		r.writePomToCache(pom)
		return pom, nil
	}
	r.writeErrToCache(coordinate, ErrArtifactResolveFailed)
	return nil, ErrArtifactResolveFailed
}

// ResolvePomFile 解析pom和其父pom
func (r *Resolver) ResolvePomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error) {
	must.True(coordinate.HasVersion())
	if ctx == nil {
		ctx = context.Background()
	}
	pom, e := r.getPomFile(ctx, coordinate)
	if e != nil {
		return nil, e
	}
	if parentCoor := pom.parentCoordinate(); parentCoor != nil && pom.parentPom == nil {
		parentPom, e := r.ResolvePomFile(ctx, *parentCoor)
		if e != nil {
			logger.Info.Println(fmt.Sprintf("Resolve parent pom failed, artifact: %s, parent: %s, error: %v", coordinate, parentCoor, e))
		}
		pom.parentPom = parentPom
	}
	return pom, nil
}

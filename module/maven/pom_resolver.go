package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
)

type PomResolver struct {
	logger        *zap.Logger
	fetcher       *fetcher
	resolvedCache *resolvedPomCache
}

func NewPomResolver(ctx context.Context, remotes []M2Remote) *PomResolver {
	return &PomResolver{
		logger:        logctx.Use(ctx),
		fetcher:       newFetcher(remotes...),
		resolvedCache: newResolvedPomCache(),
	}
}

func (r *PomResolver) fetchPom(ctx context.Context, coordinate Coordinate) (*UnresolvedPom, error) {
	return r.fetcher.FetchPom(ctx, coordinate)
}

func (r *PomResolver) ResolvePom(ctx context.Context, coordinate Coordinate) (*Pom, error) {
	if pom, e := r.resolvedCache.get(coordinate); pom != nil || e != nil {
		return pom, e
	}
	c := newResolveContext(ctx)
	c.resolver = r
	pom, err := c._resolve(ctx, coordinate)
	if err != nil {
		r.resolvedCache.storeErr(coordinate, err)
	} else {
		r.resolvedCache.storePom(coordinate, pom)
	}
	return pom, err
}

func (r *PomResolver) addPom(module *UnresolvedPom) {
	r.fetcher.cache.Put(module.Coordinate(), module, nil)
}

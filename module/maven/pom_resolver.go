package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type PomResolver struct {
	logger        *zap.Logger
	repos         []PomRepo
	pomCache      *pomCache
	stats         *resolverStats
	resolvedCache *resolvedPomCache
}

func NewPomResolver(ctx context.Context) *PomResolver {
	return &PomResolver{
		logger:        utils.UseLogger(ctx),
		repos:         nil,
		pomCache:      newPomCache(),
		stats:         newResolverStats(),
		resolvedCache: newResolvedPomCache(),
	}
}

func (r *PomResolver) AddRepo(repo PomRepo) {
	r.repos = append(r.repos, repo)
}

func (r *PomResolver) fetchPom(coordinate Coordinate) (*UnresolvedPom, error) {
	r.stats.totalReq++
	if pom, e := r.pomCache.fetch(coordinate); pom != nil || e != nil {
		r.stats.cacheHit++
		return pom, e
	}
	logger := r.logger
	logger.Debug("Fetch pom", zap.Any("coordinate", coordinate))
	for _, repo := range r.repos {
		p, e := repo.Fetch(coordinate)
		if e == nil {
			r.pomCache.write(coordinate, p, nil)
			return p, nil
		}
		if errors.Is(e, ErrArtifactNotFound) {
			continue
		}
		logger.Sugar().Infof("Fetch %s from repo[%s] failed: %s", coordinate, repo, e)
	}
	r.pomCache.write(coordinate, nil, ErrArtifactNotFound)
	return nil, ErrArtifactNotFound
}

func (r *PomResolver) ResolvePom(ctx context.Context, coordinate Coordinate) (*Pom, error) {
	if pom, e := r.resolvedCache.get(coordinate); pom != nil || e != nil {
		return pom, e
	}
	c := newResolveContext(ctx)
	c.resolver = r
	pom, err := c._resolve(coordinate)
	if err != nil {
		r.resolvedCache.storeErr(coordinate, err)
	} else {
		r.resolvedCache.storePom(coordinate, pom)
	}
	return pom, err
}

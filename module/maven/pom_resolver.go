package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
)

type PomResolver struct {
	logger   *zap.Logger
	repos    []PomRepo
	pomCache *pomCache
}

func NewPomResolver(ctx context.Context) *PomResolver {
	return &PomResolver{
		logger:   utils.UseLogger(ctx),
		repos:    nil,
		pomCache: newPomCache(),
	}
}

func (r *PomResolver) AddRepo(repo PomRepo) {
	r.repos = append(r.repos, repo)
}

func (r *PomResolver) fetchPom(coordinate Coordinate) (*UnresolvedPom, error) {
	if pom, e := r.pomCache.fetch(coordinate); pom != nil || e != nil {
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
		logger.Sugar().Infof("Fetch %s from repo[%s] failed: %s", coordinate, repo, e)
	}
	r.pomCache.write(coordinate, nil, ErrArtifactNotFound)
	return nil, ErrArtifactNotFound
}

func (r *PomResolver) ResolvePom(ctx context.Context, coordinate Coordinate) (*Pom, error) {
	c := newResolveContext(ctx)
	c.resolver = r
	return c._resolve(coordinate)
}

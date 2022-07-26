package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/inlineproperty"
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

// fetchAllPomIncludeParents 返回当前pom+所有父级, [current, current.parent, current.parent.parent, ...]
func (r *PomResolver) fetchAllPomIncludeParents(p *UnresolvedPom) []*UnresolvedPom {
	logger := r.logger
	var poms []*UnresolvedPom
	poms = append(poms, p)
	var visitedParents = make(map[Coordinate]struct{})
	parentCoordinate := p.ParentCoordinate()
	for parentCoordinate != nil {
		if _, ok := visitedParents[*parentCoordinate]; ok {
			break
		}
		visitedParents[*parentCoordinate] = struct{}{}
		parent, e := r.fetchPom(*parentCoordinate)
		if e != nil {
			logger.Warn("Fetch parent failed", zap.Any("coordinate", *parentCoordinate), zap.Error(e))
			break
		}
		poms = append(poms, parent)
		parentCoordinate = parent.ParentCoordinate()
	}
	return poms
}

func (r *PomResolver) _resolve(coordinate Coordinate, visitedCoordinate map[Coordinate]struct{}) (*Pom, error) {
	if coordinate.IsBad() {
		return nil, ErrBadCoordinate.Detailed(coordinate.String())
	}
	if _, ok := visitedCoordinate[coordinate]; ok {
		return nil, ErrPomCircularDependent.Detailed(coordinate.String())
	}
	visitedCoordinate[coordinate] = struct{}{}
	defer delete(visitedCoordinate, coordinate)

	rawPom, e := r.fetchPom(coordinate)
	if e != nil {
		return nil, e
	}
	resolved, e := r.resolveParent(rawPom)
	if e != nil {
		return nil, e
	}
	if e := r.resolveDependencyManagementImport(resolved, visitedCoordinate); e != nil {
		return nil, e
	}

	// merge dependencyManagement to dependencies
	resolved.depSet.mergeDependencyManagement(resolved.depmSet)

	return resolved, nil
}

func (r *PomResolver) resolveParent(p *UnresolvedPom) (*Pom, error) {
	var logger = r.logger
	resolved := &Pom{
		project:    p.Project,
		properties: inlineproperty.New(),
	}

	var poms []*UnresolvedPom
	poms = r.fetchAllPomIncludeParents(p)
	// reverse-iterate all parents
	utils.Reverse(poms)

	var depSet = newPomDependencySet()
	var depmSet = newPomDependencySet()
	resolved.depSet = depSet
	resolved.depmSet = depmSet
	for _, pom := range poms {
		project := pom.Project

		// merge property
		resolved.properties.PutMap(project.Properties.Entries)
		for _, profile := range project.Profiles {
			resolved.properties.PutMap(profile.Properties.Entries)
		}

		// merge dependency management
		depmSet.mergeDepsSlice(project.DependencyManagement.Dependencies)
		for _, profile := range project.Profiles {
			depmSet.mergeDepsSlice(profile.DependencyManagement.Dependencies)
		}

		// merge dependencies
		depSet.mergeDepsSlice(project.Dependencies)
		for _, profile := range project.Profiles {
			depSet.mergeDepsSlice(profile.Dependencies)
		}

	}
	coor := resolved.Coordinate()
	if !coor.Complete() {
		logger.Warn("Coordinate incomplete", zap.Any("coordinate", coor))
	}
	// merge ${project.*} ${${project.artifactId.*}}
	if coor.Complete() {
		resolved.properties.PutIfAbsent("project.version", coor.Version)
		resolved.properties.PutIfAbsent("project.artifactId", coor.ArtifactId)
		resolved.properties.PutIfAbsent("project.groupId", coor.GroupId)
		resolved.properties.PutIfAbsent(fmt.Sprintf("%s.groupId", coor.ArtifactId), coor.GroupId)
		resolved.properties.PutIfAbsent(fmt.Sprintf("%s.version", coor.ArtifactId), coor.Version)
	}
	depSet.mergeProperty(resolved.properties)
	depmSet.mergeProperty(resolved.properties)

	return resolved, nil
}

func (r *PomResolver) resolveDependencyManagementImport(pom *Pom, visited map[Coordinate]struct{}) error {
	var logger = r.logger
	for _, dependency := range pom.depmSet.listDeps() {
		if dependency.Scope != "import" {
			continue
		}
		np, e := r._resolve(Coordinate{dependency.GroupID, dependency.ArtifactID, dependency.Version}, visited)
		if e != nil {
			logger.Warn("Resolve dependencyManagement failed", zap.Error(e))
			continue
		}
		pom.depmSet.mergeDepsSlice(np.depmSet.listDeps())
	}
	return nil
}

func (r *PomResolver) ResolvePom(coordinate Coordinate) (*Pom, error) {
	pom, e := r._resolve(coordinate, map[Coordinate]struct{}{})
	if e != nil {
		return nil, e
	}
	return pom, nil
}

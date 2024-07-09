package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/vifraa/gopom"
	"go.uber.org/zap"
)

type pomBuilder struct {
	coordinate       Coordinate
	parentCoordinate *Coordinate
	properties       *properties
	deps             *pomDependencySet
	depms            *pomDependencySet
	pom              *UnresolvedPom
}

func newPomBuilder() *pomBuilder {
	return &pomBuilder{
		properties: newProperties(),
		deps:       newPomDependencySet(),
		depms:      newPomDependencySet(),
	}
}

// listDependencyManagements 返回已解析属性的依赖管理列表
func (p *pomBuilder) listDependencyManagements() []gopom.Dependency {
	var rs []gopom.Dependency
	for _, dep := range p.depms.listAll() {
		d := gopom.Dependency{
			GroupID:    p.properties.Resolve(dep.GroupID),
			ArtifactID: p.properties.Resolve(dep.ArtifactID),
			Version:    p.properties.Resolve(dep.Version),
			Type:       p.properties.Resolve(dep.Type),
			Classifier: p.properties.Resolve(dep.Classifier),
			Scope:      p.properties.Resolve(dep.Scope),
			SystemPath: p.properties.Resolve(dep.SystemPath),
			Exclusions: dep.Exclusions,
			Optional:   dep.Optional,
		}
		rs = append(rs, d)
	}
	return rs
}

func (p *pomBuilder) build() *Pom {
	return &Pom{
		dir:        "",
		project:    p.pom.Project,
		depSet:     p.deps,
		depmSet:    p.depms,
		Coordinate: p.coordinate,
		properties: p.properties,
	}
}

type resolveContext struct {
	logger   *zap.Logger
	visited  map[Coordinate]struct{}
	resolver *PomResolver
}

func newResolveContext(ctx context.Context) *resolveContext {
	return &resolveContext{
		logger:  logctx.Use(ctx),
		visited: map[Coordinate]struct{}{},
	}
}

func (r *resolveContext) visitEnter(coordinate Coordinate) bool {
	if _, ok := r.visited[coordinate]; ok {
		return true
	}
	return false
}
func (r *resolveContext) visitExit(coordinate Coordinate) {
	delete(r.visited, coordinate)
}

func (r *resolveContext) _resolve(ctx context.Context, coordinate Coordinate) (*Pom, error) {
	if r.visitEnter(coordinate) {
		return nil, ErrPomCircularDependent.Detailed(coordinate.String())
	}
	defer r.visitExit(coordinate)
	builder := newPomBuilder()
	rawPom, e := r.resolver.fetchPom(ctx, coordinate)
	if e != nil {
		return nil, e
	}
	builder.pom = rawPom
	builder.parentCoordinate = rawPom.ParentCoordinate()
	r.resolveInheritance(ctx, builder)
	_ = r.resolveCoordinate(builder)
	// resolve & merge dependencyManagement.type==import
	r.resolveDependencyManagementImport(ctx, builder)
	// merge dependencyManagement into dependencies
	builder.deps.mergeAll(builder.depms.listAll(), false, true)
	return builder.build(), nil
}

func (r *resolveContext) resolveDependencyManagementImport(ctx context.Context, builder *pomBuilder) {
	for _, dependency := range builder.listDependencyManagements() {
		if dependency.Scope != "import" {
			continue
		}
		coordinate := Coordinate{dependency.GroupID, dependency.ArtifactID, dependency.Version}
		np, e := r._resolve(ctx, coordinate)
		if e != nil {
			continue
		}
		builder.depms.mergeAll(np.ListDependencyManagements(), false, false)
	}
}

func (r *resolveContext) resolveInheritance(ctx context.Context, builder *pomBuilder) {
	logger := r.logger
	coordinate := builder.parentCoordinate
	var visitedParents = make(map[Coordinate]struct{})
	var poms []*UnresolvedPom
	poms = append(poms, builder.pom)
	for coordinate != nil {
		if _, ok := visitedParents[*coordinate]; ok {
			break
		}
		visitedParents[*coordinate] = struct{}{}
		parent, e := r.resolver.fetchPom(ctx, *coordinate)
		if e != nil {
			logger.Warn("Fetch parent failed", zap.Any("coordinate", *coordinate), zap.Error(e))
			break
		}
		poms = append(poms, parent)
		coordinate = parent.ParentCoordinate()
	}

	utils.Reverse(poms)
	for _, pom := range poms {
		project := pom.Project

		// merge property
		builder.properties.PutMap(project.Properties.Entries)

		// merge dependency management
		builder.depms.mergeAll(project.DependencyManagement.Dependencies, false, false)

		// merge dependencies
		builder.deps.mergeAll(project.Dependencies, true, false)
	}

	// merge ${project.parent.*}
	if c := builder.parentCoordinate; c != nil {
		builder.properties.PutIfAbsent("project.parent.version", c.Version)
		builder.properties.PutIfAbsent("project.parent.artifactId", c.ArtifactId)
		builder.properties.PutIfAbsent("project.parent.groupId", c.GroupId)
	}
}

func (r *resolveContext) resolveCoordinate(builder *pomBuilder) error {
	pf := builder.pom.Project
	g := pf.GroupID
	if g == "" {
		g = pf.Parent.GroupID
	}
	a := pf.ArtifactID
	if a == "" {
		a = pf.Parent.ArtifactID
	}
	v := pf.Version
	if v == "" {
		v = pf.Parent.Version
	}
	coordinate := Coordinate{
		GroupId:    builder.properties.Resolve(g),
		ArtifactId: builder.properties.Resolve(a),
		Version:    builder.properties.Resolve(v),
	}
	builder.coordinate = coordinate
	// merge ${project.*} ${${project.artifactId.*}}
	builder.properties.PutIfAbsent("project.version", coordinate.Version)
	builder.properties.PutIfAbsent("project.artifactId", coordinate.ArtifactId)
	builder.properties.PutIfAbsent("project.groupId", coordinate.GroupId)
	builder.properties.PutIfAbsent(fmt.Sprintf("%s.groupId", coordinate.ArtifactId), coordinate.GroupId)
	builder.properties.PutIfAbsent(fmt.Sprintf("%s.version", coordinate.ArtifactId), coordinate.Version)
	return nil
}

type resolvedPomCache struct {
	m  map[Coordinate]*Pom
	em map[Coordinate]error
}

func newResolvedPomCache() *resolvedPomCache {
	return &resolvedPomCache{
		m:  map[Coordinate]*Pom{},
		em: map[Coordinate]error{},
	}
}

func (r *resolvedPomCache) storeErr(coordinate Coordinate, err error) {
	r.em[coordinate] = err
}

func (r *resolvedPomCache) storePom(coordinate Coordinate, pom *Pom) {
	r.m[coordinate] = pom
}

func (r *resolvedPomCache) get(coordinate Coordinate) (*Pom, error) {
	if e, ok := r.em[coordinate]; ok {
		return nil, e
	}
	if p, ok := r.m[coordinate]; ok {
		return p, nil
	}
	return nil, nil
}

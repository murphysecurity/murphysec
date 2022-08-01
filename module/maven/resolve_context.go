package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
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

func (p *pomBuilder) build() *Pom {
	return &Pom{
		dir:        "",
		project:    p.pom.Project,
		depSet:     p.deps,
		depmSet:    p.depms,
		Coordinate: p.coordinate,
	}
}

type resolveContext struct {
	logger   *zap.Logger
	visited  map[Coordinate]struct{}
	resolver *PomResolver
}

func newResolveContext(ctx context.Context) *resolveContext {
	return &resolveContext{
		logger:  utils.UseLogger(ctx),
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

func (r *resolveContext) _resolve(coordinate Coordinate) (*Pom, error) {
	if coordinate.IsBad() {
		return nil, ErrBadCoordinate.Detailed(coordinate.String())
	}
	if r.visitEnter(coordinate) {
		return nil, ErrPomCircularDependent.Detailed(coordinate.String())
	}
	defer r.visitExit(coordinate)
	builder := newPomBuilder()
	rawPom, e := r.resolver.fetchPom(coordinate)
	if e != nil {
		return nil, e
	}
	builder.pom = rawPom
	builder.parentCoordinate = rawPom.ParentCoordinate()
	r.resolveInheritance(builder)
	if e := r.resolveCoordinate(builder); e != nil {
		return nil, e
	}
	builder.depms.mergeProperty(builder.properties)
	builder.deps.mergeProperty(builder.properties)
	r.resolveDependencyManagementImport(builder)
	builder.deps.mergeDependencyManagement(builder.depms)
	return builder.build(), nil
}

func (r *resolveContext) resolveDependencyManagementImport(builder *pomBuilder) {
	var logger = r.logger
	for _, dependency := range builder.depms.listDeps() {
		if dependency.Scope != "import" {
			continue
		}
		np, e := r._resolve(Coordinate{dependency.GroupID, dependency.ArtifactID, dependency.Version})
		if e != nil {
			logger.Warn("Resolve dependencyManagement failed", zap.Error(e))
			continue
		}
		builder.depms.mergeDepsSlice(np.depmSet.listDeps())
	}
}

func (r *resolveContext) resolveInheritance(builder *pomBuilder) {
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
		parent, e := r.resolver.fetchPom(*coordinate)
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
		//for _, profile := range project.Profiles {
		//	builder.properties.PutMap(profile.Properties.Entries)
		//}

		// merge dependency management
		builder.depms.mergeDepsSlice(project.DependencyManagement.Dependencies)
		//for _, profile := range project.Profiles {
		//	builder.depms.mergeDepsSlice(profile.DependencyManagement.Dependencies)
		//}

		// merge dependencies
		builder.deps.mergeDepsSlice(project.Dependencies)
		//for _, profile := range project.Profiles {
		//	builder.deps.mergeDepsSlice(profile.Dependencies)
		//}
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
	if !coordinate.Complete() {
		return ErrCouldNotResolve.Detailed(fmt.Sprintf("bad coordinate: %s", coordinate))
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

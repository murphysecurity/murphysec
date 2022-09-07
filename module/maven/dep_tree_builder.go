package maven

import (
	"context"
	list "github.com/bahlo/generic-list-go"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/vifraa/gopom"
	"go.uber.org/zap"
)

func BuildDepTree(ctx context.Context, resolver *PomResolver, coordinate Coordinate) *Dependency {
	analyzer := &depAnalyzer{
		Context:       ctx,
		resolver:      resolver,
		logger:        utils.UseLogger(ctx),
		versionChosen: map[string]string{},
	}
	tree := analyzer.analyze(coordinate)
	//VersionReconciling(ctx, tree)
	return tree
}

type depAnalyzer struct {
	context.Context
	resolver *PomResolver
	logger   *zap.Logger
	// full resolved item should be skipped, to reduce the large tree
	fullResolved map[Coordinate]string
	// chosen version shouldn't be changed again
	versionChosen map[string]string
}

func (d *depAnalyzer) analyze(coordinate Coordinate) *Dependency {
	var logger = d.logger
	type item struct {
		Coordinate
		Children             []*item
		Exclusion            *exclusionMap
		DependencyManagement *dependencyManagementMap
		Parent               *item
	}

	var _convToTree func(it *item) *Dependency
	_convToTree = func(it *item) *Dependency {
		if it == nil {
			return nil
		}
		d := &Dependency{
			Coordinate: it.Coordinate,
			Children:   []Dependency{},
		}
		for _, it := range it.Children {
			if it == nil {
				continue
			}
			r := _convToTree(it)
			if r == nil {
				continue
			}
			d.Children = append(d.Children, *r)
		}
		return d
	}

	q := list.New[*item]()
	r := &item{Coordinate: coordinate}
	q.PushBack(r)

outer:
	for q.Len() > 0 {
		cur := q.Front().Value
		q.Remove(q.Front())

		// circular
		{
			var p = cur.Parent
			for p != nil {
				if p.Coordinate == cur.Coordinate {
					continue outer
				}
				p = p.Parent
			}
		}

		pom, e := d.resolver.ResolvePom(d.Context, cur.Coordinate)
		if e != nil {
			logger.Warn("Resolve dependency failed", zap.Error(e), zap.Any("coordinate", cur.Coordinate))
			continue
		}
		dm := newDependencyManagementMap(cur.DependencyManagement, pom.ListDependencyManagements())
		for _, dep := range pom.ListDependencies() {
			if cur.Exclusion.Has(dep.GroupID, dep.ArtifactID) {
				continue
			}
			depCoordinate := Coordinate{GroupId: dep.GroupID, ArtifactId: dep.ArtifactID}
			verKey := dep.GroupID + dep.ArtifactID
			if v := d.versionChosen[verKey]; v != "" {
				depCoordinate.Version = v
			} else if dep.Version != "" {
				depCoordinate.Version = dep.Version
				d.versionChosen[verKey] = dep.Version
			} else if v := dm.GetVersionOf(dep.GroupID, dep.ArtifactID); v != "" {
				depCoordinate.Version = v
				d.versionChosen[verKey] = v
			} else {
				logger.Warn("Resolution version failed", zap.Any("in", coordinate), zap.String("dep", dep.GroupID+":"+dep.ArtifactID))
				continue
			}
			child := &item{
				Coordinate:           depCoordinate,
				Exclusion:            newExclusionMap(cur.Exclusion, dep.Exclusions),
				DependencyManagement: dm,
				Parent:               cur,
			}
			cur.Children = append(cur.Children, child)
			q.PushBack(child)
		}
	}
	return _convToTree(r)
}

type dependencyManagementMap struct {
	parent *dependencyManagementMap
	m      map[string]string
}

func newDependencyManagementMap(parent *dependencyManagementMap, dm []gopom.Dependency) *dependencyManagementMap {
	r := &dependencyManagementMap{
		parent: parent,
		m:      map[string]string{},
	}
	for _, it := range dm {
		r.m[it.GroupID+it.ArtifactID] = it.Version
	}
	return r
}

func (m *dependencyManagementMap) GetVersionOf(groupId string, artifactId string) string {
	if m == nil {
		return ""
	}
	if m.parent != nil {
		if v := m.parent.GetVersionOf(groupId, artifactId); v != "" {
			return v
		}
	}
	return m.m[groupId+artifactId]
}

type exclusionMap struct {
	parent    *exclusionMap
	exclusion map[string]struct{}
}

func newExclusionMap(parent *exclusionMap, exclusions []gopom.Exclusion) *exclusionMap {
	r := &exclusionMap{
		parent:    parent,
		exclusion: map[string]struct{}{},
	}
	for _, it := range exclusions {
		r.exclusion[it.GroupID+it.ArtifactID] = struct{}{}
	}
	return r
}

func (m *exclusionMap) Has(groupId string, artifactId string) bool {
	if m == nil {
		return false
	}
	if _, ok := m.exclusion[groupId+artifactId]; ok {
		return true
	}
	if m.parent != nil {
		return m.parent.Has(groupId, artifactId)
	}
	return false
}

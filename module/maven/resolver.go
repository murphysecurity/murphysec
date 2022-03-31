package maven

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Resolver struct {
	repos         []Repo
	baseDir       string
	pomCache      map[Coordinate]*gopom.Project
	sy            sync.Mutex
	resolvedCache map[Coordinate]*PomFile
}

func NewResolver() *Resolver {
	r := &Resolver{
		repos:         nil,
		baseDir:       "",
		pomCache:      map[Coordinate]*gopom.Project{},
		sy:            sync.Mutex{},
		resolvedCache: map[Coordinate]*PomFile{},
	}
	op := ReadMvnOption()
	r.repos = append(r.repos, NewLocalRepo(op.LocalRepoPath))
	for _, s := range op.Remote {
		u, e := url.Parse(s)
		if e != nil {
			continue
		}
		r.repos = append(r.repos, NewHttpRepo(*u))
	}
	return r
}

var ErrCouldNotResolve = errors.New("ErrCouldNotResolve")

func (r *Resolver) fetchLocalPom(coordinate Coordinate, dir string) *gopom.Project {
	if !coordinate.Complete() {
		return nil
	}
	pom, e := gopom.Parse(filepath.Join(dir, "pom.xml"))
	if e != nil {
		return nil
	}
	if pom == nil {
		return nil
	}
	coor := Coordinate{pom.GroupID, pom.ArtifactID, pom.Version}
	if coor != coordinate {
		return nil
	}
	return pom
}

func (r *Resolver) fetchPom(coordinate Coordinate) (*gopom.Project, error) {
	for _, rep := range r.repos {
		{
			r.sy.Lock()
			p := r.pomCache[coordinate]
			r.sy.Unlock()
			if p != nil {
				logger.Debug.Println(coordinate, "from cache")
				return p, nil
			}
		}
		logger.Debug.Println("fetch", coordinate, "from", rep)
		pom, e := rep.Fetch(coordinate)
		if e == ErrArtifactNotFound {
			logger.Info.Println("not found", coordinate, "from", rep)
			continue
		}
		if e != nil {
			logger.Info.Println("fetch pom failed", coordinate, e.Error())
			return nil, e
		}
		r.sy.Lock()
		r.pomCache[coordinate] = pom
		r.sy.Unlock()
		return pom, nil
	}
	logger.Info.Println("couldn't resolve pom", coordinate)
	return nil, ErrCouldNotResolve
}

func (r *Resolver) ResolveByCoordinate(coordinate Coordinate) *PomFile {
	coordinate = coordinate.Normalize()
	if r.resolvedCache[coordinate] != nil {
		return r.resolvedCache[coordinate]
	}
	d, e := r.fetchPom(coordinate)
	if e != nil {
		return nil
	}
	pf := r.Resolve(NewPomBuilder(d), nil)
	if pf == nil {
		return nil
	}
	r.resolvedCache[coordinate] = pf
	return pf
}

func (r *Resolver) ResolveLocally(builder *PomBuilder, visitedCoor map[Coordinate]struct{}) *PomFile {
	if builder == nil {
		return nil
	}
	parentCoor := Coordinate{
		GroupId:    builder.P.Parent.GroupID,
		ArtifactId: builder.P.Parent.ArtifactID,
		Version:    builder.P.Parent.Version,
	}.Normalize()

	{
		// circle detect
		if visitedCoor == nil {
			visitedCoor = map[Coordinate]struct{}{}
		}
		if _, ok := visitedCoor[parentCoor]; ok {
			logger.Warn.Println("circular inheritance pom")
			return nil
		}
		visitedCoor[parentCoor] = struct{}{}
		defer delete(visitedCoor, parentCoor)
	}

	// resolve parent
	if parentCoor.Complete() {
		// pom relative path
		var parentPath string
		if builder.Path != "" {
			if builder.P.Parent.RelativePath == "" {
				parentPath = filepath.Join(builder.Path, "../")
			} else {
				parentPath = filepath.Join(builder.Path, builder.P.Parent.RelativePath)
			}
		}
		var parentPom *gopom.Project
		if parentPath != "" {
			parentPom = r.fetchLocalPom(parentCoor, parentPath)
		}

		if parentPom != nil {
			parentBuilder := NewPomBuilder(parentPom)
			parentBuilder.Path = parentPath
			builder.ParentPom = r.Resolve(parentBuilder, visitedCoor)
		}
	}
	pf := builder.Build()
	return pf
}

// Resolve 递归解析 pom，先查本地再查远程端
func (r *Resolver) Resolve(builder *PomBuilder, visitedCoor map[Coordinate]struct{}) *PomFile {
	if builder == nil {
		return nil
	}
	parentCoor := Coordinate{
		GroupId:    builder.P.Parent.GroupID,
		ArtifactId: builder.P.Parent.ArtifactID,
		Version:    builder.P.Parent.Version,
	}.Normalize()

	{
		// circle detect
		if visitedCoor == nil {
			visitedCoor = map[Coordinate]struct{}{}
		}
		if _, ok := visitedCoor[parentCoor]; ok {
			logger.Warn.Println("circular inheritance pom")
			return nil
		}
		visitedCoor[parentCoor] = struct{}{}
		defer delete(visitedCoor, parentCoor)
	}

	// resolve parent
	if parentCoor.Complete() {
		// pom relative path
		var parentPath string
		if builder.Path != "" {
			if builder.P.Parent.RelativePath == "" {
				parentPath = filepath.Join(builder.Path, "../")
			} else {
				parentPath = filepath.Join(builder.Path, builder.P.Parent.RelativePath)
			}
		}

		var parentPom *gopom.Project
		if parentPath != "" {
			parentPom = r.fetchLocalPom(parentCoor, parentPath)
		}
		if parentPom == nil {
			pp, e := r.fetchPom(parentCoor)
			if e != nil {
				logger.Info.Println("resolve parent failed, parent:", parentCoor, e.Error())
			}
			parentPom = pp
		}

		if parentPom != nil {
			parentBuilder := NewPomBuilder(parentPom)
			parentBuilder.Path = parentPath
			builder.ParentPom = r.Resolve(parentBuilder, visitedCoor)
		}
	}

	pf := builder.Build()
	r.resolvedCache[pf.coordinate] = pf
	return pf
}

type PomBuilder struct {
	P *gopom.Project
	// Path pom文件路径
	Path      string
	ParentPom *PomFile
}

func (p *PomBuilder) Build() *PomFile {
	pf := NewPomFile()
	pf.parentPom = p.ParentPom
	pf.path = p.Path
	{
		// fill properties
		m := map[string]string{}
		if pf.parentPom != nil {
			for k, v := range pf.parentPom.propertyMap {
				m[k] = v
			}
		}
		for k, v := range p.P.Properties.Entries {
			m[k] = v
		}
		resolved := map[string]string{}
		for k := range m {
			resolved[k] = _resolveProperty(m, nil, k)
		}
		pf.propertyMap = resolved
	}
	{
		c := Coordinate{
			GroupId:    pf.property(p.P.GroupID),
			ArtifactId: pf.property(p.P.ArtifactID),
			Version:    pf.property(p.P.Version),
		}.Normalize()
		if c.GroupId == "" {
			c.GroupId = strings.TrimSpace(pf.property(p.P.Parent.GroupID))
		}
		if c.Version == "" {
			c.Version = strings.TrimSpace(pf.property(p.P.Parent.Version))
		}
		pf.coordinate = c.Normalize()
	}
	{
		// fill dependencyManagement
		if pf.parentPom != nil {
			for k, v := range pf.parentPom.dependencyManagement {
				must.True(!k.HasVersion())
				pf.dependencyManagement[k] = v
			}
		}
		for _, it := range p.P.DependencyManagement.Dependencies {
			coor := Coordinate{
				GroupId:    pf.property(it.GroupID),
				ArtifactId: pf.property(it.ArtifactID),
			}
			if coor.IsBad() {
				continue
			}
			pf.dependencyManagement[coor.Normalize()] = strings.TrimSpace(pf.property(it.Version))
		}
	}
	{
		// fill dependencies
		type id struct {
			g string
			a string
		}
		m := map[id]PomDependencyItem{}
		if pf.parentPom != nil {
			for _, it := range pf.parentPom.dependencies {
				if it.IsBad() {
					continue
				}
				m[id{strings.TrimSpace(it.GroupId), strings.TrimSpace(it.ArtifactId)}] = it
			}
		}
		for _, it := range p.P.Dependencies {
			if it.Scope != "" && strings.TrimSpace(it.Scope) != "compile" {
				continue
			}
			if strings.TrimSpace(it.Optional) == "true" {
				continue
			}
			groupId := strings.TrimSpace(pf.property(it.GroupID))
			artifactId := strings.TrimSpace(pf.property(it.ArtifactID))
			if groupId == "" || artifactId == "" {
				continue
			}
			version := strings.TrimSpace(pf.property(it.Version))
			coor := Coordinate{GroupId: groupId, ArtifactId: artifactId}
			if coor.IsBad() {
				continue
			}
			if version == "" {
				version = pf.dependencyManagement[coor]
			}
			depItem := PomDependencyItem{
				Coordinate: Coordinate{groupId, artifactId, version}.Normalize(),
				Scope:      strings.TrimSpace(it.Scope),
			}
			if !depItem.Complete() {
				continue
			}
			m[id{groupId, artifactId}] = depItem
		}
		for _, v := range m {
			pf.dependencies = append(pf.dependencies, v)
		}
	}
	return pf
}

func _resolveProperty(m map[string]string, path map[string]struct{}, key string) string {
	if path == nil {
		path = map[string]struct{}{}
	}
	if _, ok := path[key]; ok {
		return fmt.Sprintf("${%s}", key)
	} else {
		path[key] = struct{}{}
		defer delete(path, key)
	}
	s, ok := m[key]
	if !ok {
		return fmt.Sprintf("${%s}", key)
	}
	seg := pomInlineParameterPattern.Split(s, -1)
	rs := []string{seg[0]}
	matches := pomInlineParameterPattern.FindAllStringSubmatch(s, -1)
	for index, match := range matches {
		nk := match[1]
		rs = append(rs, _resolveProperty(m, path, nk))
		rs = append(rs, seg[index+1])
	}
	return strings.Join(rs, "")
}

var pomInlineParameterPattern = regexp.MustCompile("\\$\\{([^{}]+)\\}")

func NewPomBuilder(p *gopom.Project) *PomBuilder {
	if p == nil {
		panic("p is nil")
	}
	return &PomBuilder{
		P:         p,
		Path:      "",
		ParentPom: nil,
	}
}

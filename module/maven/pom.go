package maven

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/vifraa/gopom"
	"golang.org/x/text/encoding/ianaindex"
	"io"
	"os"
)

type UnresolvedPom struct {
	Project *gopom.Project
	Path    string
}

func (u UnresolvedPom) Coordinate() Coordinate {
	c := Coordinate{
		GroupId:    u.Project.GroupID,
		ArtifactId: u.Project.ArtifactID,
		Version:    u.Project.Version,
	}
	if c.GroupId == "" {
		c.GroupId = u.Project.Parent.GroupID
	}
	if c.ArtifactId == "" {
		c.ArtifactId = u.Project.Parent.ArtifactID
	}
	if c.Version == "" || c.Version == "${parent.version}" || c.Version == "${project.parent.version}" {
		c.Version = u.Project.Parent.Version
	}
	return c
}

func (u UnresolvedPom) ParentCoordinate() *Coordinate {
	p := u.Project.Parent
	c := Coordinate{GroupId: p.GroupID, ArtifactId: p.ArtifactID, Version: p.Version}
	if !c.Complete() {
		return nil
	}
	return &c
}

type Pom struct {
	dir     string
	project *gopom.Project
	// dependencies
	depSet *pomDependencySet
	// dependencyManagement
	depmSet *pomDependencySet
	Coordinate
	properties *properties
}

// ListDependencies 返回全部已解析属性的依赖
func (p *Pom) ListDependencies() (rs []gopom.Dependency) {
	for _, dep := range p.depSet.listAll() {
		r := p.resolveDependencyProperty(dep)
		if r.Optional == "true" {
			continue
		}
		rs = append(rs, r)
	}
	return
}

// ListDependencyManagements 返回全部已解析属性的依赖管理
func (p *Pom) ListDependencyManagements() (rs []gopom.Dependency) {
	for _, dep := range p.depmSet.listAll() {
		r := p.resolveDependencyProperty(dep)
		if r.Optional == "true" {
			continue
		}
		rs = append(rs, r)
	}
	return
}

func (p Pom) resolveDependencyProperty(dep gopom.Dependency) gopom.Dependency {
	return gopom.Dependency{
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
}

func (p Pom) ParentCoordinate() *Coordinate {
	f := p.project.Parent
	c := Coordinate{GroupId: f.GroupID, ArtifactId: f.ArtifactID, Version: f.Version}
	if !c.Complete() {
		return nil
	}
	return &c
}

func readPomFile(ctx context.Context, path string) (*gopom.Project, error) {
	logger := logctx.Use(ctx)
	logger.Sugar().Debugf("Read pom: %s", path)
	data, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	logger.Sugar().Debugf("Parse pom: %s", path)
	project, e := parsePom(bytes.NewReader(data))
	if e != nil {
		return nil, e
	}
	return project, nil
}

func parsePom(reader io.Reader) (*gopom.Project, error) {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		if enc == nil {
			// Assume it's compatible with (a subset of) UTF-8 encoding
			// Bug: https://github.com/golang/go/issues/19421
			return reader, nil
		}
		return enc.NewDecoder().Reader(reader), nil
	}
	var project gopom.Project
	if e := decoder.Decode(&project); e != nil {
		return nil, ErrParsePomFailed.Wrap(e)
	}
	return &project, nil
}

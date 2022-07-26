package maven

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/inlineproperty"
	"github.com/vifraa/gopom"
	"golang.org/x/text/encoding/ianaindex"
	"io"
	"os"
)

type UnresolvedPom struct {
	Project *gopom.Project
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
	if c.Version == "" {
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
	dir        string
	properties *inlineproperty.Properties
	project    *gopom.Project
	// dependencies
	depSet *pomDependencySet
	// dependencyManagement
	depmSet *pomDependencySet
}

func (p *Pom) ListDeps() []gopom.Dependency {
	if p.depSet == nil {
		return nil
	}
	return p.depSet.listDeps()
}

func (p Pom) ParentCoordinate() *Coordinate {
	f := p.project.Parent
	c := Coordinate{GroupId: f.GroupID, ArtifactId: f.ArtifactID, Version: f.Version}
	if !c.Complete() {
		return nil
	}
	return &c
}

func (p Pom) Coordinate() Coordinate {
	pf := p.project
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
	return Coordinate{
		GroupId:    p.properties.Resolve(g),
		ArtifactId: p.properties.Resolve(a),
		Version:    p.properties.Resolve(v),
	}
}

func readPomFile(ctx context.Context, path string) (*gopom.Project, error) {
	logger := utils.UseLogger(ctx)
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

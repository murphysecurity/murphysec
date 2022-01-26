package maven

import (
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"golang.org/x/net/html/charset"
	"murphysec-cli-simple/logger"
	"regexp"
	"strings"
)

type PomFile struct {
	pom       gopom.Project
	parentPom *PomFile
}

func NewPomFileFromData(data []byte) (*PomFile, error) {
	var p gopom.Project
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	e := decoder.Decode(&p)
	//e := xml.Unmarshal(data, &p)
	if e != nil {
		return nil, errors.Wrap(e, "create PomFile failed, parse err.")
	}
	return &PomFile{
		pom: p,
	}, nil
}

var pomInlineParameterPattern = regexp.MustCompile("\\$\\{([^{}]+)\\}")

func (p *PomFile) inheritancePath() []*PomFile {
	visited := map[Coordinate]struct{}{}
	visited[p.Coordinate()] = struct{}{}
	var rs []*PomFile
	rs = append(rs, p)
	curr := p.parentPom
	for curr != nil {
		if _, ok := visited[curr.Coordinate()]; ok {
			var s []string
			for _, it := range rs {
				s = append(s, it.Coordinate().String())
			}
			logger.Warn.Println("Circular pom inheritance detected:", strings.Join(s, "->"))
			break
		}
		rs = append(rs, curr)
		curr = curr.parentPom
	}
	return rs
}

func (p *PomFile) Dependencies() []Coordinate {
	// todo: exclusion
	poms := p.inheritancePath()
	type id struct {
		groupId    string
		artifactId string
	}
	dependencyManagement := map[id]string{}
	dependencies := map[Coordinate]struct{}{}
	for _, pom := range poms {
		for _, it := range pom.pom.DependencyManagement.Dependencies {
			i := id{it.GroupID, it.ArtifactID}
			if _, ok := dependencyManagement[i]; ok {
				continue
			}
			v := pom.property(it.Version)
			if v == "" {
				continue
			}
			dependencyManagement[i] = v
		}
	}
	for _, pom := range poms {
		for _, it := range pom.pom.Dependencies {
			if it.Optional == "true" || (it.Type != "" && it.Type != "jar") {
				continue
			}
			if !(it.Scope == "" || it.Scope == "runtime" || it.Scope == "compile") {
				continue
			}
			c := Coordinate{
				GroupId:    it.GroupID,
				ArtifactId: it.ArtifactID,
				Version:    pom.property(it.Version),
			}
			if c.Version == "" {
				c.Version = dependencyManagement[id{it.GroupID, it.ArtifactID}]
			}
			dependencies[c] = struct{}{}
		}
	}
	var rs []Coordinate
	for it := range dependencies {
		rs = append(rs, it)
	}
	return rs
}

func (p *PomFile) getProperty(name string) string {
	ele := p.inheritancePath()
	for _, it := range ele {
		if it.pom.Properties.Entries == nil {
			continue
		}
		value := it.pom.Properties.Entries[name]
		if value != "" {
			return value
		}
	}
	return ""
}

func (p *PomFile) property(s string) string {
	rawStr := pomInlineParameterPattern.Split(s, -1)
	rs := []string{rawStr[0]}
	matches := pomInlineParameterPattern.FindAllStringSubmatch(s, -1)
	for index, match := range matches {
		rs = append(rs, p.getProperty(match[1]))
		rs = append(rs, rawStr[index+1])
	}
	return strings.Join(rs, "")
}

func (p *PomFile) parentCoordinate() *Coordinate {
	if p.pom.Parent.ArtifactID == "" {
		return nil
	}
	return &Coordinate{
		GroupId:    p.property(p.pom.Parent.GroupID),
		ArtifactId: p.property(p.pom.Parent.ArtifactID),
		Version:    p.property(p.pom.Parent.Version),
	}
}

func (p *PomFile) Coordinate() Coordinate {
	c := Coordinate{
		GroupId:    p.property(p.pom.GroupID),
		ArtifactId: p.property(p.pom.ArtifactID),
		Version:    p.property(p.pom.Version),
	}
	if c.GroupId == "" {
		c.GroupId = p.property(p.pom.Parent.GroupID)
	}
	if c.Version == "" {
		c.Version = p.property(p.pom.Parent.Version)
	}
	return c
}

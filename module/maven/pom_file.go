package maven

import (
	"github.com/vifraa/gopom"
	"strings"
)

type PomFile struct {
	pom                  gopom.Project
	parentPom            *PomFile
	path                 string
	propertyMap          map[string]string
	dependencyManagement map[Coordinate]string
	dependencies         []PomDependencyItem
	coordinate           Coordinate
}

type PomDependencyItem struct {
	Coordinate
	Scope string
}

func NewPomFile() *PomFile {
	return &PomFile{
		parentPom:            nil,
		path:                 "",
		propertyMap:          map[string]string{},
		dependencyManagement: map[Coordinate]string{},
	}
}

func (p *PomFile) property(s string) string {
	rawStr := pomInlineParameterPattern.Split(s, -1)
	rs := []string{rawStr[0]}
	matches := pomInlineParameterPattern.FindAllStringSubmatch(s, -1)
	for index, match := range matches {
		if s, ok := p.propertyMap[match[1]]; ok {
			rs = append(rs, s)
		} else {
			rs = append(rs, match[0])
		}
		rs = append(rs, rawStr[index+1])
	}
	return strings.Join(rs, "")
}

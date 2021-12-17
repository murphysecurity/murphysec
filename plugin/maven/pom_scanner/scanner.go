package pom_scanner

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util"
	"os"
	"path/filepath"
)

type ResolvedDependency struct {
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Dependencies []ResolvedDependency `json:"dependencies"`
}

func parsePom(path string) (*xmlquery.Node, error) {
	if !util.IsFile(path) {
		return nil, errors.New("pom file is not a file")
	}
	file, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer func() { _ = file.Close() }()
	return xmlquery.Parse(file)
}

type PomFile struct {
	artifactId              string
	groupId                 string
	version                 string
	path                    string
	dir                     string
	properties              map[string]string
	childPom                []*PomFile
	rawDependencyManagement []RawDependency
	rawDependencies         []RawDependency
	// "{groupId}:{artifactId}:{version}"
	parentModuleId string
	parentPomRef   *PomFile
}

func (f *PomFile) getModuleId() string {
	return fmt.Sprintf("%s:%s:%s", f.groupId, f.artifactId, f.version)
}

func (f *PomFile) getProperty(name string) string {
	if v, ok := f.properties[name]; ok {
		return v
	}
	if f.parentPomRef != nil {
		return f.parentPomRef.getProperty(name)
	}
	return ""
}

type RawDependency struct {
	GroupId    string
	ArtifactId string
	Version    string
}

func (f *PomFile) New() *PomFile {
	if f.properties == nil {
		f.properties = map[string]string{}
	}
	return f
}

func scanPom(dir string, scannedDirs map[string]bool) (*PomFile, error) {
	if scannedDirs == nil {
		scannedDirs = map[string]bool{}
	}
	if scannedDirs[dir] {
		return nil, errors.New("Circular module reference detected")
	}
	scannedDirs[dir] = true
	pomPath := filepath.Join(dir, "pom.xml")
	pom, err := parsePom(pomPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Read pom file failed: %s", dir))
	}
	p := new(PomFile).New()
	p.dir = dir
	p.path = pomPath

	// read moduleId
	{
		var groupId string
		var artifactId string
		var version string
		if n := xmlquery.FindOne(pom, "/project/groupId"); n != nil {
			groupId = n.InnerText()
		}

		if n := xmlquery.FindOne(pom, "/project/artifactId"); n != nil {
			artifactId = n.InnerText()
		}

		if n := xmlquery.FindOne(pom, "/project/version"); n != nil {
			version = n.InnerText()
		}

		if groupId == "" || artifactId == "" || version == "" {
			return nil, errors.New(fmt.Sprintf("module info incomplete: %s", pomPath))
		}
		p.groupId = groupId
		p.artifactId = artifactId
		p.version = version
	}

	// read parent
	if parent := xmlquery.FindOne(pom, "/project/parent"); parent != nil {
		var groupId string
		var artifactId string
		var version string

		if n := xmlquery.FindOne(parent, "/groupId"); n != nil {
			groupId = n.InnerText()
		}
		if n := xmlquery.FindOne(parent, "/artifactId"); n != nil {
			artifactId = n.InnerText()
		}
		if n := xmlquery.FindOne(parent, "/version"); n != nil {
			version = n.InnerText()
		}
		if groupId == "" || artifactId == "" || version == "" {
			return nil, errors.New(fmt.Sprintf("parent info incomplete: %s", pomPath))
		}
		p.parentModuleId = fmt.Sprintf("%s:%s:%s", groupId, artifactId, version)
	}

	// read properties
	for _, it := range xmlquery.Find(pom, "/project/properties/*") {
		if it != nil && it.Type == xmlquery.ElementNode {
			p.properties[it.Data] = it.InnerText()
		}
	}

	// read dependencyManagement
	for _, it := range xmlquery.Find(pom, "/project/dependencyManagement/dependencies/dependency") {
		if it != nil && it.Type == xmlquery.ElementNode {
			var artifactId string
			var groupId string
			var version string
			if n := xmlquery.FindOne(it, "/artifactId"); n != nil {
				artifactId = n.InnerText()
			}
			if n := xmlquery.FindOne(it, "/groupId"); n != nil {
				groupId = n.InnerText()
			}
			if n := xmlquery.FindOne(it, "/version"); n != nil {
				version = n.InnerText()
			}
			if artifactId == "" || groupId == "" || version == "" {
				continue
			}
			p.rawDependencyManagement = append(p.rawDependencyManagement, RawDependency{
				GroupId:    groupId,
				ArtifactId: artifactId,
				Version:    version,
			})
		}
	}

	// read dependencies
	for _, it := range xmlquery.Find(pom, "/project/dependencies/dependency") {
		if it != nil && it.Type == xmlquery.ElementNode {
			var artifactId string
			var groupId string
			var version string
			if n := xmlquery.FindOne(it, "/groupId"); n != nil {
				groupId = n.InnerText()
			}
			if n := xmlquery.FindOne(it, "/artifactId"); n != nil {
				artifactId = n.InnerText()
			}
			if n := xmlquery.FindOne(it, "version"); n != nil {
				version = n.InnerText()
			}
			if artifactId == "" || groupId == "" {
				continue
			}
			p.rawDependencies = append(p.rawDependencies, RawDependency{
				GroupId:    groupId,
				ArtifactId: artifactId,
				Version:    version,
			})
		}
	}

	// read module
	var modules []string
	for _, it := range xmlquery.Find(pom, "/project/modules/module") {
		if it != nil || it.Type == xmlquery.ElementNode {
			modules = append(modules, it.InnerText())
		}
	}

	// recursive parse pom in modules
	for _, it := range modules {
		cpom, err := scanPom(filepath.Join(dir, it), scannedDirs)
		if err != nil {
			return nil, err
		}
		p.childPom = append(p.childPom, cpom)
	}

	return p, nil
}

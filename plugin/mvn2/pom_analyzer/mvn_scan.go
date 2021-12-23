package pom_analyzer

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func PomAnalyze(dir string) (*PomFile, map[string]*PomFile, error) {
	rootPom, e := openPom(dir)
	if e != nil {
		return nil, nil, e
	}
	e = _findModules(rootPom, []string{})
	if e != nil {
		return nil, nil, e
	}
	poms, e := _flattenModulePom(rootPom, []*PomFile{})
	if e != nil {
		return nil, nil, e
	}
	pomMap := map[string]*PomFile{}
	for _, it := range poms {
		pomMap[it.Id()] = it
	}
	e = resolveInheritance(rootPom)
	if e != nil {
		return nil, nil, e
	}
	for _, it := range pomMap {
		parseDependencyManagement(it)
	}
	for _, it := range pomMap {
		parseDependency(it, pomMap)
	}
	return rootPom, pomMap, nil
}

type Dependency struct {
	GroupId      string
	ArtifactId   string
	Version      string
	PomFile      *PomFile
	Dependencies []*Dependency
}

func (d *Dependency) Id() string {
	return fmt.Sprintf("%s:%s:%s", d.GroupId, d.ArtifactId, d.Version)
}

var inlineParameterPattern = func() *regexp.Regexp {
	return regexp.MustCompile("\\$\\{([^{}]+)\\}")
}()

func parseDependency(p *PomFile, localPoms map[string]*PomFile) {
	d := p.XmlDoc
	var deps []*Dependency
	for _, it := range xmlquery.Find(d, "/project/dependencies/dependency") {
		var groupId string
		var artifactId string
		var version string
		if n := xmlquery.FindOne(it, "/groupId"); n != nil {
			groupId = resolvePropertiesValue(p, n.InnerText())
		}
		if n := xmlquery.FindOne(it, "/artifactId"); n != nil {
			artifactId = resolvePropertiesValue(p, n.InnerText())
		}
		if n := xmlquery.FindOne(it, "/version"); n != nil {
			version = resolvePropertiesValue(p, n.InnerText())
		}
		if version == "" {
			version = p.getDependencyVersion(groupId, artifactId)
		}
		deps = append(deps, &Dependency{
			GroupId:    groupId,
			ArtifactId: artifactId,
			Version:    version,
			PomFile:    localPoms[groupId+":"+artifactId+":"+version],
		})
	}
	p.Dependencies = deps
}

func parseDependencyManagement(p *PomFile) {
	doc := p.XmlDoc
	m := map[string]string{}
	for _, it := range xmlquery.Find(doc, "/project/dependencyManagement/dependencies/dependency") {
		var groupId string
		var artifactId string
		var version string
		if n := xmlquery.FindOne(it, "/groupId"); n != nil {
			groupId = resolvePropertiesValue(p, n.InnerText())
		}
		if n := xmlquery.FindOne(it, "/artifactId"); n != nil {
			artifactId = resolvePropertiesValue(p, n.InnerText())
		}
		if n := xmlquery.FindOne(it, "/version"); n != nil {
			version = resolvePropertiesValue(p, n.InnerText())
		}
		m[groupId+":"+artifactId] = version
	}
	p.DependencyManagement = m
}

func resolvePropertiesValue(p *PomFile, s string) string {
	rawStr := inlineParameterPattern.Split(s, -1)
	rs := []string{rawStr[0]}
	matches := inlineParameterPattern.FindAllStringSubmatch(s, -1)
	for index, match := range matches {
		key := match[1]
		if value, ok := p.Properties[key]; !ok {
			value = match[0]
		} else {
			rs = append(rs, value)
		}
		rs = append(rs, rawStr[index+1])
	}
	return strings.Join(rs, "")
}

func resolveInheritance(root *PomFile) error {
	r, e := _flattenModulePom(root, []*PomFile{})
	if e != nil {
		return e
	}
	rmap := map[string]*PomFile{}
	for _, it := range r {
		rmap[it.Id()] = it
	}
	for _, it := range r {
		if it.ParentPomDescriptor == nil {
			continue
		}
		pNode := rmap[it.ParentPomDescriptor.Id()]
		if pNode == nil {
			continue
		}
		it.ParentPom = pNode
	}
	for _, it := range r {
		e := _checkInheritanceCircular(it, []*PomFile{})
		if e != nil {
			return e
		}
	}
	return nil
}

func _checkInheritanceCircular(p *PomFile, accessPoms []*PomFile) error {
	if p == nil {
		return nil
	}
	for _, it := range accessPoms {
		if it == p {
			return errors.New("circular inheritance")
		}
	}
	if p.ParentPom != nil {
		return _checkInheritanceCircular(p.ParentPom, append(accessPoms, p))
	}
	return nil
}

func _flattenModulePom(root *PomFile, accessedPom []*PomFile) ([]*PomFile, error) {
	for _, it := range accessedPom {
		if it == root {
			return nil, errors.New("circular modules")
		}
	}
	rs := []*PomFile{root}
	for _, it := range root.Modules {
		r, e := _flattenModulePom(it, append(accessedPom, root))
		if e != nil {
			return nil, e
		}
		rs = append(rs, r...)
	}
	return rs, nil
}
func openPom(dir string) (*PomFile, error) {
	pomFilePath := filepath.Join(dir, "pom.xml")
	d, e := parsePom(pomFilePath)
	if e != nil {
		return nil, e

	}
	r := PomFile{Dir: dir, Path: pomFilePath, XmlDoc: d, Properties: map[string]string{}}

	if !checkPomVersion(d) {
		return nil, errors.New("only support modelVersion: 4.0.0")
	}

	for _, it := range xmlquery.Find(d, "/project/properties/*") {
		r.Properties[it.Data] = it.InnerText()
	}

	if n := xmlquery.FindOne(d, "/project/groupId"); n != nil {
		r.GroupId = n.InnerText()
	}
	if n := xmlquery.FindOne(d, "/project/artifactId"); n != nil {
		r.ArtifactId = n.InnerText()
	}
	if n := xmlquery.FindOne(d, "/project/version"); n != nil {
		r.Version = n.InnerText()
	}

	if parent := xmlquery.FindOne(d, "/project/parent"); parent != nil {
		var artifactId string
		var groupId string
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
		r.ParentPomDescriptor = &ParentPomDescriptor{
			GroupId:    groupId,
			ArtifactId: artifactId,
			Version:    version,
		}
	}

	return &r, nil
}

func _findModules(p *PomFile, accessedFile []string) error {
	var moduleDirs []string
	for _, it := range xmlquery.Find(p.XmlDoc, "/project/modules/module") {
		if it == nil {
			continue
		}
		d := filepath.Join(p.Dir, it.InnerText())
		if !util.IsDir(d) {
			return errors.New(fmt.Sprintf("analyze module failed: %s", it.InnerText()))
		}
		if util.InStringSlice(accessedFile, d) {
			return errors.New("circular module references detected")
		}
		moduleDirs = append(moduleDirs, d)
	}
	for _, it := range moduleDirs {
		mp, e := openPom(it)
		if e != nil {
			return e
		}
		e = _findModules(mp, append(accessedFile, it))
		if e != nil {
			return e
		}
		p.Modules = append(p.Modules, mp)
	}
	return nil
}

func (r *PomFile) Id() string {
	return fmt.Sprintf("%s:%s:%s", r.GroupId, r.ArtifactId, r.Version)
}

type PomFile struct {
	Path                 string
	Dir                  string
	Modules              []*PomFile
	XmlDoc               *xmlquery.Node `json:"-"`
	GroupId              string
	ArtifactId           string
	Version              string
	Properties           map[string]string
	ParentPom            *PomFile `json:"-"`
	ParentPomDescriptor  *ParentPomDescriptor
	DependencyManagement map[string]string
	Dependencies         []*Dependency
}

func (r *PomFile) getDependencyVersion(groupId string, artifactId string) string {
	k := groupId + ":" + artifactId
	v := r.DependencyManagement[k]
	if v == "" && r.ParentPom != nil {
		return r.ParentPom.getDependencyVersion(groupId, artifactId)
	}
	return v
}

func (r *ParentPomDescriptor) Id() string {
	return fmt.Sprintf("%s:%s:%s", r.GroupId, r.ArtifactId, r.Version)
}

type ParentPomDescriptor struct {
	GroupId    string
	ArtifactId string
	Version    string
}

func checkPomVersion(doc *xmlquery.Node) bool {
	mv := xmlquery.FindOne(doc, "/project/modelVersion")
	return mv != nil && mv.InnerText() == "4.0.0"
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

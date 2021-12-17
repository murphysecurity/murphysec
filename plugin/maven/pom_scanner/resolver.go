package pom_scanner

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func resolve(root *PomFile) (*ResolvedDependency, error) {
	resolveInheritance(root)
	fmt.Println("========")
	if e := detectInheritanceCircular(root, map[*PomFile]bool{}); e != nil {
		return nil, e
	}
	modules := map[string]*PomFile{}
	// flatten modules
	q := []*PomFile{root}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for _, it := range cur.childPom {
			q = append(q, it)
		}
		modules[cur.getModuleId()] = cur
	}
	rootDependency := ResolvedDependency{
		Name:         root.getModuleId(),
		Version:      "",
		Dependencies: []ResolvedDependency{},
	}
	for _, it := range modules {
		if r, e := resolvePom(it, map[string]bool{}, modules); e != nil {
			return nil, e
		} else {
			rootDependency.Dependencies = append(rootDependency.Dependencies, ResolvedDependency{
				Name:         it.groupId + ":" + it.artifactId,
				Version:      it.version,
				Dependencies: r,
			})
		}
	}
	return &rootDependency, nil
}

func resolvePom(p *PomFile, resolvedDep map[string]bool, pomFiles map[string]*PomFile) ([]ResolvedDependency, error) {
	moduleId := p.getModuleId()
	if resolvedDep[moduleId] {
		return nil, errors.New("circular dependency detected")
	}
	resolvedDep[moduleId] = true
	defer delete(resolvedDep, moduleId)
	dm := collectDependencyManagementVersionMapping(p)
	rs := make([]ResolvedDependency, 0)
	for _, it := range p.rawDependencies {
		name := fmt.Sprintf("%s:%s", resolveStrProperty(it.GroupId, p), resolveStrProperty(it.ArtifactId, p))
		version := resolveStrProperty(it.Version, p)
		if version == "" {
			version = dm[name]
		}
		rd := ResolvedDependency{
			Name:         name,
			Version:      version,
			Dependencies: []ResolvedDependency{},
		}
		if pd := pomFiles[name+":"+version]; pd != nil {
			if rs, e := resolvePom(pd, resolvedDep, pomFiles); e != nil {
				return nil, e
			} else {
				rd.Dependencies = append(rd.Dependencies, rs...)
			}
		}
		rs = append(rs, rd)
	}
	return rs, nil
}

func collectDependencyManagementVersionMapping(p *PomFile) map[string]string {
	var rs map[string]string
	if p.parentPomRef != nil {
		// collect inherited dependencyManagement
		rs = collectDependencyManagementVersionMapping(p.parentPomRef)
	} else {
		rs = map[string]string{}
	}
	for _, it := range p.rawDependencyManagement {
		depName := fmt.Sprintf("%s:%s", resolveStrProperty(it.GroupId, p), resolveStrProperty(it.ArtifactId, p))
		version := resolveStrProperty(it.Version, p)
		rs[depName] = version
	}
	return rs
}

func resolveInheritance(rootPom *PomFile) {
	m := map[string]*PomFile{}

	q := []*PomFile{rootPom}
	curIndex := 0
	for curIndex < len(q) {
		for i := range q[curIndex].childPom {
			q = append(q, q[curIndex].childPom[i])
		}
		curIndex++
	}

	for _, it := range q {
		m[it.getModuleId()] = it
	}
	for _, it := range q {
		if it.parentModuleId != "" {
			it.parentPomRef = m[it.parentModuleId]
		}
	}
}

func detectInheritanceCircular(pom *PomFile, visitedSet map[*PomFile]bool) error {
	if visitedSet[pom] {
		return errors.New("circular module inheritance detected")
	}
	visitedSet[pom] = true
	defer delete(visitedSet, pom)
	if pom.parentPomRef == nil {
		return nil
	}
	return detectInheritanceCircular(pom.parentPomRef, visitedSet)
}

var propertyPattern = func() *regexp.Regexp { return regexp.MustCompile("\\$\\{([^{}]+)\\}") }()

func resolveStrProperty(str string, pom *PomFile) string {
	p := propertyPattern.Split(str, -1)
	rs := []string{p[0]}
	ptr := 1
	for _, ints := range propertyPattern.FindAllSubmatchIndex([]byte(str), -1) {
		key := str[ints[2]:ints[3]]
		//raw:=str[ints[0]:ints[1]]
		v := pom.getProperty(key)
		rs = append(rs, v, p[ptr])
		ptr++
	}
	return strings.Join(rs, "")
}

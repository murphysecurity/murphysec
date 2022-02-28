package maven

import (
	"github.com/vifraa/gopom"
	"murphysec-cli-simple/logger"
	"path/filepath"
	"strings"
)

func InspectModule(dir string) map[string]*PomBuilder {
	rs := map[string]*PomBuilder{}
	_inspectModule(dir, dir, nil, rs)
	return rs
}

func _inspectModule(dir string, basedir string, visited map[string]struct{}, rs map[string]*PomBuilder) {
	if s, e := filepath.Rel(basedir, dir); e != nil || strings.HasPrefix(s, "../") {
		return
	}
	{
		if visited == nil {
			visited = map[string]struct{}{}
		}
		if _, ok := visited[dir]; ok {
			return
		}
		visited[dir] = struct{}{}
		defer delete(visited, dir)
	}
	p, e := gopom.Parse(filepath.Join(dir, "pom.xml"))
	if e != nil {
		logger.Warn.Println("parse local pom failed", e.Error(), dir)
		return
	}
	builder := NewPomBuilder(p)
	builder.Path = dir
	rs[dir] = builder

	for _, module := range p.Modules {
		_inspectModule(filepath.Join(dir, module), basedir, visited, rs)
	}
}

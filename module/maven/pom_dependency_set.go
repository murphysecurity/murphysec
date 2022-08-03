package maven

import (
	"fmt"
	"github.com/vifraa/gopom"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"sort"
	"strings"
)

type pomDependencySet struct {
	m *orderedmap.OrderedMap[string, *gopom.Dependency]
}

func newPomDependencySet() *pomDependencySet {
	return &pomDependencySet{m: orderedmap.New[string, *gopom.Dependency]()}
}

func (p *pomDependencySet) String() string {
	var rs = make([]string, 0, p.m.Len())
	for pair := p.m.Oldest(); pair != nil; pair = pair.Next() {
		rs = append(rs, fmt.Sprintf("%v -> %v", pair.Key, pair.Value))
	}
	return strings.Join(rs, "\n")
}

func (p *pomDependencySet) listAll() []gopom.Dependency {
	var rs []gopom.Dependency
	for pair := p.m.Oldest(); pair != nil; pair = pair.Next() {
		rs = append(rs, *pair.Value)
	}
	return rs
}

func (p *pomDependencySet) mergeProperty(property *properties) {
	for pair := p.m.Oldest(); pair != nil; pair = pair.Next() {
		pair.Value.ArtifactID = property.Resolve(pair.Value.ArtifactID)
		pair.Value.GroupID = property.Resolve(pair.Value.GroupID)
		pair.Value.Version = property.Resolve(pair.Value.Version)
	}
}

func (p *pomDependencySet) mergeAll(b []gopom.Dependency, override bool, ignoreIfNotExists bool) {
	for _, dependency := range b {
		p.mergeItem(dependency, override, ignoreIfNotExists)
	}
}

func (p *pomDependencySet) mergeItem(b gopom.Dependency, override bool, ignoreIfNotExists bool) {
	k := b.GroupID + ":" + b.ArtifactID
	var a, ok = p.m.Get(k)
	if !ok && ignoreIfNotExists {
		return
	}
	if !ok {
		p.m.Set(k, &b)
		return
	}
	if (b.Version != "" && override) || a.Version == "" {
		a.Version = b.Version
	}
	if (b.Scope != "" && override) || a.Scope == "" {
		a.Scope = b.Scope
	}
	a.Exclusions = _mergeExclusions(a.Exclusions, b.Exclusions)
}

func _mergeExclusions(old []gopom.Exclusion, new []gopom.Exclusion) []gopom.Exclusion {
	if len(old) == 0 {
		return new
	}
	if len(new) == 0 {
		return old
	}
	var m = map[gopom.Exclusion]struct{}{}
	for _, it := range old {
		m[it] = struct{}{}
	}
	for _, it := range new {
		m[it] = struct{}{}
	}
	var rs []gopom.Exclusion
	for exclusion := range m {
		rs = append(rs, exclusion)
	}
	sort.Slice(rs, func(i, j int) bool {
		if rs[i].ArtifactID == rs[j].ArtifactID {
			return rs[i].GroupID < rs[j].GroupID
		}
		return rs[i].ArtifactID <= rs[j].ArtifactID
	})
	return rs
}

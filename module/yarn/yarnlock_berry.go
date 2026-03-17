package yarn

import (
	neturl "net/url"
	"strings"

	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/module/pkgjs"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
)

type berryLockData struct {
	Lockfile yarnlock.LockFile
	Index    map[string]string
}

func parseYarnLockYaml(text string) (yarnlock.LockFile, error) {
	data, e := parseYarnLockYamlWithIndex(text)
	if e != nil {
		return nil, e
	}
	return data.Lockfile, nil
}

func parseYarnLockYamlWithIndex(text string) (*berryLockData, error) {
	type Element struct {
		Version      string            `yaml:"version"`
		Dependencies map[string]string `yaml:"dependencies"`
		OptionalDeps map[string]string `yaml:"optionalDependencies"`
	}
	var pkgs map[string]Element
	e := shared.ParseYaml([]byte(text), &pkgs)
	if e != nil {
		return nil, e
	}

	lockfile := make(yarnlock.LockFile)
	index := map[string]string{}
	for key, value := range pkgs {
		for keyEl := range strings.SplitSeq(key, ", ") {
			keyEl = strings.TrimSpace(keyEl)
			if keyEl == "" {
				continue
			}
			entry := yarnlock.LockFileEntry{
				Version:              value.Version,
				Dependencies:         value.Dependencies,
				OptionalDependencies: value.OptionalDeps,
			}
			lockfile[keyEl] = entry
			if pkgName, pkgVer := parsePkgName(keyEl); pkgName != "" && pkgVer != "" {
				for _, normalized := range normalizeRangeCandidates(pkgVer) {
					k := berryIndexKey(pkgName, normalized)
					if _, ok := index[k]; !ok {
						index[k] = keyEl
					}
				}
			}
		}
	}
	return &berryLockData{
		Lockfile: lockfile,
		Index:    index,
	}, nil
}

var yarnVersionProtocols = map[string]struct{}{
	"npm":       {},
	"workspace": {},
	"patch":     {},
	"file":      {},
	"link":      {},
	"portal":    {},
	"exec":      {},
	"git":       {},
	"github":    {},
	"http":      {},
	"https":     {},
	"git+https": {},
	"git+ssh":   {},
}

func normalizeRangeCandidates(pkgVer string) []string {
	seen := map[string]struct{}{}
	queue := []string{pkgVer}
	candidates := make([]string, 0, 8)
	add := func(v string) {
		v = strings.TrimSpace(v)
		v = strings.Trim(v, "\"")
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		candidates = append(candidates, v)
		queue = append(queue, v)
	}
	add(pkgVer)
	for i := 0; i < len(queue); i++ {
		v := queue[i]
		if decoded, err := neturl.QueryUnescape(v); err == nil && decoded != v {
			add(decoded)
		}
		if idx := strings.Index(v, ":"); idx > 0 {
			if _, ok := yarnVersionProtocols[v[:idx]]; ok {
				add(v[idx+1:])
			}
		}
		if idx := strings.Index(v, "@npm:"); idx >= 0 {
			raw := v[idx+len("@npm:"):]
			add(raw)
			if hash := strings.Index(raw, "#"); hash > 0 {
				add(raw[:hash])
			}
			if amp := strings.Index(raw, "&"); amp > 0 {
				add(raw[:amp])
			}
		}
		if strings.HasPrefix(v, "patch:") {
			inner := strings.TrimPrefix(v, "patch:")
			if hash := strings.Index(inner, "#"); hash > 0 {
				inner = inner[:hash]
			}
			if p, rng := parsePkgName(inner); p != "" && rng != "" {
				add(rng)
			}
		}
	}
	return candidates
}

func berryIndexKey(name, rng string) string {
	return name + "@" + rng
}

func resolveBerryEntry(data *berryLockData, name, rng string) (string, yarnlock.LockFileEntry, bool) {
	if data == nil {
		return "", yarnlock.LockFileEntry{}, false
	}
	for _, normalized := range normalizeRangeCandidates(rng) {
		if selector, ok := data.Index[berryIndexKey(name, normalized)]; ok {
			if entry, ok := data.Lockfile[selector]; ok {
				return selector, entry, true
			}
		}
	}
	return "", yarnlock.LockFileEntry{}, false
}

func buildDepTreeBerry(data *berryLockData, pkg *pkgjs.Pkg) []Dep {
	type id struct {
		name    string
		version string
	}
	var rs []Dep
	repeatedElement := map[id]struct{}{}
	for n, v := range pkg.Dependencies {
		node := _buildDepTreeBerry(data, n, v, map[string]struct{}{}, 32)
		if node == nil {
			continue
		}
		key := id{node.Name, node.Version}
		if _, ok := repeatedElement[key]; ok {
			continue
		}
		repeatedElement[key] = struct{}{}
		rs = append(rs, *node)
	}
	for n, v := range pkg.DevDependencies {
		node := _buildDepTreeBerry(data, n, v, map[string]struct{}{}, 32)
		if node == nil {
			continue
		}
		key := id{node.Name, node.Version}
		if _, ok := repeatedElement[key]; ok {
			continue
		}
		repeatedElement[key] = struct{}{}
		rs = append(rs, *node)
	}
	return rs
}

func _buildDepTreeBerry(data *berryLockData, pkgName string, pkgRange string, visitedKey map[string]struct{}, depth int) *Dep {
	if depth < 0 {
		return nil
	}
	selector, info, ok := resolveBerryEntry(data, pkgName, pkgRange)
	if !ok {
		return nil
	}
	if _, ok := visitedKey[selector]; ok {
		return nil
	}
	visitedKey[selector] = struct{}{}
	defer delete(visitedKey, selector)

	node := &Dep{
		Name:    pkgName,
		Version: info.Version,
	}
	type id struct {
		name    string
		version string
	}
	repeatedElement := map[id]struct{}{}
	for childComp, childVer := range info.Dependencies {
		c := _buildDepTreeBerry(data, childComp, childVer, visitedKey, depth-1)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Children = append(node.Children, *c)
	}
	for childComp, childVer := range info.OptionalDependencies {
		c := _buildDepTreeBerry(data, childComp, childVer, visitedKey, depth-1)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Children = append(node.Children, *c)
	}
	return node
}

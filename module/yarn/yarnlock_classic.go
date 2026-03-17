package yarn

import (
	"regexp"

	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/module/pkgjs"
)

func buildDepTree(lkFile yarnlock.LockFile, pkg *pkgjs.Pkg) []Dep {
	type id struct {
		name    string
		version string
	}
	var rs []Dep
	repeatedElement := map[id]struct{}{}
	for n, v := range pkg.Dependencies {
		node := _buildDepTree(lkFile, n+"@"+v, map[string]struct{}{}, 5)
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
		node := _buildDepTree(lkFile, n+"@"+v, map[string]struct{}{}, 5)
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

func _buildDepTree(lkFile yarnlock.LockFile, element string, visitedKey map[string]struct{}, depth int) *Dep {
	if depth < 0 {
		return nil
	}
	if _, ok := visitedKey[element]; ok {
		return nil
	}
	visitedKey[element] = struct{}{}
	defer delete(visitedKey, element)

	info, ok := lkFile[element]
	if !ok {
		return nil
	}
	pkgName, pkgVer := parsePkgName(element)
	if pkgName == "" || pkgVer == "" {
		return nil
	}
	node := &Dep{
		Name:    pkgName,
		Version: info.Version,
	}
	type id struct {
		name    string
		version string
	}
	repeatedElement := map[id]struct{}{}
	for childComp, childVer := range lkFile[element].Dependencies {
		childKey := childComp + "@" + childVer
		c := _buildDepTree(lkFile, childKey, visitedKey, depth-1)
		if c == nil {
			continue
		}
		if _, ok := repeatedElement[id{c.Name, c.Version}]; ok {
			continue
		}
		repeatedElement[id{c.Name, c.Version}] = struct{}{}
		node.Children = append(node.Children, *c)
	}
	for childComp, childVer := range lkFile[element].OptionalDependencies {
		childKey := childComp + "@" + childVer
		c := _buildDepTree(lkFile, childKey, visitedKey, depth-1)
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

var versionNpmPattern = regexp.MustCompile(`^npm:(.+?)@(.+)`)
var parsePkgNamePattern = regexp.MustCompile("(@?[^@]+)@(.+)")

func parsePkgName(input string) (pkgName string, pkgVersion string) {
	m := parsePkgNamePattern.FindStringSubmatch(input)
	if m == nil {
		return "", ""
	}
	if m := versionNpmPattern.FindStringSubmatch(m[2]); m != nil {
		return m[1], m[2]
	}
	return m[1], m[2]
}

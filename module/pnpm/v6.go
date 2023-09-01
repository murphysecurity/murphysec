package pnpm

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"strings"
)

type v6Lockfile struct {
	LockfileVersion string               `json:"lockfile_version" yaml:"lockfileVersion"`
	Dependencies    map[string]v6Package `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]v6Package `json:"dev_dependencies" yaml:"devDependencies"`
	Packages        map[string]v6Package `json:"packages" yaml:"packages"`
	pathMapping     map[[2]string]string
}

func parseV6Lockfile(data []byte, strict bool) (*v6Lockfile, error) {
	var e error
	var r v6Lockfile
	e = yaml.Unmarshal(data, &r)
	if e != nil {
		return nil, fmt.Errorf("parseV6Lockfile: %w", e)
	}
	e = r.postprocess(strict)
	if e != nil {
		return nil, fmt.Errorf("parseV6Lockfile: %w", e)
	}
	return &r, nil
}

func (v *v6Lockfile) buildDependencyTree(strict bool) ([]model.DependencyItem, error) {
	var deps [][2]string
	for name, pkg := range v.Dependencies {
		if name == "" {
			if !strict {
				continue
			}
			return nil, fmt.Errorf("buildDependencyTree: dependencies name is empty")
		}
		var version = trimVersionSuffix(pkg.Version)
		if version == "" {
			if !strict {
				continue
			}
			return nil, fmt.Errorf("buildDependencyTree: dependencies %s version is empty", name)
		}
		deps = append(deps, [2]string{name, version})
	}
	for name, pkg := range v.DevDependencies {
		if name == "" {
			if !strict {
				continue
			}
			return nil, fmt.Errorf("buildDependencyTree: devDependencies name is empty")
		}
		var version = trimVersionSuffix(pkg.Version)
		if version == "" {
			if !strict {
				continue
			}
			return nil, fmt.Errorf("buildDependencyTree: devDependencies %s version is empty", name)
		}
		deps = append(deps, [2]string{name, version})
	}
	deps = lo.Uniq(deps)
	var rdeps []model.DependencyItem
	for _, dep := range deps {
		node, e := v.buildDependencyTree0(dep[0], dep[1], strict, nil)
		if e != nil {
			if !strict {
				continue
			}
			return nil, e
		}
		rdeps = append(rdeps, *node)
	}
	return rdeps, nil
}

func (v *v6Lockfile) buildDependencyTree0(name, version string, strict bool, visited *visitedSeg) (*model.DependencyItem, error) {
	key := [2]string{name, version}
	if visited != nil && visited.Contain(name, version) {
		return nil, &recursiveVisitError{visited}
	}
	nextVisitObj := visitedSegOf(name, version, visited)
	path, ok := v.pathMapping[key]
	if !ok {
		return nil, &unknownPkgError{name, version}
	}
	pkg, ok := v.Packages[path]
	if !ok {
		panic("v.Packages[lockfile] !ok")
	}
	r := &model.DependencyItem{
		Component: model.Component{
			CompName:    name,
			CompVersion: version,
			EcoRepo:     EcoRepo,
		},
		Dependencies: nil,
	}
	r.IsOnline.SetOnline(!pkg.Optional && !pkg.Dev)
	for childName, childVersion := range pkg.Dependencies {
		child, e := v.buildDependencyTree0(childName, trimVersionSuffix(childVersion), strict, nextVisitObj)
		if e != nil {
			if !strict {
				continue
			}
			return nil, e
		}
		r.Dependencies = append(r.Dependencies, *child)
	}
	return r, nil
}

func (v *v6Lockfile) postprocess(strict bool) error {
	v.pathMapping = map[[2]string]string{}
	for path, pkg := range v.Packages {
		var e error
		var name = pkg.Name
		if name == "" {
			name, e = getNameFromPath(path)
			if e != nil {
				if !strict {
					continue
				}
				return postprocessCreateError(path, e)
			}
			if name == "" {
				if !strict {
					continue
				}
				return postprocessCreateError(path, fmt.Errorf("name is empty"))
			}
		}
		var version = trimVersionSuffix(pkg.Version)
		if version == "" {
			version, e = getVersionFromPath(path)
			if e != nil {
				if !strict {
					continue
				}
				return postprocessCreateError(path, e)
			}
			if version == "" {
				if !strict {
					continue
				}
				return postprocessCreateError(path, fmt.Errorf("version is empty"))
			}
		}
		v.pathMapping[[2]string{name, version}] = path
	}
	return nil
}

func postprocessCreateError(path string, e error) error {
	return fmt.Errorf("postprocess lockfile \"%s\": %w", path, e)
}

type v6Package struct {
	Name                 string            `json:"name" yaml:"name"`
	Version              string            `json:"version" yaml:"version"`
	Dev                  bool              `json:"dev" yaml:"dev"`
	Optional             bool              `json:"optional" yaml:"optional"`
	Dependencies         map[string]string `json:"dependencies" yaml:"dependencies"`
	OptionalDependencies map[string]string `json:"optional_dependencies" yaml:"optionalDependencies"`
}

type unknownPkgError struct {
	Name    string
	Version string
}

func (e unknownPkgError) Error() string {
	return fmt.Sprintf("unknown package of name=%s, version=%s", e.Name, e.Version)
}

type recursiveVisitError struct {
	seg *visitedSeg
}

func (e recursiveVisitError) Error() string {
	var r = e.seg
	var arr []*visitedSeg
	for r != nil {
		arr = append(arr, r)
		r = r.Parent
	}
	utils.Reverse(arr)
	return "circular visit: " + strings.Join(fp.Map(func(seg *visitedSeg) string { return seg.Name + "@" + seg.Version })(arr), " -> ")
}

type visitedSeg struct {
	Name    string
	Version string
	Parent  *visitedSeg
}

func visitedSegOf(name, version string, parent *visitedSeg) *visitedSeg {
	if parent != nil && parent.Contain(name, version) {
		panic("parent.Contain(name, version)")
	}
	return &visitedSeg{
		Name:    name,
		Version: version,
		Parent:  parent,
	}
}

func (v *visitedSeg) Contain(name, version string) bool {
	if v.Name == name && v.Version == version {
		return true
	}
	if v.Parent != nil {
		return v.Parent.Contain(name, version)
	}
	return false
}

var unparsableDependencyPath = fmt.Errorf("parsing pnpm dependencyPath, unparsable")

func getNameFromPath(input string) (string, error) {
	var e error
	input, e = trimPathParentSegments(input)
	if e != nil {
		return "", e
	}
	if idx := strings.Index(input, "@"); idx != -1 {
		input = input[:idx]
	}
	return input, nil
}

func getVersionFromPath(input string) (string, error) {
	var e error
	input, e = trimPathParentSegments(input)
	if e != nil {
		return "", e
	}
	if idx := strings.Index(input, "@"); idx != -1 {
		if idx+1 >= len(input) {
			return "", fmt.Errorf("%w: %s", unparsableDependencyPath, input)
		}
		input = input[idx+1:]
	}
	if idx := strings.IndexAny(input, "(+_"); idx != -1 {
		input = input[:idx]
	}
	return input, nil
}

func trimVersionSuffix(input string) string {
	if strings.Contains(input, "/") { // legacy mode
		s, _ := getVersionFromPath(input)
		return s
	}
	if idx := strings.IndexAny(input, "(+_"); idx != -1 {
		return input[:idx]
	}
	return input
}

func trimPathParentSegments(input string) (string, error) {
	if idx := strings.LastIndex(input, "/"); idx != -1 {
		if idx+1 >= len(input) {
			return "", fmt.Errorf("%w: %s", unparsableDependencyPath, input)
		}
		input = input[idx+1:]
	}
	return input, nil
}

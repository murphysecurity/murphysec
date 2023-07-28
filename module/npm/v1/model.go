package v1

import (
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/module/npm/shared"
)

type Lockfile struct {
	root lockRoot
}

func (l *Lockfile) Build(requires [][2]string, strict bool) ([]*shared.Node, error) {
	var r []*shared.Node
	for _, require := range requires {
		name := require[0]
		versionCons := require[1]
		n, e := buildTree(name, versionCons, &l.root.lockPkg, nil, strict)
		if e != nil {
			if !strict {
				continue
			}
			return nil, fmt.Errorf("v1.ParseLockfile: %w", e)
		}
		r = append(r, n)
	}
	return r, nil
}

type lockPkg struct {
	Version      string              `json:"version"`
	Optional     *bool               `json:"optional"`
	Requires     map[string]string   `json:"requires"`
	Dependencies map[string]*lockPkg `json:"dependencies"`
	Dev          *bool               `json:"dev"`
	parent       *lockPkg
}

type lockRoot struct {
	Name string `json:"name"`
	lockPkg
	Requires bool `json:"requires"`
}

func postprocessPkg(pkg *lockPkg, parent *lockPkg) {
	pkg.parent = parent
	for _, p := range pkg.Dependencies {
		postprocessPkg(p, pkg)
	}
}

func buildTree(name string, versionConstraint string, current *lockPkg, visited *shared.Visited, strict bool) (*shared.Node, error) {
	childVisited := visited.CreateSub(name, versionConstraint)
	if childVisited == nil {
		return nil, shared.CreateRevisitError(visited)
	}
	for current != nil {
		childPkg := current.Dependencies[name]
		if childPkg == nil {
			current = current.parent
			continue
		}
		// found
		var node shared.Node
		node.Name = name
		node.Version = childPkg.Version
		if childPkg.Optional != nil {
			node.IsOnline.SetOnline(!*childPkg.Optional)
		}
		if childPkg.Dev != nil {
			node.Dev = *childPkg.Dev
		} else {
			node.Dev = false
		}
		for childName, versionCons := range childPkg.Requires {
			childNode, e := buildTree(childName, versionCons, childPkg, childVisited, strict)
			if e != nil {
				if !strict {
					continue
				}
				return nil, e
			}
			if childNode == nil {
				panic("childNode == nil")
			}
			node.Children = append(node.Children, childNode)
		}
		return &node, nil
	}
	return nil, shared.CreateDependencyNotFoundError(name, versionConstraint)
}

func ParseLockfile(data []byte) (*Lockfile, error) {
	var e error
	var lockRoot lockRoot
	e = json.Unmarshal(data, &lockRoot)
	if e != nil {
		return nil, fmt.Errorf("v1.ParseLockfile: %w", e)
	}
	postprocessPkg(&lockRoot.lockPkg, nil)
	return &Lockfile{lockRoot}, e
}

package conan

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/murphysecurity/murphysec/model"
)

type _ConanGraphInfoJsonFile struct {
	Graph struct {
		Nodes json.RawMessage `json:"nodes"`
	} `json:"graph"`
}

type _ConanGraphNode struct {
	ID       string `json:"id"`
	Ref      string `json:"ref"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Requires []_ConanGraphRequirement `json:"requires"`
	// Conan2 graph json uses "dependencies", which can be a map or array.
	Dependencies json.RawMessage `json:"dependencies"`
}

type _ConanGraphRequirement struct {
	ID      string `json:"id"`
	Ref     string `json:"ref"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (t *_ConanGraphInfoJsonFile) ReadFromFile(path string) error {
	data, e := os.ReadFile(path)
	if e != nil {
		return fmt.Errorf("read conan graph json failed: %w", e)
	}
	if e := json.Unmarshal(data, &t); e != nil {
		return fmt.Errorf("parse conan graph json failed: %w", e)
	}
	return nil
}

func parseConanGraphNodes(raw json.RawMessage) map[string]_ConanGraphNode {
	if len(raw) == 0 {
		return nil
	}
	var nodesMap map[string]_ConanGraphNode
	if e := json.Unmarshal(raw, &nodesMap); e == nil && len(nodesMap) > 0 {
		for k, v := range nodesMap {
			if v.ID == "" {
				v.ID = k
				nodesMap[k] = v
			}
		}
		return nodesMap
	}
	var nodesArray []_ConanGraphNode
	if e := json.Unmarshal(raw, &nodesArray); e == nil && len(nodesArray) > 0 {
		nodesMap = make(map[string]_ConanGraphNode, len(nodesArray))
		for i, v := range nodesArray {
			k := v.ID
			if k == "" {
				k = fmt.Sprint(i)
				v.ID = k
			}
			nodesMap[k] = v
		}
		return nodesMap
	}
	return nil
}

func normalizeConanRef(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	// remove revisions and user/channel for dependency id purposes
	if i := strings.Index(ref, "#"); i > -1 {
		ref = ref[:i]
	}
	if i := strings.Index(ref, "@"); i > -1 {
		ref = ref[:i]
	}
	return ref
}

func parseConanRefToComponent(ref, name, version string) (string, string) {
	ref = normalizeConanRef(ref)
	if ref != "" {
		parts := strings.SplitN(ref, "/", 2)
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
		return ref, ""
	}
	return name, version
}

func parseConanGraphRequirements(n _ConanGraphNode) []_ConanGraphRequirement {
	reqs := append([]_ConanGraphRequirement{}, n.Requires...)
	if len(n.Dependencies) == 0 {
		return reqs
	}
	var depMap map[string]_ConanGraphRequirement
	if e := json.Unmarshal(n.Dependencies, &depMap); e == nil && len(depMap) > 0 {
		for k, v := range depMap {
			if v.ID == "" {
				v.ID = k
			}
			reqs = append(reqs, v)
		}
		return reqs
	}
	var depArray []_ConanGraphRequirement
	if e := json.Unmarshal(n.Dependencies, &depArray); e == nil && len(depArray) > 0 {
		reqs = append(reqs, depArray...)
	}
	return reqs
}

func (t _ConanGraphInfoJsonFile) Tree() (*model.DependencyItem, error) {
	nodes := parseConanGraphNodes(t.Graph.Nodes)
	if len(nodes) == 0 {
		return nil, ErrRootNodeNotFound
	}
	inDegree := make(map[string]int, len(nodes))
	for id := range nodes {
		inDegree[id] = 0
	}
	for _, n := range nodes {
		for _, req := range parseConanGraphRequirements(n) {
			if req.ID == "" {
				continue
			}
			if _, ok := inDegree[req.ID]; ok {
				inDegree[req.ID]++
			}
		}
	}
	rootID := ""
	for id, d := range inDegree {
		if d == 0 {
			rootID = id
			if nodes[id].Ref == "" {
				break
			}
		}
	}
	if rootID == "" {
		return nil, ErrRootNodeNotFound
	}
	var build func(id string, depth int, visited map[string]bool) *model.DependencyItem
	build = func(id string, depth int, visited map[string]bool) *model.DependencyItem {
		n, ok := nodes[id]
		if !ok {
			return nil
		}
		if visited[id] {
			return nil
		}
		visited[id] = true
		defer delete(visited, id)
		compName, compVersion := parseConanRefToComponent(n.Ref, n.Name, n.Version)
		item := &model.DependencyItem{
			Component: model.Component{
				CompName:    compName,
				CompVersion: compVersion,
				EcoRepo:     EcoRepo,
			},
		}
		if depth == 1 {
			item.DependencyRelation = model.DependencyRelationDirect
		} else if depth > 1 {
			item.DependencyRelation = model.DependencyRelationTransitive
		}
		for _, req := range parseConanGraphRequirements(n) {
			if req.ID == "" {
				continue
			}
			child := build(req.ID, depth+1, visited)
			if child == nil {
				cn, cv := parseConanRefToComponent(req.Ref, req.Name, req.Version)
				if cn == "" {
					continue
				}
				child = &model.DependencyItem{
					Component: model.Component{
						CompName:    cn,
						CompVersion: cv,
						EcoRepo:     EcoRepo,
					},
				}
				if depth+1 == 1 {
					child.DependencyRelation = model.DependencyRelationDirect
				} else if depth+1 > 1 {
					child.DependencyRelation = model.DependencyRelationTransitive
				}
			}
			item.Dependencies = append(item.Dependencies, *child)
		}
		return item
	}
	root := build(rootID, 0, map[string]bool{})
	if root == nil {
		return nil, ErrRootNodeNotFound
	}
	return root, nil
}

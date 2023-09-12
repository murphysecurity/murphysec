package maven

import (
	"encoding/json"
	"os"
)

// PluginGraphOutput maven dependency-graph.json file
type PluginGraphOutput struct {
	GraphName string `json:"graphName"`
	Artifacts []struct {
		GroupId    string   `json:"groupId"`
		ArtifactId string   `json:"artifactId"`
		Optional   bool     `json:"optional"`
		Scopes     []string `json:"scopes"`
		Version    string   `json:"version"`
	} `json:"artifacts"`
	Dependencies []struct {
		NumericFrom int `json:"numericFrom"`
		NumericTo   int `json:"numericTo"`
	} `json:"dependencies"`
}

// ReadFromFile dependency-graph.json
func (d *PluginGraphOutput) ReadFromFile(path string) error {
	data, e := os.ReadFile(path)
	if e != nil {
		return ErrBadDepsGraph.DetailedWrap("read graph file", e)
	}
	g := PluginGraphOutput{}
	if e := json.Unmarshal(data, &g); e != nil {
		return ErrBadDepsGraph.Wrap(e)
	}
	*d = g
	return nil
}

func (d PluginGraphOutput) Tree() (*Dependency, error) {
	// from -> listOf to
	edges := d.edgesMap()
	root, e := d.findRootNode()
	if e != nil {
		return nil, e
	}
	visited := make([]bool, len(d.Artifacts))
	t := d._tree(root, visited, edges)
	if t != nil {
		return t, nil
	}
	return nil, ErrBadDepsGraph.Detailed("empty graph")
}

func (d PluginGraphOutput) _tree(id int, visitedId []bool, edges map[int][]int) *Dependency {
	if visitedId[id] {
		return nil
	}
	visitedId[id] = true
	defer func() { visitedId[id] = false }()

	r := &Dependency{
		Coordinate: Coordinate{
			GroupId:    d.Artifacts[id].GroupId,
			ArtifactId: d.Artifacts[id].ArtifactId,
			Version:    d.Artifacts[id].Version,
		},
		Children: nil,
	}
	if len(d.Artifacts[id].Scopes) > 0 {
		r.Scope = d.Artifacts[id].Scopes[0]
	}
	for _, toNum := range edges[id] {
		t := d._tree(toNum, visitedId, edges)
		if t == nil {
			continue
		}
		r.Children = append(r.Children, *t)
	}
	return r
}

func (d PluginGraphOutput) edgesMap() (m map[int][]int) {
	m = map[int][]int{}
	distinctM := map[int64]struct{}{}
	for _, it := range d.Dependencies {
		d := int64(it.NumericFrom<<32) | int64(it.NumericTo)
		if _, ok := distinctM[d]; ok {
			continue
		}
		distinctM[d] = struct{}{}
		m[it.NumericFrom] = append(m[it.NumericFrom], it.NumericTo)
	}
	return
}

func (d PluginGraphOutput) findRootNode() (int, error) {
	candidate := make([]bool, len(d.Artifacts))
	for _, it := range d.Dependencies {
		if it.NumericTo > len(candidate) {
			return 0, ErrBadDepsGraph.Detailed("numeric_to > len(artifacts)")
		}
		candidate[it.NumericTo] = true
	}
	for idx, it := range candidate {
		if !it {
			return idx, nil
		}
	}
	return 0, ErrBadDepsGraph.Detailed("root node not found")
}

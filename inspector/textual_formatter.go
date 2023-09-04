package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"time"
)

type DependencyGraph struct {
	nodes     []model.Component
	nodeIndex map[model.Component]int
	edges     map[[2]int]struct{}
}

func (d *DependencyGraph) indexOf(component model.Component) int {
	if i, ok := d.nodeIndex[component]; ok {
		return i
	}
	i := len(d.nodes)
	d.nodeIndex[component] = i
	d.nodes = append(d.nodes, component)
	return i
}

func (d *DependencyGraph) putEdge(from, to model.Component) {
	var fromIndex = d.indexOf(from)
	var toIndex = d.indexOf(to)
	d.edges[[2]int{fromIndex, toIndex}] = struct{}{}
}

func (d *DependencyGraph) _visit(parent model.Component, dep model.DependencyItem) {
	d.putEdge(parent, dep.Component)
	for _, it := range dep.Dependencies {
		d._visit(it.Component, it)
	}
}

func BuildSpdx(task *model.ScanTask) []byte {
	var g = DependencyGraph{
		nodes:     nil,
		nodeIndex: make(map[model.Component]int),
		edges:     make(map[[2]int]struct{}),
	}
	for _, module := range task.Modules {
		for _, it := range module.Dependencies {
			g.indexOf(it.Component)
			for _, dependency := range it.Dependencies {
				g._visit(it.Component, dependency)
			}
		}
	}

	var packages = make([]any, 0)
	for i, node := range g.nodes {
		var m = map[string]any{
			"SPDXID":       fmt.Sprint("SPDXRef-package-", i),
			"name":         node.CompName,
			"versionInfo":  node.CompVersion,
			"fileAnalyzed": false,
		}
		var externalRefs map[string]any
		if node.Ecosystem == "maven" {
			externalRefs = map[string]any{
				"referenceCategory": "PACKAGE-MANAGER",
				"referenceLocator":  node.CompName + ":" + node.CompVersion,
				"referenceType":     "maven-central",
			}
		}
		if node.Ecosystem == "npm" {
			externalRefs = map[string]any{
				"referenceCategory": "PACKAGE-MANAGER",
				"referenceLocator":  node.CompName + "@" + node.CompVersion,
				"referenceType":     "npm",
			}
		}
		if node.Ecosystem == "go" {
			externalRefs = map[string]any{
				"referenceCategory": "PACKAGE-MANAGER",
				"referenceLocator":  "pkg:golang/" + node.CompName + "@" + node.CompVersion,
				"referenceType":     "purl",
			}
		}
		if externalRefs != nil {
			m["externalRefs"] = externalRefs
		}
		packages = append(packages, m)
	}
	var relationships = make([]any, 0)
	for ints := range g.edges {
		var from = ints[0]
		var to = ints[1]
		relationships = append(relationships, map[string]any{
			"relationshipType":   "DEPENDS_ON",
			"spdxElementId":      fmt.Sprint("SPDXRef-package-", from),
			"relatedSpdxElement": fmt.Sprint("SPDXRef-package-", to),
		})
	}
	var spdx = map[string]any{
		"SPDXID":      "SPDXRef-DOCUMENT",
		"spdxVersion": "SPDX-2.3",
		"creationInfo": map[string]any{
			"created":  time.Now(),
			"creators": []string{"Tool: murphysec-cli"},
		},
		"name":          "com.github.murphysecurity/murphysec",
		"dataLicense":   "CC0-1.0",
		"packages":      packages,
		"relationships": relationships,
	}
	data, e := json.MarshalIndent(spdx, "", "  ")
	if e != nil {
		panic("marshal spdx failed: " + e.Error())
	}
	return data
}

package model

type DependencyItem struct {
	Component
	Dependencies       []DependencyItem   `json:"dependencies,omitempty"`
	DependencyRelation DependencyRelation `json:"dependency_relation"`
	MavenScope         string             `json:"maven_scope,omitempty"`
	IsOnline           IsOnline           `json:"is_online"`
	IsDirectDependency bool               `json:"is_direct_dependency,omitempty"`
}

func (d *DependencyItem) Postprocess() {
	for i := range d.Dependencies {
		if d.Dependencies[i].DependencyRelation == DependencyRelationDirect {
			d.IsDirectDependency = true
		}
		d.Dependencies[i].Postprocess()
	}
}

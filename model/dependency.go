package model

type DependencyItem struct {
	Component
	Dependencies       []DependencyItem   `json:"dependencies,omitempty"`
	DependencyRelation DependencyRelation `json:"dependency_relation"`
	MavenScope         string             `json:"maven_scope,omitempty"`
	IsOnline           IsOnline           `json:"is_online"`
}

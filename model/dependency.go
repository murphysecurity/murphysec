package model

type DependencyItem struct {
	Component
	Dependencies       []DependencyItem `json:"dependencies,omitempty"`
	IsDirectDependency bool             `json:"is_direct_dependency,omitempty"`
	MavenScope         string           `json:"maven_scope,omitempty"`
	IsOnline           bool             `json:"is_online"`
}

package model

type DependencyItem struct {
	Component
	Dependencies       []DependencyItem `json:"dependencies,omitempty"`
	IsDirectDependency bool             `json:"is_direct_dependency,omitempty"`
}

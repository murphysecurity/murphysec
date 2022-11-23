package model

type DependencyItem struct {
	Component
	Dependencies []DependencyItem `json:"dependencies,omitempty"`
}

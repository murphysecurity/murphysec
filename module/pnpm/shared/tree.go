package shared

import "github.com/murphysecurity/murphysec/model"

type DepTree struct {
	Name         string
	Dependencies []model.DependencyItem
}

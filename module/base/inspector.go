package base

import (
	"context"
	"fmt"
)

type Inspector interface {
	fmt.Stringer
	CheckDir(dir string) bool
	InspectProject(ctx context.Context) error
	SupportFeature(feature InspectorFeature) bool
}

type InspectorFeature int

const (
	InspectorFeatureAllowNested InspectorFeature = 1 << iota
)

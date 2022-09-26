package model

import (
	"context"
	"fmt"
)

type InspectorFeature int

type Inspector interface {
	fmt.Stringer
	CheckDir(dir string) bool
	InspectProject(ctx context.Context) error
	SupportFeature(feature InspectorFeature) bool
}

const (
	InspectorFeatureAllowNested InspectorFeature = 1 << iota
)

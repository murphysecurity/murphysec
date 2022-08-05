package base

import (
	"context"
	"fmt"
)

type Inspector interface {
	fmt.Stringer
	CheckDir(dir string) bool
	InspectProject(ctx context.Context) error
	SupportFeature(feature Feature) bool
}

type Feature int

const (
	FeatureAllowNested Feature = 1 << iota
)

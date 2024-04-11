package shared

import "github.com/murphysecurity/murphysec/model"

type Component struct {
	Name     string
	Version  string
	Dev      bool
	Children []*Component
}

type GComponent struct {
	Name    string
	Version string
	Dev     bool
}

type GVisitor[T any] func(visitor DoVisit[T], parent *GComponent, child *GComponent, arg T) error
type DoVisit[T any] func(T) error

var EcoRepo = model.EcoRepo{Ecosystem: "npm"}

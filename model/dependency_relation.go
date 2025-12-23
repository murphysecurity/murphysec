package model

import (
	"encoding"
	"errors"
)

//go:generate stringer -linecomment -type DependencyRelation -output dependency_relation_string.go
type DependencyRelation int

func (i DependencyRelation) MarshalText() (text []byte, err error) {
	return []byte(i.String()), nil
}

func (i *DependencyRelation) UnmarshalText(text []byte) error {
	switch string(text) {
	case DependencyRelationUnknown.String():
		*i = DependencyRelationUnknown
		return nil
	case DependencyRelationDirect.String():
		*i = DependencyRelationDirect
		return nil
	case DependencyRelationTransitive.String():
		*i = DependencyRelationTransitive
		return nil
	default:
		return errors.New("bad dependency relation")
	}
}

const (
	DependencyRelationEmptyStr   DependencyRelation = 0 //
	DependencyRelationUnknown    DependencyRelation = 1 // Unknown
	DependencyRelationDirect     DependencyRelation = 2 // Direct
	DependencyRelationTransitive DependencyRelation = 3 // Transitive
)

var _ encoding.TextUnmarshaler = (*DependencyRelation)(nil)
var _ encoding.TextMarshaler = DependencyRelation(0)

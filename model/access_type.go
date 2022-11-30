package model

import (
	"errors"
	"strings"
)

type AccessType string

const (
	AccessTypeCli  AccessType = "cli"
	AccessTypeIdea AccessType = "idea"
)

func (i AccessType) Valid() bool {
	switch i {
	case AccessTypeIdea, AccessTypeCli:
		return true
	}
	return false
}

func (i *AccessType) Of(s string) error {
	switch strings.ToLower(s) {
	case "cli", "":
		*i = AccessTypeCli
	case "idea":
		*i = AccessTypeIdea
	default:
		return errors.New("bad access type")
	}
	return nil
}

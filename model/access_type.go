package model

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

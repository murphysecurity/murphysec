package model

type basePathKey struct{}

var BasePathKey = &basePathKey{}

func (basePathKey) String() string {
	return "basePath"
}

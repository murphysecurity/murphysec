package gradle

type ProjectFilter struct {
	ProjectNames []string
}

type _ProjectFilterKey struct{}

func (_ProjectFilterKey) String() string {
	return "GradleProjectFilter"
}

var ProjectFilterCtxKey = &_ProjectFilterKey{}

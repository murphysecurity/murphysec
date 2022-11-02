package maven

import "fmt"

const (
	ErrMvnDisabled          = mvnError("mvn command disabled")
	ErrMvnNotFound          = mvnError("mvn command not found")
	ErrCheckMvnVersion      = mvnError("eval mvn version failed")
	ErrBadDepsGraph         = mvnError("bad dependency graph")
	ErrInvalidCoordinate    = mvnError("invalid coordinate")
	ErrArtifactNotFound     = mvnError("artifact not found")
	ErrGetArtifactFailed    = mvnError("get artifact failed")
	ErrParsePomFailed       = mvnError("parse pom failed")
	ErrOpenProject          = mvnError("open project failed")
	ErrPomCircularDependent = mvnError("pom file circular dependent")
	ErrBadCoordinate        = mvnError("bad coordinate")
	ErrCouldNotResolve      = mvnError("couldn't resolve")
	ErrMvnExitErr           = mvnError("mvn command exit with non-zero code")
	ErrMvnCmd               = mvnError("error during mvn execution")
	ErrInspection           = mvnError("can't inspect the maven project")
)

type mvnError string

func (m mvnError) DetailedWrap(detail string, cause error) error {
	return &wrappedMvnErr{
		k: m,
		e: &detailedMvnErr{
			cause:  cause,
			detail: detail,
		},
	}
}

func (m mvnError) Detailed(detail string) error {
	return &detailedMvnErr{
		cause:  m,
		detail: detail,
	}
}

func (m mvnError) Wrap(cause error) error {
	return &wrappedMvnErr{
		k: m,
		e: cause,
	}
}

func (m mvnError) Error() string {
	return string(m)
}

type wrappedMvnErr struct {
	k mvnError
	e error
}

func (w wrappedMvnErr) Is(target error) bool {
	return w.k == target || w == target
}

func (w wrappedMvnErr) Unwrap() error {
	return w.e
}

func (w wrappedMvnErr) Error() string {
	return fmt.Sprintf("%s: %s", w.k.Error(), w.e.Error())
}

type detailedMvnErr struct {
	cause  error
	detail string
}

func (d detailedMvnErr) Unwrap() error {
	return d.cause
}
func (d detailedMvnErr) Error() string {
	return fmt.Sprintf("%s: %s", d.cause.Error(), d.detail)
}

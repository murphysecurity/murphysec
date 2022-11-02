package maven

import "fmt"

const (
	_                       mvnError = iota
	ErrMvnDisabled                   // maven: mvn command disabled
	ErrMvnNotFound                   // maven: mvn command not found
	ErrCheckMvnVersion               // maven: eval mvn version failed
	ErrBadDepsGraph                  // maven: bad dependency graph
	ErrInvalidCoordinate             // maven: invalid coordinate
	ErrArtifactNotFound              // maven: artifact not found
	ErrGetArtifactFailed             // maven: get artifact failed
	ErrParsePomFailed                // maven: parse pom failed
	ErrOpenProject                   // maven: open project failed
	ErrPomCircularDependent          // maven: pom file circular dependent
	ErrBadCoordinate                 // maven: bad coordinate
	ErrCouldNotResolve               // maven: couldn't resolve
	ErrMvnExitErr                    // maven: mvn command exit with non-zero code
	ErrMvnCmd                        // maven: error during mvn execution
	ErrInspection                    // maven: can't inspect the maven project
)

//go:generate stringer -linecomment -type mvnError -output err_string.go
type mvnError int

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
	return m.String()
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

package arkts

//go:generate stringer -linecomment -type textError -output error_string.go
type textError int

func (t textError) Error() string {
	return t.String()
}

const (
	ErrReadPackage     textError = iota + 1 // arkts: read oh-package.json5 failed
	ErrReadPackageLock                      // arkts: read oh-package-lock.json5 failed
	ErrRootNotFound                         // arkts: root package not found, and no oh-package.json5 found
)

func (t textError) withCause(e error) error {
	return Error{t, e}
}

type Error struct {
	Kind   textError
	Unwrap error
}

func (a Error) Error() string {
	if a.Unwrap == nil {
		return a.Kind.String()
	}
	return a.Kind.String() + ", " + a.Unwrap.Error()
}

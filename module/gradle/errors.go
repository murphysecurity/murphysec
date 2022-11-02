package gradle

//go:generate stringer -type _e -linecomment -output errors_string.go
type _e int

func (e _e) Error() string {
	return e.String()
}

const (
	_                    _e = iota
	ErrNoGradle             // gradle: No gradle found
	ErrEvalGradleVersion    // gradle: Eval gradle version failed
)

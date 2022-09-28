package gradle

type _e string

func (e _e) Error() string {
	return string(e)
}

const (
	ErrNoGradle          _e = "No gradle found"
	ErrEvalGradleVersion _e = "Eval gradle version failed"
)

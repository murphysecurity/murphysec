package logger

const (
	ErrCreateLogFileFailed _e = "create log file failed"
)

type _e string

func (e _e) Error() string {
	return string(e)
}

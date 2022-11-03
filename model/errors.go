package model

//go:generate stringer -linecomment -type _e -output errors_string.go
type _e int

const (
	_ _e = iota
	ErrNoGitRemoteFound // no git remote found
	ErrNoGitRepo        // no git repo found
)

func (e _e) Error() string {
	return e.String()
}

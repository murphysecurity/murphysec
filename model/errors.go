package model

type _err string

const (
	ErrNoGitRemoteFound _err = "No git remote found"
	ErrNoGitRepo        _err = "No git repo found"
)

func (e _err) Error() string {
	return string(e)
}

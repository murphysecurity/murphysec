package gitinfo

//go:generate stringer -type _e -linecomment -output error_string.go
type _e int

const (
	_                    _e = iota
	ErrNoRepoFound          // gitinfo: no repo
	_ErrNoRemoteURLFound    // gitinfo: no remote url
)

func (i _e) Error() string {
	return i.String()
}

package maven

//go:generate stringer -type remoteError -linecomment -output remote_error_string.go
type remoteError int

const (
	_                   remoteError = iota
	ErrRemoteNoResource             // m2remote: no resource
)

func (i remoteError) Error() string {
	return i.String()
}

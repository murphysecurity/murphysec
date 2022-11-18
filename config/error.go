package config

//go:generate stringer -type _e -output error_string.go -linecomment
type _e int

const (
	_                       _e = iota
	_ErrTokenFileReadFailed    // config: token can't be read
	ErrNoToken                 // config: no token
	ErrBadToken                // config: bad token
	_ErrTokenFileNotAFile      // config: token file not a file
)

func (i _e) Error() string {
	return i.String()
}

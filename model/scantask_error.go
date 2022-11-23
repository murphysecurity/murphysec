package model

//go:generate stringer -type scanTaskError -linecomment -output scantask_error_string.go
type scanTaskError int

const (
	_                     scanTaskError = iota
	ErrPathIsNotAbsolute                // task: path must be absolute
	ErrMustBeDirectory                  // task: path must be a directory
	ErrMustBeFile                       // task: path must be a file
	ErrSubtaskNameIsEmpty               // task: subtask name is empty
)

func (i scanTaskError) Error() string { return i.String() }

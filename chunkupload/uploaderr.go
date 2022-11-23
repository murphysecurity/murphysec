package chunkupload

//go:generate stringer -type uploadErr -linecomment
type uploadErr int

const (
	_              uploadErr = iota
	ErrDirInvalid            // uploader: dir invalid
	ErrEvalAbsPath           // uploader: cannot evaluate absolute path
)

func (i uploadErr) Error() string {
	return i.String()
}

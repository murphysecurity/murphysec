package bundler

//go:generate stringer -type parserError -linecomment -output gemlock_parser_error_string.go
type parserError int

const (
	_            parserError = iota
	ErrBadIndent             // Bad indent
	ErrParseFail             // ParseFail
)

func (i parserError) Error() string {
	return i.String()
}

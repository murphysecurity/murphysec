package errors

import goerr "errors"

func Is(err, target error) bool {
	return goerr.Is(err, target)
}

func As(err error, target any) bool {
	//goland:noinspection GoErrorsAs
	return goerr.As(err, target)
}

func Unwrap(err error) error {
	return goerr.Unwrap(err)
}

func New(text string) error {
	return goerr.New(text)
}

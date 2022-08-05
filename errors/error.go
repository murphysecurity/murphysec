package errors

import (
	goerr "errors"
)

type detailedError struct {
	detail string
	e      error
}

func (e *detailedError) Unwrap() error {
	return e.e
}

func (e *detailedError) Error() string {
	return e.e.Error() + ": " + e.detail
}

// WithDetail returns an error with detailed text
func WithDetail(e error, detail string) error {
	return &detailedError{
		detail: detail,
		e:      e,
	}
}

type causeErr struct {
	err   error
	cause error
}

func (e *causeErr) Is(target error) bool {
	return e == target || goerr.Is(e.err, target)
}

func (e *causeErr) Unwrap() error {
	return e.cause
}

func (e *causeErr) Error() string {
	a := e.err.Error()
	b := e.cause.Error()
	if a == "" {
		return b
	}
	return a + ": " + b
}

// WithCause returns an error wrap a cause
func WithCause(err error, cause error) error {
	if err == nil {
		panic("err mustn't be nil")
	}
	if cause == nil {
		panic("cause mustn't be nil")
	}
	return &causeErr{
		err:   err,
		cause: cause,
	}
}

type wrapped struct {
	prefix string
	e      error
}

func (e *wrapped) Unwrap() error {
	return e.e
}

func (e *wrapped) Error() string {
	if e.prefix == "" {
		return e.e.Error()
	}
	return e.prefix + ": " + e.e.Error()
}

// Wrap returns an error with prefix
func Wrap(e error, prefix string) error {
	if e == nil {
		panic("e mustn't be nil")
	}
	return &wrapped{
		prefix: prefix,
		e:      e,
	}
}

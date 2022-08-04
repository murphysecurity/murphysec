package errors

import (
	goerr "errors"
)

func New(s string) error {
	return &E{s}
}

type E struct {
	s string
}

func (e *E) Error() string {
	return e.s
}

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
	m string
	e error
}

func (e *wrapped) Unwrap() error {
	return e.e
}

func (e *wrapped) Error() string {
	if e.m == "" {
		return e.e.Error()
	}
	return e.m + ": " + e.e.Error()
}

func Wrap(e error, m string) error {
	if e == nil {
		panic("e mustn't be nil")
	}
	return &wrapped{
		m: m,
		e: e,
	}
}

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

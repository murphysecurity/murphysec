package api

import (
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/tlserr"
	"net/url"
	"strconv"
)

//go:generate stringer -type apiError -linecomment
type apiError int

const (
	_                        apiError = iota
	ErrTLSError                       // api: tls error
	ErrTimeout                        // api: timeout
	ErrRequest                        // api: request fail
	ErrUnprocessableResponse          // api: cannot process server response
	ErrTokenInvalid                   // api: token invalid
	ErrBadURL                         // api: bad URL
	ErrServerFail                     // api: server fail
	ErrGeneralError                   // api: general error
	ErrValidateFail                   // api: validation fail
)

func (i apiError) Error() string {
	return i.String()
}

type Error struct {
	Cause error `json:"cause"`
	// HTTPStatus copy the response status code
	HTTPStatus int `json:"http_status"`
	// UnprocessableResponse indicates an unprocessable response body, the response is not a recognized error
	UnprocessableResponse bool `json:"unprocessable_response"`
	// Server side error code
	Code int `json:"code,omitempty"`
	// Server side message
	Message string `json:"message,omitempty"`
}

func (i *Error) Error() string {
	if i.UnprocessableResponse {
		s := "API error, http status: " + strconv.Itoa(i.HTTPStatus) + ". The response can't be processed"
		if i.Cause != nil {
			s += ": " + i.Cause.Error()
		} else {
			s += "."
		}
		return s
	}
	if i.Cause != nil {
		return i.Cause.Error()
	}

	return fmt.Sprintf("[%d/%d]%s", i.HTTPStatus, i.Code, i.Message)
}

func (i *Error) Is(target error) bool {
	switch target {
	case ErrTokenInvalid:
		return i.HTTPStatus == 401
	case ErrUnprocessableResponse:
		return i.UnprocessableResponse
	case ErrServerFail:
		return i.HTTPStatus >= 500
	case ErrGeneralError:
		return i.HTTPStatus >= 400 && i.HTTPStatus < 500
	case ErrTLSError:
		return tlserr.IsTLSError(i.Cause)
	case ErrRequest:
		return true
	case ErrTimeout:
		r, ok := i.Cause.(*url.Error)
		return ok && r.Timeout()
	}
	if r, ok := target.(*Error); ok {
		if i.Code == r.Code {
			return true
		}
		if i.Cause != nil && r.Cause != nil {
			return errors.Is(i.Cause, r.Cause)
		}
	}
	return false
}

func (i *Error) Unwrap() error {
	return i.Cause
}

func errorOf(e error) error {
	return &Error{
		Cause: e,
	}
}

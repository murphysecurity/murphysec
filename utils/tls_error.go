package utils

import (
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

func IsTlsCertError(e error) bool {
	var a0 x509.HostnameError
	if errors.As(e, &a0) {
		return true
	}
	fmt.Println(reflect.TypeOf(e), reflect.TypeOf(errors.Unwrap(e)), reflect.TypeOf(errors.Unwrap(errors.Unwrap(e))), errors.As(e, &a0))
	var a1 x509.CertificateInvalidError
	if errors.As(e, &a1) {
		return true
	}
	var a2 x509.UnknownAuthorityError
	if errors.As(e, &a2) {
		return true
	}
	var a3 x509.SystemRootsError
	if errors.As(e, &a3) {
		return true
	}
	var a4 x509.ConstraintViolationError
	if errors.As(e, &a4) {
		return true
	}
	var a5 x509.ConstraintViolationError
	if errors.As(e, &a5) {
		return true
	}
	return false
}

package utils

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"regexp"
	"runtime"
)

var __x509Pattern = regexp.MustCompile("x509: .+ certificate is not standards compliant")

func IsTlsCertError(e error) bool {
	var a0 x509.HostnameError
	if errors.As(e, &a0) {
		return true
	}
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
	var a5 x509.InsecureAlgorithmError
	if errors.As(e, &a5) {
		return true
	}
	// Workaround for macOS because of: https://github.com/golang/go/issues/51991
	if runtime.GOOS == "darwin" && e != nil && __x509Pattern.MatchString(e.Error()) {
		return true
	}
	return false
}

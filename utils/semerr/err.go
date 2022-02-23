package semerr

import (
	"fmt"
	"github.com/savsgio/gotils/nocopy"
	"io"
)

func New(name string) *SemErr {
	return &SemErr{
		NoCopy: nocopy.NoCopy{},
		name:   name,
	}
}

type SemErr struct {
	nocopy.NoCopy
	name string
}

func (s *SemErr) Is(e error) bool {
	if s == e {
		return true
	}
	if ee, ok := e.(*Instance); ok {
		return ee.semErr == s || ee.wrappedErr == s
	}
	return false
}

func (s *SemErr) Decorate(e error) error {
	if _, ok := e.(*SemErr); ok {
		panic("Don't decorate a SemErr")
	}
	return &Instance{
		semErr:     s,
		wrappedErr: e,
	}
}

func (s *SemErr) Error() string {
	return s.name
}

type Instance struct {
	semErr     *SemErr
	wrappedErr error
}

func (i *Instance) Is(e error) bool {
	if e == i || e == i.wrappedErr || e == i.semErr {
		return true
	}
	if ee, ok := e.(*Instance); ok {
		return ee.semErr == i.semErr
	}
	return false
}

func (i *Instance) Unwrap() error {
	return i.wrappedErr
}

func (i *Instance) Error() string {
	return i.semErr.name + ": " + i.wrappedErr.Error()
}
func (i *Instance) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, i.wrappedErr.Error())
			_, _ = fmt.Fprintf(s, "%+v", i.wrappedErr)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, i.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", i.Error())
	default:
		_, _ = fmt.Fprintf(s, "[bad format]")
	}
}

package semerr

import (
	"github.com/savsgio/gotils/nocopy"
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

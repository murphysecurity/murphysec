package must

import (
	"io"
	"reflect"
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func M(e error) {
	if e != nil {
		panic(e)
	}
}

func M1[T1 any](t1 T1, e error) T1 {
	if e != nil {
		panic(e)
	}
	return t1
}

func M2[T1, T2 any](t1 T1, t2 T2, e error) (T1, T2) {
	if e != nil {
		panic(e)
	}
	return t1, t2
}

func A[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

// True checks b is true.
func True(b bool) {
	if !b {
		panic("assertion not true")
	}
}

// Close closes the file and panic on error.
// Useful in defer statement.
func Close(c io.Closer) {
	Must(c.Close())
}

func NotNil(a any) {
	if reflect.ValueOf(a).IsNil() {
		panic("value is nil")
	}
}

func NotZero(a any) {
	if reflect.ValueOf(a).IsZero() {
		panic("value is zero")
	}
}

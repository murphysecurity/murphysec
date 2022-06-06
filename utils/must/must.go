package must

import (
	"io"
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
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

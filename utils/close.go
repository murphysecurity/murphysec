package utils

import "io"

func Close(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

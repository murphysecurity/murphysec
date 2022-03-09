package utils

import "io"

func Close(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

func CloseToChan(closer io.Closer, errCh chan error) error {
	if e := closer.Close(); e != nil {
		errCh <- e
		return e
	}
	return nil
}

package utils

import "net/url"

func IsHttpTimeout(e error) bool {
	r, ok := e.(*url.Error)
	return ok && r.Timeout()
}

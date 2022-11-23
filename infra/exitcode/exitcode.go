package exitcode

import "os"

var _code int

func Set(code int) {
	_code = code
}

func Exit() {
	os.Exit(_code)
}

func Get() int {
	return _code
}

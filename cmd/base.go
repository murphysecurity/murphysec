package cmd

var _exitCode = 0

func SetGlobalExitCode(code int) {
	_exitCode = code
}

func GetGlobalExitCode() int {
	return _exitCode
}

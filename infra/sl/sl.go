package sl

func StringIsEmpty(s string) bool {
	return s == ""
}

func StringNotEmpty(s string) bool {
	return s != ""
}

func NotF1[T any](f func(T) bool) func(T) bool {
	return func(t T) bool {
		return !f(t)
	}
}

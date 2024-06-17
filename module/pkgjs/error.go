package pkgjs

type wrapped struct {
	e error
}

func (e wrapped) Unwrap() error {
	return e.e
}

func (e wrapped) Error() string {
	return "pkgjs: " + e.Error()
}

func wrap(e error) error {
	if //goland:noinspection GoTypeAssertionOnErrors
	_, ok := e.(wrapped); ok {
		return e
	}
	return wrapped{e}
}

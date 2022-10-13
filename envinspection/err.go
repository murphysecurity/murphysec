package envinspection

type _err string

func (e _err) Error() string {
	return string(e)
}

const (
	ErrGetHostname _err = "Cannot get hostname"
)

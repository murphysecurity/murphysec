// Code generated by "stringer -type remoteError -linecomment -output remote_error_string.go"; DO NOT EDIT.

package maven

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrRemoteNoResource-1]
}

const _remoteError_name = "m2remote: no resource"

var _remoteError_index = [...]uint8{0, 21}

func (i remoteError) String() string {
	i -= 1
	if i < 0 || i >= remoteError(len(_remoteError_index)-1) {
		return "remoteError(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _remoteError_name[_remoteError_index[i]:_remoteError_index[i+1]]
}
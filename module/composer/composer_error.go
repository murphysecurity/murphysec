package composer

import (
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
)

//go:generate stringer -type _e -linecomment -output composer_error_string.go
type _e int

const (
	_                        _e = iota
	ErrReadComposerManifest     // read composer.json failed
	ErrParseComposerManifest    // parsing composer.json failed
	ErrNoComposerFound          // no composer found
)

func (i _e) Error() string {
	return i.String()
}

type ce struct {
	key    error
	reason error
}

func (c *ce) Error() string {
	return fmt.Sprintf("%s Caused by: %s", c.key.Error(), c.reason.Error())
}

func (c *ce) Is(target error) bool {
	return errors.Is(c.key, target)
}

func (c *ce) Unwrap() error {
	return c.reason
}

func wrapErr(key error, reason error) error {
	if key == nil || reason == nil {
		panic("key == nil || reason == nil")
	}
	return &ce{
		key:    key,
		reason: reason,
	}
}

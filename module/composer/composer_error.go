package composer

import (
	"fmt"
	"github.com/pkg/errors"
)

var ErrReadComposerManifest = errors.New("Read composer.json failed")
var ErrParseComposerManifest = errors.New("Parsing composer.json failed")
var ErrComposerResolveFail = errors.New("PHP composer resolve failed")

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

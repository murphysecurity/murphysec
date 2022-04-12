package eg

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type EG errgroup.Group

func (e *EG) Wait() error {
	return (*errgroup.Group)(e).Wait()
}

func (e *EG) Go(f func()) {
	(*errgroup.Group)(e).Go(func() (err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprint("Panic: ", r))
			}
		}()
		f()
		return nil
	})
}

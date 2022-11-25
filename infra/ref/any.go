package ref

import "reflect"

func OmitZero[T any](t T) *T {
	if reflect.ValueOf(t).IsZero() {
		return nil
	}
	return &t
}

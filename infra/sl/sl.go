package sl

import "github.com/repeale/fp-go"

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

func FilterNotNull[T any](input []*T) []*T {
	return fp.Filter(func(t *T) bool { return t != nil })(input)
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	var r []Entry[K, V]
	for k, v := range m {
		r = append(r, Entry[K, V]{k, v})
	}
	return r
}

func AssociateBy[K comparable, V any](input []V, keySelector func(v V) K) map[K]V {
	var m = make(map[K]V)
	for _, v := range input {
		m[keySelector(v)] = v
	}
	return m
}

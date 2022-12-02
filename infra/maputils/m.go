package maputils

import (
	"sort"
)

type sortable interface {
	~uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | uintptr | string
}

type w[T sortable, E any] struct {
	ks []T
	vs []E
}

func (s w[T, E]) Len() int {
	return len(s.vs)
}

func (s w[T, E]) Swap(i, j int) {
	s.ks[i], s.ks[j] = s.ks[j], s.ks[i]
	s.vs[i], s.vs[j] = s.vs[j], s.vs[i]
}

func (s w[T, E]) Less(i, j int) bool {
	return s.ks[i] < s.ks[j]
}

func KeysSortedByValue[K comparable, V sortable](m map[K]V) []K {
	var r = w[V, K]{
		ks: make([]V, 0, len(m)),
		vs: make([]K, 0, len(m)),
	}
	for k, v := range m {
		r.ks = append(r.ks, v)
		r.vs = append(r.vs, k)
	}
	sort.Sort(r)
	return r.vs
}

func ValuesSortedByKey[K sortable, V any](m map[K]V) []V {
	var r = w[K, V]{
		ks: make([]K, 0, len(m)),
		vs: make([]V, 0, len(m)),
	}
	for k, v := range m {
		r.ks = append(r.ks, k)
		r.vs = append(r.vs, v)
	}
	sort.Sort(r)
	return r.vs
}

func Keys[K comparable, V any](m map[K]V) []K {
	var r = make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func Values[K comparable, V any](m map[K]V) []V {
	var r = make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

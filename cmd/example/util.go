package main

import (
	"cmp"
	"slices"

	"github.com/octohelm/gio-compose/pkg/compose"
)

func MapIndexed[E any, T any](list []E, m func(e E, i int) T) []T {
	out := make([]T, len(list))
	for i := range list {
		out[i] = m(list[i], i)
	}
	return out
}

func SortedMap[K cmp.Ordered, E any, O any](m map[K]E, each func(e E, k K) O) []O {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	out := make([]O, len(keys))

	for i, k := range keys {
		out[i] = each(m[k], k)
	}

	return out
}

func SafeBuild[T any](v T, build func(c T) compose.VNode) compose.VNode {
	if any(v) != nil {
		return build(v)
	}
	return nil
}

package iter

func Map[E any, T any](list []E, m func(e E) T) []T {
	out := make([]T, len(list))
	for i := range list {
		out[i] = m(list[i])
	}
	return out
}

func MapIndexed[E any, T any](list []E, m func(e E, i int) T) []T {
	out := make([]T, len(list))
	for i := range list {
		out[i] = m(list[i], i)
	}
	return out
}

func Filter[E any](list []E, filter func(e E) bool) []E {
	out := make([]E, 0, len(list))
	for i := range list {
		if filter(list[i]) {
			out = append(out, list[i])
		}
	}
	return out
}

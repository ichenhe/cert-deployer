package domain

func MapSlice[S any, T any](src []S, mapper func(S) T) []T {
	n := make([]T, len(src))
	for i, e := range src {
		n[i] = mapper(e)
	}
	return n
}

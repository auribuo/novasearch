package sql

func Skip[T any](source []T, count int) []T {
	return source[count:]
}

func SkipWhile[T any](source []T, predicate func(T) bool) []T {
	for i, item := range source {
		if !predicate(item) {
			return source[i:]
		}
	}
	return []T{}
}

func SkipOut[T any](source []T) ([]T, T) {
	return source[1:], source[0]
}

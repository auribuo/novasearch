package sql

func Take[T any](source []T, count int) []T {
	return source[:len(source)-count]
}

func TakeWhile[T any](source []T, predicate func(T) bool) []T {
	res := make([]T, 0)
	for i, item := range source {
		if predicate(item) {
			res = append(res, source[i])
		}
	}
	return source
}

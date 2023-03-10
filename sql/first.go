package sql

func First[T any](source []T) T {
	return source[0]
}

func FirstF[T any](source []T, predicate func(T) bool) *T {
	for _, item := range source {
		if predicate(item) {
			return &item
		}
	}
	return nil
}

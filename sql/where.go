package sql

func Where[T any](source []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range source {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

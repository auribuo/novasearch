package sql

func Map[T, R any](source []T, mapper func(T) R) []R {
	result := make([]R, len(source))
	for i, item := range source {
		result[i] = mapper(item)
	}
	return result
}

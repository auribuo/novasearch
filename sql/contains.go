package sql

func Contains[T comparable](source []T, item T) bool {
	for _, i := range source {
		if i == item {
			return true
		}
	}
	return false
}

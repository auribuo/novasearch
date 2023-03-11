package sql

func Keys[K comparable, V any](source map[K]V) []K {
	res := make([]K, 0, len(source))
	for k := range source {
		res = append(res, k)
	}
	return res
}

func Values[K comparable, V any](source map[K]V) []V {
	res := make([]V, 0, len(source))
	for _, v := range source {
		res = append(res, v)
	}
	return res
}

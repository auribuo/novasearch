package sql

func Keys[TKey comparable, TValue any](source map[TKey]TValue) []TKey {
	result := make([]TKey, 0, len(source))
	for key := range source {
		result = append(result, key)
	}
	return result
}

func Values[TKey comparable, TValue any](source map[TKey]TValue) []TValue {
	result := make([]TValue, 0, len(source))
	for _, value := range source {
		result = append(result, value)
	}
	return result
}

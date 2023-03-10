package sql

func Chunk[T any](source []T, size int) [][]T {
	result := make([][]T, 0)
	for i := 0; i < len(source); i += size {
		end := i + size
		if end > len(source) {
			end = len(source)
		}
		result = append(result, source[i:end])
	}
	return result
}

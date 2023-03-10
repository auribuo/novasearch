package util

func SliceBasedOn(source string, lengths []int) []string {
	result := make([]string, len(lengths))
	start := 0

	for i, length := range lengths {
		if start+length > len(source) {
			result[i] = ""
			continue
		}

		token := source[start : start+length]
		result[i] = token
		atLast := 0
		if lengths[len(lengths)-1] == length {
			atLast = 1
		}
		start += length + 1 + atLast
	}

	return result
}

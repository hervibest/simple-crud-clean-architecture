package converter

// Generic map function
func Map[T any, U any](input []T, transform func(T) U) []U {
	var result []U
	for _, item := range input {
		result = append(result, transform(item))
	}
	return result
}

package common

func Filter[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Remove[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if !f(v) {
			result = append(result, v)
		}
	}
	return result
}

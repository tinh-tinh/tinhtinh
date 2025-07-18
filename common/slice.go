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

func Map[T any, R any](in []T, fn func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

func Find[T any](list []T, match func(T) bool) (T, bool) {
	var zero T
	for _, item := range list {
		if match(item) {
			return item, true
		}
	}
	return zero, false
}

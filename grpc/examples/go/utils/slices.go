package utils

func At[T any](s []T, i int) (T, bool) {
	l := len(s)
	if i < 0 {
		i = l + i
	}
	if i < 0 || i >= l {
		var zero T
		return zero, false
	}
	return s[i], true
}

func ForEach[T any](collection []T, iteratee func(item T, index int)) {
	for i := range collection {
		iteratee(collection[i], i)
	}
}

func Filter[T any](collection []T, predicate func(item T, index int) bool) []T {
	result := make([]T, 0, len(collection))
	for i := range collection {
		if predicate(collection[i], i) {
			result = append(result, collection[i])
		}
	}
	return result
}

func Map[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))
	for i := range collection {
		result[i] = iteratee(collection[i], i)
	}
	return result
}

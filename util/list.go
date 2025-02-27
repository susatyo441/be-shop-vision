package util

func Find[T any](list []T, cb func(i T) bool) *T {
	for _, item := range list {
		if cb(item) {
			return &item
		}
	}

	return nil
}

func Map[T, R any](list []T, cb func(i T) R) []R {
	result := []R{}
	for _, item := range list {
		result = append(result, cb(item))
	}

	return result
}

package list

// //////////////////////////////////////////////////
// filter

func Filter[T any](items []T, filterFn func(item T) bool) []T {
	if items == nil {
		return nil
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		if filterFn == nil || filterFn(item) {
			result = append(result, item)
		}
	}
	return result
}

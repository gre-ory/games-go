package dict

// //////////////////////////////////////////////////
// filter

func Filter[T comparable, U any](items map[T]U, filterFn func(item U) bool) []U {
	if items == nil {
		return nil
	}
	result := make([]U, 0, len(items))
	for _, item := range items {
		if filterFn(item) {
			result = append(result, item)
		}
	}
	return result
}

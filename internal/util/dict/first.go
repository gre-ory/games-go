package dict

// //////////////////////////////////////////////////
// first

func First[T comparable, U any](items map[T]U, filterFn func(item U) bool) (U, bool) {
	if items == nil {
		var empty U
		return empty, false
	}
	for _, item := range items {
		if filterFn(item) {
			return item, true
		}
	}
	var empty U
	return empty, false
}

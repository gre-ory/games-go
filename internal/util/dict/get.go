package dict

// //////////////////////////////////////////////////
// get

func Get[T comparable, U any](items map[T]U, key T) (U, bool) {
	if items == nil {
		var empty U
		return empty, false
	}
	if item, ok := items[key]; ok {
		return item, true
	}
	var empty U
	return empty, false
}

func MustGet[T comparable, U any](items map[T]U, key T) U {
	result, _ := Get(items, key)
	return result
}

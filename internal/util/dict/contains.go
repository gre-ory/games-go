package dict

// //////////////////////////////////////////////////
// contains

func ContainsKey[T comparable, U any](items map[T]U, key T) bool {
	if _, ok := items[key]; ok {
		return true
	}
	return false
}

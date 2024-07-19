package dict

// //////////////////////////////////////////////////
// convert

func ConvertToList[T comparable, U any, V any](items map[T]U, convertFn func(key T, value U) V) []V {
	if items == nil {
		return nil
	}
	result := make([]V, 0, len(items))
	for key, value := range items {
		result = append(result, convertFn(key, value))
	}
	return result
}

func Key[T comparable, U any](key T, value U) T {
	return key
}

func Value[T comparable, U any](key T, value U) U {
	return value
}

func Values[T comparable, U any](items map[T]U) []U {
	return ConvertToList(items, Value)
}

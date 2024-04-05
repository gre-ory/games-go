package list

// //////////////////////////////////////////////////
// converter

// A Converter is a function that converts one object into one other object
type Converter[V any, O any] func(V) O

// //////////////////////////////////////////////////
// convert

// Convert converts a list of items into a list of other items
func Convert[T any, U any](items []T, convert Converter[T, U]) []U {
	converted := make([]U, 0, len(items))
	for _, item := range items {
		converted = append(converted, convert(item))
	}
	return converted
}

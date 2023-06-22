package slice

// Contains returns true if a slice contains an element
func Contains[T comparable](slice []T, e T) bool {
	for _, elem := range slice {
		if elem == e {
			return true
		}
	}
	return false
}

// Transform takes a slice of a given type, and a function to
// perform on each element to transform it into a different type
func Transform[T any, V any](slice []T, fn func(T) V) []V {
	vs := []V{}
	for _, e := range slice {
		vs = append(vs, fn(e))
	}
	return vs
}

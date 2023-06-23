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

// Map takes a slice of a given type and a function
// to extract a map key and value from slice elements,
// returns a new map with those key-value pairs. Note
// that if multiple elements in the slice return the
// same key, the element that appears last in the slice
// will be the (only) element that maps to the key.
func Map[T any, K comparable, V any](slice []T, fn func(T) (K, V)) map[K]V {
	m := map[K]V{}
	for _, e := range slice {
		k, v := fn(e)
		m[k] = v
	}
	return m
}

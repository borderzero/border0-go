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

// CountWhere returns the number of elements in a slice that satisfy a given condition.
func CountWhere[T any](slice []T, condition func(T) bool) int {
	count := 0
	for _, e := range slice {
		if condition(e) {
			count++
		}
	}
	return count
}

// Diff finds the difference between two slices, `original` and `changed`. It returns two new slices for the changes.
// The first slice contains new items that are in `changed` but not in `original`, and the second slice contains removed
// items that are in `original` but not in `changed`.
func Diff[T comparable](original, changed []T) (newItems, removedItems []T) {
	originalMap := make(map[T]bool)
	changedMap := make(map[T]bool)

	for _, item := range original {
		originalMap[item] = true
	}
	for _, item := range changed {
		changedMap[item] = true
	}

	// find items that are in b but not in a (new items)
	for _, item := range changed {
		if _, exists := originalMap[item]; !exists {
			newItems = append(newItems, item)
		}
	}

	// find items that are in a but not in b (removed items)
	for _, item := range original {
		if _, exists := changedMap[item]; !exists {
			removedItems = append(removedItems, item)
		}
	}

	return newItems, removedItems
}

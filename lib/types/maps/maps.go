package maps

// EnsureNotNil ensures a given map is not nil. Returns a
// new empty (but non-nil) map if the original map was nil.
func EnsureNotNil[K comparable, V any](m map[K]V) map[K]V {
	if m != nil {
		return m
	}
	return make(map[K]V)
}

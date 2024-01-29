package maps

// Reverse returns a new map where the keys and values are reversed
// from the original input map. This is useful when we want to index
// a particular search on both keys and values e.g. mapping of public
// identifiers to internal identifiers and viceversa.
//
// Note that if values are not unique in the input map, the retured
// map will only have a key-value pair for the last ocurrence of the
// value (which is the key in the generated map) in the input map.
func Reverse[K, V comparable](m map[K]V) map[V]K {
	if m == nil {
		return nil
	}
	reversed := make(map[V]K)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

// EnsureNotNil ensures a given map is not nil. Returns a
// new empty (but non-nil) map if the original map was nil.
func EnsureNotNil[K comparable, V any](m map[K]V) map[K]V {
	if m != nil {
		return m
	}
	return make(map[K]V)
}

// MatchesFilters returns true if the given map
// matches the given inclusion and exclusion filters.
// This requires that the map's values are comparable.
func MatchesFilters[K comparable, V comparable](
	kv map[K]V,
	inclusion map[K][]V,
	exclusion map[K][]V,
) bool {
	included := (inclusion == nil)
	excluded := false

	if inclusion != nil {
		for key, value := range kv {
			if matchesFilter(key, value, inclusion) {
				included = true
				break
			}
		}
	}

	if exclusion != nil {
		for key, value := range kv {
			if matchesFilter(key, value, exclusion) {
				excluded = true
				break
			}
		}
	}

	return included && !excluded
}

// matchesFilter returns true if a given key-value
// pair matches a given filter of key-value options.
func matchesFilter[K comparable, V comparable](
	key K,
	value V,
	filter map[K][]V,
) bool {
	for filterKey, filterValues := range filter {
		if key == filterKey {
			// we interpret empty filter values
			// as "match any value of the key"
			if len(filterValues) == 0 {
				return true
			}
			for _, filterValue := range filterValues {
				if value == filterValue {
					return true
				}
			}
		}
	}
	return false
}

package set

// Set represents a set of unique elements
type Set[T comparable] map[T]struct{}

// New returns a new set
func New[T comparable](ss ...T) Set[T] {
	return make(Set[T]).Add(ss...)
}

// Has returns true if an element is in a set
func (s Set[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

// Add adds a list of elements to a set
func (s Set[T]) Add(ss ...T) Set[T] {
	for _, e := range ss {
		s[e] = struct{}{}
	}
	return s
}

// Remove removes a list of elements from a set
func (s Set[T]) Remove(ss ...T) Set[T] {
	for _, e := range ss {
		delete(s, e)
	}
	return s
}

// Join joins two sets
func (s Set[T]) Join(ss Set[T]) Set[T] {
	for k := range ss {
		s[k] = struct{}{}
	}
	return s
}

// Copy returns a copy of a set
func (s Set[T]) Copy() Set[T] {
	return New[T]().Join(s)
}

// Slice returns the set as a slice
func (s Set[T]) Slice() []T {
	slice := []T{}
	for k := range s {
		slice = append(slice, k)
	}
	return slice
}

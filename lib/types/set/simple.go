package set

// SimpleSet represents a set of unique elements
type SimpleSet[T comparable] map[T]struct{}

// ensures SimpleSet implements Set at compile time
var _ Set[interface{}] = (SimpleSet[interface{}])(nil)

// New returns a new simple set
func New[T comparable](ss ...T) SimpleSet[T] {
	set := make(SimpleSet[T])
	set.Add(ss...)
	return set
}

// Has returns true if an element is in a set
func (s SimpleSet[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

// Add adds a list of elements to a set
func (s SimpleSet[T]) Add(ss ...T) Set[T] {
	for _, e := range ss {
		s[e] = struct{}{}
	}
	return s
}

// Remove removes a list of elements from a set
func (s SimpleSet[T]) Remove(ss ...T) Set[T] {
	for _, e := range ss {
		delete(s, e)
	}
	return s
}

// Join joins two sets
func (s SimpleSet[T]) Join(ss Set[T]) Set[T] {
	if ss != nil {
		for _, k := range ss.Slice() {
			s[k] = struct{}{}
		}
	}
	return s
}

// Copy returns a copy of a set
func (s SimpleSet[T]) Copy() Set[T] {
	return New[T]().Join(s)
}

// Slice returns the set as a slice
func (s SimpleSet[T]) Slice() []T {
	slice := []T{}
	for k := range s {
		slice = append(slice, k)
	}
	return slice
}

// Size returns the number of elements in the set
func (s SimpleSet[T]) Size() int {
	return len(s)
}

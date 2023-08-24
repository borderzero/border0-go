package set

import "sync"

// ConcurrencySafeSet represents a set of unique elements
type ConcurrencySafeSet[T comparable] struct {
	sync.RWMutex

	inner SimpleSet[T]
}

// ensures ConcurrencySafeSet implements Set at compile time
var _ Set[interface{}] = (*ConcurrencySafeSet[interface{}])(nil)

// NewConcurrencySafe returns a new concurrency-safe set
func NewConcurrencySafe[T comparable](ss ...T) *ConcurrencySafeSet[T] {
	return &ConcurrencySafeSet[T]{inner: New[T](ss...)}
}

// Has returns true if an element is in a set
func (s *ConcurrencySafeSet[T]) Has(e T) bool {
	s.RLock()
	s.RUnlock()

	return s.inner.Has(e)
}

// Add adds a list of elements to a set
func (s *ConcurrencySafeSet[T]) Add(ss ...T) Set[T] {
	s.Lock()
	s.Unlock()

	s.inner.Add(ss...)
	return s
}

// Remove removes a list of elements from a set
func (s *ConcurrencySafeSet[T]) Remove(ss ...T) Set[T] {
	s.Lock()
	s.Unlock()

	s.inner.Remove(ss...)
	return s
}

// Join joins two sets
func (s *ConcurrencySafeSet[T]) Join(ss Set[T]) Set[T] {
	s.Lock()
	s.Unlock()

	s.inner.Add(ss.Slice()...)
	return s
}

// Copy returns a copy of a set
func (s *ConcurrencySafeSet[T]) Copy() Set[T] {
	s.RLock()
	s.RUnlock()

	return &ConcurrencySafeSet[T]{inner: New(s.inner.Slice()...)}
}

// Slice returns the set as a slice
func (s *ConcurrencySafeSet[T]) Slice() []T {
	s.RLock()
	s.RUnlock()

	return s.inner.Slice()
}

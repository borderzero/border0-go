package set

import (
	"iter"
	"sync"
)

// ConcurrencySafeSet represents a set of unique elements.
type ConcurrencySafeSet[T comparable] struct {
	sync.RWMutex

	inner SimpleSet[T]
}

// ensures ConcurrencySafeSet implements Set at compile time.
var _ Set[any] = (*ConcurrencySafeSet[any])(nil)

// NewConcurrencySafe returns a new concurrency-safe set.
func NewConcurrencySafe[T comparable](ss ...T) *ConcurrencySafeSet[T] {
	return &ConcurrencySafeSet[T]{inner: New(ss...)}
}

// Has returns true if an element is in a set.
func (s *ConcurrencySafeSet[T]) Has(e T) bool {
	s.RLock()
	defer s.RUnlock()

	return s.inner.Has(e)
}

// Add adds a list of elements to a set.
func (s *ConcurrencySafeSet[T]) Add(ss ...T) Set[T] {
	s.Lock()
	defer s.Unlock()

	s.inner.Add(ss...)
	return s
}

// Remove removes a list of elements from a set.
func (s *ConcurrencySafeSet[T]) Remove(ss ...T) Set[T] {
	s.Lock()
	defer s.Unlock()

	s.inner.Remove(ss...)
	return s
}

// Join joins two sets.
func (s *ConcurrencySafeSet[T]) Join(ss Set[T]) Set[T] {
	// NOTE: this is done before acquiring the lock
	// to avoid the case where the set is being joined
	// with itself (would cause deadlock otherwise)
	elemsToJoin := ss.Slice()

	s.Lock()
	defer s.Unlock()

	s.inner.Add(elemsToJoin...)
	return s
}

// Iter returns an iterator for a copy of the set. Note that
// elements are iterated through in no particular order.
func (s *ConcurrencySafeSet[T]) Iter() iter.Seq[T] {
	s.RLock()
	defer s.RUnlock()

	return New(s.inner.Slice()...).Iter()
}

// Copy returns a copy of a set.
func (s *ConcurrencySafeSet[T]) Copy() Set[T] {
	s.RLock()
	defer s.RUnlock()

	return &ConcurrencySafeSet[T]{inner: New(s.inner.Slice()...)}
}

// Slice returns the set as a slice.
func (s *ConcurrencySafeSet[T]) Slice() []T {
	if s == nil || s.inner == nil {
		return []T{}
	}

	s.RLock()
	defer s.RUnlock()

	return s.inner.Slice()
}

// Size returns the number of elements in the set.
func (s *ConcurrencySafeSet[T]) Size() int {
	if s == nil || s.inner == nil {
		return 0
	}

	s.RLock()
	defer s.RUnlock()

	return s.inner.Size()
}

// Equals returns true if a given set is equal (has the same elements).
func (s *ConcurrencySafeSet[T]) Equals(comp Set[T]) bool {
	// must check if both sets reference the same
	// object, otherwise deadlock is possible.
	if concurrencySafeSetB, sameType := comp.(*ConcurrencySafeSet[T]); sameType {
		if s == concurrencySafeSetB {
			return true
		}
	}

	s.Lock()
	defer s.Unlock()

	return s.inner.Equals(comp)
}

// NewOfSameType returns a new empty set of the same type as the parent set.
func (s *ConcurrencySafeSet[T]) NewOfSameType(ss ...T) Set[T] {
	return NewConcurrencySafe(ss...)
}

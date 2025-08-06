package set

import (
	"iter"
)

// Set represents a set of unique elements.
type Set[T comparable] interface {
	Has(T) bool
	Add(...T) Set[T]
	Remove(...T) Set[T]
	Join(Set[T]) Set[T]
	Iter() iter.Seq[T]
	Copy() Set[T]
	Slice() []T
	Size() int
	Equals(Set[T]) bool
	NewOfSameType(...T) Set[T]
}

// Union returns elements in either a or b.
func Union[T comparable](a, b Set[T]) Set[T] {
	return a.Copy().Join(b)
}

// Intersection returns elements in either a or b.
func Intersection[T comparable](a, b Set[T]) Set[T] {
	var base, other Set[T]
	if a.Size() <= b.Size() {
		base, other = a, b
	} else {
		base, other = b, a
	}

	result := base.NewOfSameType()
	for v := range base.Iter() {
		if other.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns elements in a that are not in b.
func Difference[T comparable](a, b Set[T]) Set[T] {
	result := a.NewOfSameType()
	for v := range a.Iter() {
		if !b.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// SymmetricDifference returns elements in either a or b, but not in both.
func SymmetricDifference[T comparable](a, b Set[T]) Set[T] {
	return Union(Difference(a, b), Difference(b, a))
}

// Complement returns the relative complement of b in a
// i.e. all the elements that are in a but not in b.
func Complement[T comparable](a, b Set[T]) Set[T] {
	result := a.Copy()
	for v := range b.Iter() {
		result.Remove(v)
	}
	return result
}

// AddedOrRemovedSlice returns the items that were added or removed (symmetric difference)
// between the before and after sets, as a slice. It avoids extra allocations.
func AddedOrRemovedSlice[T comparable](before, after Set[T]) []T {
	// preallocate slice with upper bound of size
	result := make([]T, 0, before.Size()+after.Size())

	for v := range after.Iter() {
		if !before.Has(v) {
			result = append(result, v)
		}
	}
	for v := range before.Iter() {
		if !after.Has(v) {
			result = append(result, v)
		}
	}
	return result
}

// AddedOrRemoved returns the set of items that were added or removed between the before and after sets.
// This is equivalent to the symmetric difference: (after ∖ before) ∪ (before ∖ after).
func AddedOrRemoved[T comparable](before, after Set[T]) Set[T] {
	return before.NewOfSameType(AddedOrRemovedSlice(before, after)...)
}

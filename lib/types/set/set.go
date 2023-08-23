package set

// Set represents a set of unique elements
type Set[T comparable] interface {
	Has(T) bool
	Add(...T) Set[T]
	Remove(...T) Set[T]
	Join(Set[T]) Set[T]
	Copy() Set[T]
	Slice() []T
}

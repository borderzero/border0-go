package pointer

// To returns the address of any given value.
// This is useful because in Go it is illegal to take the address of
// a literal value (i.e. is not possible to do &"string") - this
// helper function allows us to do pointer.To("string") in-line.
func To[T any](value T) *T {
	return &value
}

// ValueOrZero returns the value of dereferencing a pointer if
// the pointer is not nil, or the zero value of the type if nil.
func ValueOrZero[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	return *new(T)
}

// ValueOrDefault returns the value of dereferencing a pointer
// if the pointer is not nil, or a given default value if nil.
func ValueOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

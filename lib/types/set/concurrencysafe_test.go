package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConcurrencySafe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		elems      []string
		expectSize int
	}{
		{
			name:       "New with nil argument should return empty set",
			elems:      nil,
			expectSize: 0,
		},
		{
			name:       "New with empty slice argument should return empty set",
			elems:      []string{},
			expectSize: 0,
		},
		{
			name:       "New with one element should return set with one element",
			elems:      []string{"a"},
			expectSize: 1,
		},
		{
			name:       "New with two unique elements should return set with two elements",
			elems:      []string{"a", "b"},
			expectSize: 2,
		},
		{
			name:       "New with multiple non unique elements should return set of unique elements",
			elems:      []string{"a", "b", "c", "a", "a", "b", "c"},
			expectSize: 3, // a, b, c
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := NewConcurrencySafe(test.elems...)

			assert.NotNil(t, s)
			assert.Equal(t, s.Size(), test.expectSize)
			for _, elem := range test.elems {
				assert.True(t, s.Has(elem))
			}
		})
	}
}

func Test_ConcurrencySafeSetHas(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"
	mockElemC := "c"
	mockSet := NewConcurrencySafe(mockElemA, mockElemB)

	tests := []struct {
		name          string
		set           *ConcurrencySafeSet[string]
		elem          string
		shouldBeFound bool
	}{
		{
			name:          "Should return true when element is in set",
			set:           mockSet,
			elem:          mockElemA,
			shouldBeFound: true,
		},
		{
			name:          "Should return false when element is not in set",
			set:           mockSet,
			elem:          mockElemC,
			shouldBeFound: false,
		},
		{
			name:          "Should return false when set is empty",
			set:           NewConcurrencySafe[string](),
			elem:          mockElemA,
			shouldBeFound: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.set.Has(test.elem), test.shouldBeFound)
		})
	}
}

func Test_ConcurrencySafeSetAdd(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name   string
		before *ConcurrencySafeSet[string]
		toAdd  []string
		after  *ConcurrencySafeSet[string]
	}{
		{
			name:   "Should not do anything when Add with nil arguments",
			before: NewConcurrencySafe[string](),
			toAdd:  nil,
			after:  NewConcurrencySafe[string](),
		},
		{
			name:   "Should not do anything when Add with no arguments",
			before: NewConcurrencySafe[string](),
			toAdd:  []string{},
			after:  NewConcurrencySafe[string](),
		},
		{
			name:   "Should add when Add has a single argument",
			before: NewConcurrencySafe[string](),
			toAdd:  []string{mockElemA},
			after:  NewConcurrencySafe[string](mockElemA),
		},
		{
			name:   "Should add when Add has multiple unique argument",
			before: NewConcurrencySafe[string](),
			toAdd:  []string{mockElemA, mockElemB},
			after:  NewConcurrencySafe[string](mockElemA, mockElemB),
		},
		{
			name:   "Should add only unique elements when Add has repeated arguments",
			before: NewConcurrencySafe[string](),
			toAdd:  []string{mockElemA, mockElemA, mockElemA},
			after:  NewConcurrencySafe[string](mockElemA),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.before.Add(test.toAdd...), test.after)
		})
	}
}

func Test_ConcurrencySafeSetRemove(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name     string
		before   *ConcurrencySafeSet[string]
		toRemove []string
		after    *ConcurrencySafeSet[string]
	}{
		{
			name:     "Should not do anything when Remove nil",
			before:   NewConcurrencySafe[string](),
			toRemove: nil,
			after:    NewConcurrencySafe[string](),
		},
		{
			name:     "Should not do anything when Remove empty",
			before:   NewConcurrencySafe[string](),
			toRemove: []string{},
			after:    NewConcurrencySafe[string](),
		},
		{
			name:     "Should not do anything when Remove element already not in set",
			before:   NewConcurrencySafe[string](mockElemB),
			toRemove: []string{mockElemA},
			after:    NewConcurrencySafe[string](mockElemB),
		},
		{
			name:     "Should remove element when element in set",
			before:   NewConcurrencySafe[string](mockElemA, mockElemB),
			toRemove: []string{mockElemA},
			after:    NewConcurrencySafe[string](mockElemB),
		},
		{
			name:     "Should remove multiple at once",
			before:   NewConcurrencySafe[string](mockElemA, mockElemB),
			toRemove: []string{mockElemA, mockElemB},
			after:    NewConcurrencySafe[string](),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.before.Remove(test.toRemove...), test.after)
		})
	}
}

func Test_ConcurrencySafeSetJoin(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	mockSet := NewConcurrencySafe[string](mockElemA, mockElemB)

	tests := []struct {
		name   string
		before *ConcurrencySafeSet[string]
		toJoin *ConcurrencySafeSet[string]
		after  *ConcurrencySafeSet[string]
	}{
		{
			name:   "Should not do anything when joining nil",
			before: NewConcurrencySafe[string](),
			toJoin: nil,
			after:  NewConcurrencySafe[string](),
		},
		{
			name:   "Should not do anything when joining empty set",
			before: NewConcurrencySafe[string](),
			toJoin: NewConcurrencySafe[string](),
			after:  NewConcurrencySafe[string](),
		},
		{
			name:   "Should join set with single element in set B",
			before: NewConcurrencySafe[string](),
			toJoin: NewConcurrencySafe[string](mockElemA),
			after:  NewConcurrencySafe[string](mockElemA),
		},
		{
			name:   "Should join set with single element in set A",
			before: NewConcurrencySafe[string](mockElemA),
			toJoin: NewConcurrencySafe[string](),
			after:  NewConcurrencySafe[string](mockElemA),
		},
		{
			name:   "Should join set with single repeated element",
			before: NewConcurrencySafe[string](mockElemA),
			toJoin: NewConcurrencySafe[string](mockElemA),
			after:  NewConcurrencySafe[string](mockElemA),
		},
		{
			name:   "Should join set with repeated elements",
			before: NewConcurrencySafe[string](mockElemA, mockElemB),
			toJoin: NewConcurrencySafe[string](mockElemA, mockElemB),
			after:  NewConcurrencySafe[string](mockElemA, mockElemB),
		},
		{
			name:   "Should join set with unique elements",
			before: NewConcurrencySafe[string](mockElemA),
			toJoin: NewConcurrencySafe[string](mockElemB),
			after:  NewConcurrencySafe[string](mockElemA, mockElemB),
		},
		{
			name:   "Should join self without deadlock",
			before: mockSet,
			toJoin: mockSet,
			after:  mockSet,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.before.Join(test.toJoin), test.after)
		})
	}
}

func Test_ConcurrencySafeSetCopy(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name   string
		toCopy *ConcurrencySafeSet[string]
		copy   *ConcurrencySafeSet[string]
	}{
		{
			name:   "Should copy empty set",
			toCopy: NewConcurrencySafe[string](),
			copy:   NewConcurrencySafe[string](),
		},
		{
			name:   "Should copy non empty set, single element",
			toCopy: NewConcurrencySafe[string](mockElemA),
			copy:   NewConcurrencySafe[string](mockElemA),
		},
		{
			name:   "Should copy non empty set, multiple element",
			toCopy: NewConcurrencySafe[string](mockElemA, mockElemB),
			copy:   NewConcurrencySafe[string](mockElemA, mockElemB),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.toCopy.Copy(), test.copy)
		})
	}
}

func Test_ConcurrencySafeSetSlice(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name  string
		set   *ConcurrencySafeSet[string]
		slice []string
	}{
		{
			name:  "Empty set should yield empty slice",
			set:   NewConcurrencySafe[string](),
			slice: []string{},
		},
		{
			name:  "Non empty set should yield non empty slice (single element)",
			set:   NewConcurrencySafe[string](mockElemA),
			slice: []string{mockElemA},
		},
		{
			name:  "Non empty set should yield non empty slice (multiple elements)",
			set:   NewConcurrencySafe[string](mockElemA, mockElemB),
			slice: []string{mockElemA, mockElemB},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, test.set.Slice(), test.slice)
		})
	}
}

func Test_ConcurrencySafeSetSize(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name string
		set  *ConcurrencySafeSet[string]
		size int
	}{
		{
			name: "Empty set should yield size 0",
			set:  NewConcurrencySafe[string](),
			size: 0,
		},
		{
			name: "Non empty set should yield size 1 (single element)",
			set:  NewConcurrencySafe[string](mockElemA),
			size: 1,
		},
		{
			name: "Non empty set should yield size 2 (multiple elements)",
			set:  NewConcurrencySafe[string](mockElemA, mockElemB),
			size: 2,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.set.Size(), test.size)
		})
	}
}

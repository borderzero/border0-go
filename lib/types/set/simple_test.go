package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
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

			s := New(test.elems...)

			assert.NotNil(t, s)
			assert.Equal(t, len(s), test.expectSize)
			for _, elem := range test.elems {
				assert.Contains(t, s, elem)
			}
		})
	}
}

func Test_SimpleSetHas(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"
	mockElemC := "c"
	mockSet := New(mockElemA, mockElemB)

	tests := []struct {
		name          string
		set           SimpleSet[string]
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
			set:           New[string](),
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

func Test_SimpleSetAdd(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name   string
		before SimpleSet[string]
		toAdd  []string
		after  SimpleSet[string]
	}{
		{
			name:   "Should not do anything when Add with nil arguments",
			before: New[string](),
			toAdd:  nil,
			after:  New[string](),
		},
		{
			name:   "Should not do anything when Add with no arguments",
			before: New[string](),
			toAdd:  []string{},
			after:  New[string](),
		},
		{
			name:   "Should add when Add has a single argument",
			before: New[string](),
			toAdd:  []string{mockElemA},
			after:  New[string](mockElemA),
		},
		{
			name:   "Should add when Add has multiple unique argument",
			before: New[string](),
			toAdd:  []string{mockElemA, mockElemB},
			after:  New[string](mockElemA, mockElemB),
		},
		{
			name:   "Should add only unique elements when Add has repeated arguments",
			before: New[string](),
			toAdd:  []string{mockElemA, mockElemA, mockElemA},
			after:  New[string](mockElemA),
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

func Test_SimpleSetRemove(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name     string
		before   SimpleSet[string]
		toRemove []string
		after    SimpleSet[string]
	}{
		{
			name:     "Should not do anything when Remove nil",
			before:   New[string](),
			toRemove: nil,
			after:    New[string](),
		},
		{
			name:     "Should not do anything when Remove empty",
			before:   New[string](),
			toRemove: []string{},
			after:    New[string](),
		},
		{
			name:     "Should not do anything when Remove element already not in set",
			before:   New[string](mockElemB),
			toRemove: []string{mockElemA},
			after:    New[string](mockElemB),
		},
		{
			name:     "Should remove element when element in set",
			before:   New[string](mockElemA, mockElemB),
			toRemove: []string{mockElemA},
			after:    New[string](mockElemB),
		},
		{
			name:     "Should remove multiple at once",
			before:   New[string](mockElemA, mockElemB),
			toRemove: []string{mockElemA, mockElemB},
			after:    New[string](),
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

func Test_SimpleSetJoin(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name   string
		before SimpleSet[string]
		toJoin SimpleSet[string]
		after  SimpleSet[string]
	}{
		{
			name:   "Should not do anything when joining nil",
			before: New[string](),
			toJoin: nil,
			after:  New[string](),
		},
		{
			name:   "Should not do anything when joining empty set",
			before: New[string](),
			toJoin: New[string](),
			after:  New[string](),
		},
		{
			name:   "Should join set with single element in set B",
			before: New[string](),
			toJoin: New[string](mockElemA),
			after:  New[string](mockElemA),
		},
		{
			name:   "Should join set with single element in set A",
			before: New[string](mockElemA),
			toJoin: New[string](),
			after:  New[string](mockElemA),
		},
		{
			name:   "Should join set with single repeated element",
			before: New[string](mockElemA),
			toJoin: New[string](mockElemA),
			after:  New[string](mockElemA),
		},
		{
			name:   "Should join set with repeated elements",
			before: New[string](mockElemA, mockElemB),
			toJoin: New[string](mockElemA, mockElemB),
			after:  New[string](mockElemA, mockElemB),
		},
		{
			name:   "Should join set with unique elements",
			before: New[string](mockElemA),
			toJoin: New[string](mockElemB),
			after:  New[string](mockElemA, mockElemB),
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

func Test_SimpleSetCopy(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name   string
		toCopy SimpleSet[string]
		copy   SimpleSet[string]
	}{
		{
			name:   "Should copy empty set",
			toCopy: New[string](),
			copy:   New[string](),
		},
		{
			name:   "Should copy non empty set, single element",
			toCopy: New[string](mockElemA),
			copy:   New[string](mockElemA),
		},
		{
			name:   "Should copy non empty set, multiple element",
			toCopy: New[string](mockElemA, mockElemB),
			copy:   New[string](mockElemA, mockElemB),
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

func Test_SimpleSetSlice(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name  string
		set   SimpleSet[string]
		slice []string
	}{
		{
			name:  "Empty set should yield empty slice",
			set:   New[string](),
			slice: []string{},
		},
		{
			name:  "Non empty set should yield non empty slice (single element)",
			set:   New[string](mockElemA),
			slice: []string{mockElemA},
		},
		{
			name:  "Non empty set should yield non empty slice (multiple elements)",
			set:   New[string](mockElemA, mockElemB),
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

func Test_SimpleSetSize(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name string
		set  SimpleSet[string]
		size int
	}{
		{
			name: "Empty set should yield size 0",
			set:  New[string](),
			size: 0,
		},
		{
			name: "Non empty set should yield size 1 (single element)",
			set:  New[string](mockElemA),
			size: 1,
		},
		{
			name: "Non empty set should yield size 2 (multiple elements)",
			set:  New[string](mockElemA, mockElemB),
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

func Test_SimpleSetEquals(t *testing.T) {
	t.Parallel()

	mockElemA := "a"
	mockElemB := "b"

	tests := []struct {
		name string
		seta SimpleSet[string]
		setb SimpleSet[string]
		eq   bool
	}{
		{
			name: "Empty sets should be equal",
			seta: New[string](),
			setb: New[string](),
			eq:   true,
		},
		{
			name: "Sets with equal elements should be equal",
			seta: New[string](mockElemA),
			setb: New[string](mockElemA),
			eq:   true,
		},
		{
			name: "Sets with different elements should not be equal",
			seta: New[string](mockElemA),
			setb: New[string](mockElemB),
			eq:   false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.seta.Equals(test.setb), test.eq)
		})
	}
}

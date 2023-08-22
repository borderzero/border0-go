package nilcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AreAllNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []any
		expect bool
	}{
		{
			name:   "All with no values should return true",
			values: []any{},
			expect: true,
		},
		{
			name:   "All with single nil value should return true",
			values: []any{nil},
			expect: true,
		},
		{
			name:   "All with multiple nil value should return true",
			values: []any{nil, nil},
			expect: true,
		},
		{
			name:   "All with single nil struct value should return true",
			values: []any{(*struct{ Field string })(nil)},
			expect: true,
		},
		{
			name:   "All with single non nil struct value should return false",
			values: []any{struct{ Field string }{Field: "hello"}},
			expect: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expect, AreAllNil(test.values...))
		})
	}
}

func Test_AnyNotNil(t *testing.T) {
	t.Parallel()

	type (
		test1 struct{ field string }
		test2 struct{ field string }
	)

	var (
		nilStruct1 *test1
		nilStruct2 *test2

		nonNilStruct1 = new(test1)
		nonNilStruct2 = new(test2)
	)

	tests := []struct {
		name   string
		values []any
		expect bool
	}{
		{
			name:   "no values",
			values: []any{},
			expect: false,
		},
		{
			name:   "all nils",
			values: []any{nil, nilStruct1, nilStruct2},
			expect: false,
		},
		{
			name:   "mixed nils and non nils",
			values: []any{"hello", nonNilStruct1, nilStruct2},
			expect: true,
		},
		{
			name:   "all non nils with mixed types",
			values: []any{"hello", nonNilStruct1, nonNilStruct2, 123, []byte("hello"), true, false},
			expect: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expect, AnyNotNil(test.values...))
		})
	}
}

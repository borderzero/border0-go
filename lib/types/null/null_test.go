package null

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
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
			assert.Equal(t, test.expect, All(test.values...))
		})
	}
}

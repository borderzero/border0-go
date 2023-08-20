package slice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Diff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		givenOriginal []any
		givenChanged  []any
		wantNew       []any
		wantRemoved   []any
	}{
		{
			name:          "no change",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  []any{"111", "222", "333"},
			wantNew:       nil,
			wantRemoved:   nil,
		},
		{
			name:          "add items",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  []any{"111", "222", "333", "444", "555"},
			wantNew:       []any{"444", "555"},
			wantRemoved:   nil,
		},
		{
			name:          "remove items",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  []any{"111"},
			wantNew:       nil,
			wantRemoved:   []any{"222", "333"},
		},
		{
			name:          "add and remove items",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  []any{"111", "444", "555"},
			wantNew:       []any{"444", "555"},
			wantRemoved:   []any{"222", "333"},
		},
		{
			name:          "add and remove items in different order",
			givenOriginal: []any{"222", "111", "333"},
			givenChanged:  []any{"444", "555", "111"},
			wantNew:       []any{"444", "555"},
			wantRemoved:   []any{"222", "333"},
		},
		{
			name:          "add and remove items with int type",
			givenOriginal: []any{1, 2, 3},
			givenChanged:  []any{1, 4, 5},
			wantNew:       []any{4, 5},
			wantRemoved:   []any{2, 3},
		},
		{
			name:          "add and remove items with time.Time type",
			givenOriginal: []any{time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC), time.Date(2022, 2, 2, 2, 2, 2, 2, time.UTC)},
			givenChanged:  []any{time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC), time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)},
			wantNew:       []any{time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)},
			wantRemoved:   []any{time.Date(2022, 2, 2, 2, 2, 2, 2, time.UTC)},
		},
		{
			name:          "remove everything",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  nil,
			wantNew:       nil,
			wantRemoved:   []any{"111", "222", "333"},
		},
		{
			name:          "was empty and add new items",
			givenOriginal: nil,
			givenChanged:  []any{"111", "222", "333"},
			wantNew:       []any{"111", "222", "333"},
			wantRemoved:   nil,
		},
		{
			name:          "remove everything and add new",
			givenOriginal: []any{"111", "222", "333"},
			givenChanged:  []any{"444", "555", "666"},
			wantNew:       []any{"444", "555", "666"},
			wantRemoved:   []any{"111", "222", "333"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotNew, gotRemoved := Diff(test.givenOriginal, test.givenChanged)

			assert.EqualValues(t, test.wantNew, gotNew)
			assert.EqualValues(t, test.wantRemoved, gotRemoved)
		})
	}
}

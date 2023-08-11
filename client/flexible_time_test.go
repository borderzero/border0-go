package client

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FlexibleTimeFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		given    string
		wantTime FlexibleTime
		wantErr  error
	}{
		{
			name:    "failed to parse input date and time string",
			given:   "not a valid date and time",
			wantErr: errors.New(`parsing time "not a valid date and time" as "2006-01-02T15:04:05Z07:00": cannot parse "not a valid date and time" as "2006"`),
		},
		{
			name:     "happy path",
			given:    "2020-01-02T15:04:05Z",
			wantTime: FlexibleTime{time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotTime, gotErr := FlexibleTimeFrom(test.given)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantTime, gotTime)
		})
	}
}

func Test_FlexibleTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given FlexibleTime
		want  []byte
	}{
		{
			name:  "inner time has zero value",
			given: FlexibleTime{time.Time{}},
			want:  []byte("0"),
		},
		{
			name:  "inner time has non-zero value",
			given: FlexibleTime{time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)},
			want:  []byte("1577977445"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.given.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_FlexibleTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		given    []byte
		wantTime time.Time
		wantErr  error
	}{
		// test string value input
		{
			name:    "not a valid string json value",
			given:   []byte(`"`),
			wantErr: errors.New("unexpected end of JSON input"),
		},
		{
			name:    "string value is valid json but not valid time",
			given:   []byte(`"not a valid date and time"`),
			wantErr: errors.New(`parsing time "not a valid date and time" as "2006-01-02T15:04:05Z07:00": cannot parse "not a valid date and time" as "2006"`),
		},
		{
			name:     "happy path - string time value",
			given:    []byte(`"2020-01-02T15:04:05Z"`),
			wantTime: time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		// test integer value input
		{
			name:    "not a valid integer json value",
			given:   []byte("0101"),
			wantErr: errors.New("invalid character '1' after top-level value"),
		},
		{
			name:     "happy path - integer time value",
			given:    []byte("1577977445"),
			wantTime: time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var got FlexibleTime
			gotErr := got.UnmarshalJSON(test.given)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantTime, got.Time.UTC())
		})
	}
}

func Test_FlexibleTime_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given FlexibleTime
		want  string
	}{
		{
			name:  "inner time has zero value",
			given: FlexibleTime{time.Time{}},
			want:  "",
		},
		{
			name:  "inner time has non-zero value",
			given: FlexibleTime{time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)},
			want:  "2020-01-02T15:04:05Z",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.given.String()
			assert.Equal(t, test.want, got)
		})
	}
}

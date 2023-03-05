package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ExponentialBackoff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		givenMin     time.Duration
		givenMax     time.Duration
		givenAttempt int
		want         time.Duration
	}{
		{
			name:         "first attempt",
			givenMin:     time.Second,
			givenMax:     10 * time.Second,
			givenAttempt: 0,
			want:         time.Second, // it should not be the min
		},
		{
			name:         "third attempt",
			givenMin:     time.Second,
			givenMax:     10 * time.Second,
			givenAttempt: 2,
			want:         4 * time.Second, // it should be the min * 2^attempt
		},
		{
			name:         "Large Attempt",
			givenMin:     time.Second,
			givenMax:     10 * time.Second,
			givenAttempt: 10,
			want:         10 * time.Second, // it should not be more than max
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			backoff := ExponentialBackoff(test.givenMin, test.givenMax, test.givenAttempt)
			assert.Equal(t, test.want, backoff)
		})
	}
}

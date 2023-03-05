package client

import (
	"math"
	"time"
)

// Backoff is a callback function which will be called by APIClient when
// performing retries. It is passed the minimum and maximum durations to
// backoff between, as well as the attempt number (starting at zero)
type Backoff func(min, max time.Duration, attempt int) time.Duration

// ExponentialBackoff is a Backoff function which will backoff exponentially
// between the given minimum and maximum durations. The attempt number is used
// as the exponent base, so the first attempt will backoff by the minimum
// duration, the second attempt will backoff by twice the minimum duration, the
// third attempt will backoff by four times the minimum duration, and so on.
// The maximum duration is used as a cap, so the backoff will never exceed the
// maximum duration.
func ExponentialBackoff(min, max time.Duration, attempt int) time.Duration {
	multiple := math.Pow(2, float64(attempt)) * float64(min)
	backoff := time.Duration(multiple)
	if float64(backoff) != multiple || backoff > max {
		backoff = max
	}
	return backoff
}

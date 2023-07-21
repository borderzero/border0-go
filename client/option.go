package client

import (
	"time"
)

// Option is a function that can be passed to NewAPIClient to configure it.
type Option func(*APIClient)

// WithTimeout sets the timeout for the underlying http client.
func WithTimeout(timeout time.Duration) Option {
	return func(api *APIClient) {
		api.timeout = timeout
	}
}

// WithAuthToken sets the auth token for Border0 api calls.
func WithAuthToken(token string) Option {
	return func(api *APIClient) {
		api.authToken = token
	}
}

// WithBaseURL sets the base url for Border0 api calls.
func WithBaseURL(url string) Option {
	return func(api *APIClient) {
		api.baseURL = url
	}
}

// WithRetryWaitMin sets the minimum wait time between retries of failed api calls.
func WithRetryWaitMin(wait time.Duration) Option {
	return func(api *APIClient) {
		api.retryWaitMin = wait
	}
}

// WithRetryWaitMax sets the maximum wait time between retries of failed api calls.
func WithRetryWaitMax(wait time.Duration) Option {
	return func(api *APIClient) {
		api.retryWaitMax = wait
	}
}

// WithRetryMax sets the maximum number of retries of failed api calls.
func WithRetryMax(attempts int) Option {
	return func(api *APIClient) {
		api.retryMax = attempts
	}
}

// WithBackoff sets the backoff function that's used to calculate the wait time
// between retries of failed api calls.
func WithBackoff(fn Backoff) Option {
	return func(api *APIClient) {
		api.backoff = fn
	}
}

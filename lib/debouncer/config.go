package debouncer

import "time"

const (
	defaultDebounceTime = time.Second
	defaultMaxWaitTime  = time.Duration(-1)
)

type config struct {
	debounceTime time.Duration
	maxWaitTime  time.Duration
}

func newDefaultConfig() *config {
	return &config{
		debounceTime: defaultDebounceTime,
		maxWaitTime:  defaultMaxWaitTime,
	}
}

func (c *config) Apply(opts ...Option) *config {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

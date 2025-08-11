package debouncer

import "time"

type Option func(*config)

// WithDebounceTime sets the debounce time i.e. how long to wait
// without further invocations before firing. Defaults to 1 second.
func WithDebounceTime(d time.Duration) Option { return func(c *config) { c.debounceTime = d } }

// WithMaxWaitTime sets the max wait time i.e. maximum time spent deferring invocations
// due to new invocations. Defaults to no maximum i.e. never fire for infinite successive
// invocations within the debounce time.
func WithMaxWaitTime(d time.Duration) Option { return func(c *config) { c.maxWaitTime = d } }

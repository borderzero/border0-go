package debouncer

import (
	"sync"
	"time"
)

// Debouncer represents an entity
// capable of debouncing function calls.
type Debouncer[T any] struct {
	debounceTime   time.Duration
	debounceChan   chan T
	mu             sync.Mutex
	timer          *time.Timer
	pendingCallArg T
	fn             func(T)
}

// NewDebouncer is the debouncer constructor.
func NewDebouncer[T any](debounceTime time.Duration, fn func(T)) chan<- T {
	debounceChan := make(chan T, 1)
	debouncer := &Debouncer[T]{
		debounceTime: debounceTime,
		debounceChan: debounceChan,
		fn:           fn,
	}

	go debouncer.start()
	return debounceChan
}

// Start starts the debouncer.
func (d *Debouncer[T]) start() {
	for msg := range d.debounceChan {
		d.mu.Lock()

		// Update the pending argument.
		d.pendingCallArg = msg

		// Stop and drain the existing timer, if needed.
		if d.timer != nil {
			if !d.timer.Stop() {
				select {
				case <-d.timer.C:
				default:
				}
			}
		}

		// Start a new timer.
		d.timer = time.AfterFunc(d.debounceTime, func() {
			d.mu.Lock()
			arg := d.pendingCallArg
			d.mu.Unlock()

			// Call the debounced function.
			d.fn(arg)
		})

		d.mu.Unlock()
	}
}

package debouncer

import (
	"sync"
	"time"
)

// Debouncer represents an entity
// capable of debouncing function calls.
type Debouncer[T any] interface {
	Do(T)
	Flush()
	Close()
}

// debouncer is the default Debouncer implementation.
type debouncer[T any] struct {
	debounceTime   time.Duration
	maxWaitTime    time.Duration
	debounceChan   chan T
	wg             sync.WaitGroup
	mu             sync.Mutex
	timer          *time.Timer
	maxWaitTimer   *time.Timer
	pendingCallArg T
	hasPending     bool
	fn             func(T)
}

// New is the Debouncer constructor.
func New[T any](debounceTime, maxWaitTime time.Duration, fn func(T)) Debouncer[T] {
	debounceChan := make(chan T, 1)
	d := &debouncer[T]{
		debounceTime: debounceTime,
		maxWaitTime:  maxWaitTime,
		debounceChan: debounceChan,
		fn:           fn,
	}
	d.wg.Add(1)
	go d.start()
	return d
}

func (d *debouncer[T]) start() {
	defer d.wg.Done()
	for msg := range d.debounceChan {
		func() {
			d.mu.Lock()
			defer d.mu.Unlock()

			d.pendingCallArg = msg
			d.hasPending = true

			// Reset debounce timer
			if d.timer != nil {
				if !d.timer.Stop() {
					select {
					case <-d.timer.C:
					default:
					}
				}
			}
			d.timer = time.AfterFunc(d.debounceTime, d.Flush)

			// Start or keep max wait timer
			if d.maxWaitTime > 0 && d.maxWaitTimer == nil {
				d.maxWaitTimer = time.AfterFunc(d.maxWaitTime, d.Flush)
			}
		}()
	}
}

// Do submits a new function request.
func (d *debouncer[T]) Do(t T) {
	select {
	case d.debounceChan <- t:
	default:
		// channel is full, drop the old message
		select {
		case <-d.debounceChan:
		default:
		}
		// now send the new one (guaranteed to succeed)
		d.debounceChan <- t
	}
}

// Flush flushes the debouncer immediately (calling the function) without debouncing.
func (d *debouncer[T]) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// stop debounce timer
	if d.timer != nil {
		if !d.timer.Stop() {
			select {
			case <-d.timer.C:
			default:
			}
		}
		d.timer = nil
	}

	// stop max wait timer
	if d.maxWaitTimer != nil {
		if !d.maxWaitTimer.Stop() {
			select {
			case <-d.maxWaitTimer.C:
			default:
			}
		}
		d.maxWaitTimer = nil
	}

	// call the function
	if d.hasPending {
		d.hasPending = false
		d.fn(d.pendingCallArg)
	}
}

// Close gracefully closes and flushes the debouncer.
func (d *debouncer[T]) Close() {
	close(d.debounceChan)
	d.wg.Wait()
	d.Flush()
}

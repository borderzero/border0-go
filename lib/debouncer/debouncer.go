package debouncer

import (
	"sync"
	"sync/atomic"
	"time"
)

// Debouncer represents an entity
// capable of debouncing function calls.
type Debouncer[T any] interface {
	Flushable

	Do(T) Flushable
	Close()
}

// Flushable represents an entity
// capable of being flushed.
type Flushable interface {
	Flush()
}

// debouncer is the default Debouncer implementation.
type debouncer[T any] struct {
	// timer for debouncing
	debounceTime time.Duration
	timer        *time.Timer

	// timer for preventing infinite debouncing
	maxWaitTime  time.Duration
	maxWaitTimer *time.Timer

	// argument for next invocation
	pendingCallArg atomic.Pointer[T]

	// lock for timers, current pending arg and pending indicator
	updateLock sync.Mutex

	// wait group for waiting for the debouncer to be flushed
	processorWG sync.WaitGroup

	// channel for triggering invocations
	invocationChan chan struct{}

	// lock for enforcing the arg policy during Do()
	sendLock sync.Mutex

	// user provided function
	fn func(T)

	// closed indicator
	closed atomic.Bool
}

// New is the Debouncer constructor.
func New[T any](fn func(T), opts ...Option) Debouncer[T] {
	return newWithOpts(fn, newDefaultConfig().Apply(opts...))
}

func newWithOpts[T any](fn func(T), cfg *config) Debouncer[T] {
	d := &debouncer[T]{
		debounceTime:   cfg.debounceTime,
		maxWaitTime:    cfg.maxWaitTime,
		invocationChan: make(chan struct{}, 1),
		fn:             fn,
	}
	d.pendingCallArg.Store(nil)
	d.processorWG.Add(1)
	go d.processor()
	return d
}

func (d *debouncer[T]) processor() {
	defer d.processorWG.Done()
	for range d.invocationChan {
		func() {
			d.updateLock.Lock()
			defer d.updateLock.Unlock()

			// reset debounce timer
			if d.timer != nil {
				if !d.timer.Stop() {
					select {
					case <-d.timer.C:
					default:
					}
				}
			}
			d.timer = time.AfterFunc(d.debounceTime, d.Flush)

			// start or keep max wait timer
			if d.maxWaitTime > 0 && d.maxWaitTimer == nil {
				d.maxWaitTimer = time.AfterFunc(d.maxWaitTime, d.Flush)
			}
		}()
	}
}

// Do submits a new function request.
func (d *debouncer[T]) Do(t T) Flushable {
	d.sendLock.Lock()
	defer d.sendLock.Unlock()

	if d.closed.Load() {
		return d
	}

	// update the arg
	d.pendingCallArg.Store(&t)

	select {
	case d.invocationChan <- struct{}{}:
	default:
	}
	return d
}

// Flush flushes the debouncer immediately (calling the function) without debouncing.
func (d *debouncer[T]) Flush() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

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

	// consume the pending arg exactly once i.e. second Flush() is a no-op
	if arg := d.pendingCallArg.Swap(nil); arg != nil {
		d.fn(*arg)
	}
}

// Close gracefully closes and flushes the debouncer.
func (d *debouncer[T]) Close() {
	if !d.closed.CompareAndSwap(false, true) {
		// already closed
		return
	}

	d.closeChannel()
	d.processorWG.Wait()
	d.Flush()
}

func (d *debouncer[T]) closeChannel() {
	d.sendLock.Lock()
	defer d.sendLock.Unlock()

	close(d.invocationChan)
}

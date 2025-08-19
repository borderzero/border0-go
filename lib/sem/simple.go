package sem

import "context"

// semaphore is a concrete implementation of
// the Semaphore interface using a buffered channel.
type semaphore struct {
	ch chan struct{}
}

// New returns a Semaphore that allows up to n concurrent holders.
func New(n uint) Semaphore {
	return &semaphore{
		// NOTE: This channel is intentionally never closed.
		// It will be garbage collected when there are no goroutines
		// blocked on it and no references to the semaphore object.
		ch: make(chan struct{}, n),
	}
}

// Acquire takes a slot in the semaphore, blocking if none are available.
func (s *semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot in the semaphore.
func (s *semaphore) Release() { <-s.ch }

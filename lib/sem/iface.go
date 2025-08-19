package sem

import "context"

// Semaphore is the interface that wraps basic Acquire and Release
// methods to control access to a limited resource pool.
type Semaphore interface {
	// Acquire blocks until a slot is available in the semaphore.
	Acquire(context.Context) error

	// Release frees a slot in the semaphore.
	Release()
}

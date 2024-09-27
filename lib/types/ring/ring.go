package ring

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// Number is an interface representing all numerical types in Go.
type Number interface {
	constraints.Integer | constraints.Float
}

// Ring is an entity capable of keeping only the a number of most recent values.
type Ring[N Number] interface {
	Put(N)            // Put adds a sample to the Ring.
	Average() float64 // Average computes the average value of all values in the Ring.
}

// ring is the default implementation of the Ring interface.
type ring[N Number] struct {
	mu      sync.RWMutex
	entries []N
	window  int
	puts    int
}

// New returns a new default implementation of the Ring interface.
func New[N Number](window int) Ring[N] {
	return &ring[N]{
		mu:      sync.RWMutex{},
		entries: make([]N, window),
		window:  window,
		puts:    0,
	}
}

// Put adds a sample to the Ring.
func (r *ring[N]) Put(n N) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries[r.puts%r.window] = n
	r.puts++
}

// Average computes the average value of all values in the Ring.
func (r *ring[N]) Average() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// account for the case of not having any
	// elements to avoid a division by zero.
	if r.puts == 0 {
		return 0
	}

	// account for the case of not having enough
	// elements to compute over the whole window.
	items := r.window
	if r.puts < items {
		items = r.puts
	}

	sum := N(0)
	for i := 0; i < items; i++ {
		sum += r.entries[i]
	}
	return float64(sum) / float64(items)
}

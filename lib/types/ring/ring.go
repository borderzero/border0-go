package ring

import (
	"math"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// Number is an interface representing all numerical types in Go.
type Number interface {
	constraints.Integer | constraints.Float
}

// Ring is an entity capable of keeping only the a number of most recent values.
type Ring[N Number] interface {
	Put(N)                                  // Put adds a sample to the Ring.
	Min() float64                           // Min returns the smallest value of all values in the Ring.
	Max() float64                           // Max returns the largest value of all values in the Ring.
	Average() float64                       // Average computes the average value of all values in the Ring.
	MinMaxAvg() (float64, float64, float64) // MinMaxAvg computes the min, max, and average values of all values in the Ring.
	Percentile(p float32) float64           // Percentile computes a given percentile value of all values in the Ring.
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

// Min returns the smallest value of all values in the Ring.
// If the ring is empty, retuns math.MaxFloat64.
func (r *ring[N]) Min() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	smallest := math.MaxFloat64

	// account for the case of not having any elems.
	if r.puts == 0 {
		return smallest
	}

	// account for the case of not having enough
	// elements to compute over the whole window.
	items := min(r.window, r.puts)

	for i := 0; i < items; i++ {
		smallest = min(smallest, float64(r.entries[i]))
	}
	return smallest
}

// Min returns the smallest value of all values in the Ring.
// If the ring is empty, retuns (-1)*math.MaxFloat64.
func (r *ring[N]) Max() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	largest := -math.MaxFloat64

	// account for the case of not having any elems.
	if r.puts == 0 {
		return largest
	}

	// account for the case of not having enough
	// elements to compute over the whole window.
	items := min(r.window, r.puts)

	for i := 0; i < items; i++ {
		largest = max(largest, float64(r.entries[i]))
	}
	return largest
}

// Average computes the average value of all values in the Ring.
// If the ring is empty, returns 0.
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
	items := min(r.window, r.puts)

	sum := N(0)
	for i := 0; i < items; i++ {
		sum += r.entries[i]
	}
	return float64(sum) / float64(items)
}

// MinMaxAvg computes the min, max, and average values of all values in the Ring.
// It is more performant than invoking Min(), Max(), and Average() independently.
// If the ring is empty, it returns math.MaxFloat64, -math.MaxFloat64, and 0.
func (r *ring[N]) MinMaxAvg() (float64, float64, float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	smallest := math.MaxFloat64
	largest := -math.MaxFloat64
	avg := float64(0)

	// account for the case of not having any elems.
	if r.puts == 0 {
		return smallest, largest, avg
	}

	// account for the case of not having enough
	// elements to compute over the whole window.
	items := min(r.window, r.puts)

	sum := N(0)
	for i := 0; i < items; i++ {
		smallest = min(smallest, float64(r.entries[i]))
		largest = max(largest, float64(r.entries[i]))
		sum += r.entries[i]
	}
	return smallest, largest, float64(sum) / float64(items)
}

// Percentile computes a given percentile value of all values in the Ring.
func (r *ring[N]) Percentile(p float32) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// account for an invalid percentile value.
	if p > 1.00 || p < 0.00 {
		return 0
	}

	// account for the case of not having any elems.
	if r.puts == 0 {
		return 0
	}

	// account for the case of not having enough
	// elements to compute over the whole window.
	items := min(r.window, r.puts)

	// copy the relevant part and sort it.
	sortedEntries := make([]N, items)
	copy(sortedEntries, r.entries[:items])
	sort.Slice(sortedEntries, func(i, j int) bool { return sortedEntries[i] < sortedEntries[j] })

	// compute the index of the percentile element.
	index := int(float32(items-1) * p)

	// return the element at the computed index.
	return float64(sortedEntries[index])
}

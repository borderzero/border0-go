package ring

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Put(t *testing.T) {
	tests := []struct {
		name                 string
		ring                 *ring[int]
		insertElem           int
		expectedEntriesAfter []int
		expectedPutsAfter    int
	}{
		{
			name: "putting into empty ring",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{0, 0, 0, 0, 0},
				window:  5,
				puts:    0,
			},
			insertElem:           10,
			expectedPutsAfter:    1,
			expectedEntriesAfter: []int{10, 0, 0, 0, 0},
		},
		{
			name: "putting into non-empty (but not full) ring",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 0, 0, 0, 0},
				window:  5,
				puts:    1,
			},
			insertElem:           20,
			expectedPutsAfter:    2,
			expectedEntriesAfter: []int{10, 20, 0, 0, 0},
		},
		{
			name: "putting into full ring (should overwrite oldest elem)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			insertElem:           35,
			expectedPutsAfter:    6,
			expectedEntriesAfter: []int{35, 15, 20, 25, 30}, // 5 % 5 is 0, so index 0 is overwritten
		},
		{
			name: "putting into overly full ring (should overwrite module window size elem)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    12,
			},
			insertElem:           35,
			expectedPutsAfter:    13,
			expectedEntriesAfter: []int{10, 15, 35, 25, 30}, // 12 % 5 is 2, so index 2 is overwritten
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.ring.Put(test.insertElem)

			assert.EqualValues(t, test.expectedEntriesAfter, test.ring.entries)
			assert.Equal(t, test.expectedPutsAfter, test.ring.puts)
		})
	}
}

func Test_Min(t *testing.T) {
	tests := []struct {
		name     string
		ring     *ring[int]
		expected float64
	}{
		{
			name: "returns math.MaxFloat64 when there are no elements",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{},
				window:  5,
				puts:    0,
			},
			expected: math.MaxFloat64,
		},
		{
			name: "returns correct min when window is not full",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			expected: 10,
		},
		{
			name: "returns correct min when window is filled exactly once",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			expected: 10,
		},
		{
			name: "returns correct min when window is filled and beyond",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			expected: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.ring.Min())
		})
	}
}

func Test_Max(t *testing.T) {
	tests := []struct {
		name     string
		ring     *ring[int]
		expected float64
	}{
		{
			name: "returns (-1)*math.MaxFloat64 when there are no elements",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{},
				window:  5,
				puts:    0,
			},
			expected: -math.MaxFloat64,
		},
		{
			name: "returns correct max when window is not full",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			expected: 30,
		},
		{
			name: "returns correct max when window is filled exactly once",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			expected: 30,
		},
		{
			name: "returns correct max when window is filled and beyond",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			expected: 30,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.ring.Max())
		})
	}
}

func Test_Average(t *testing.T) {
	tests := []struct {
		name     string
		ring     *ring[int]
		expected float64
	}{
		{
			name: "returns zero when there are no elements",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{},
				window:  5,
				puts:    0,
			},
			expected: 0,
		},
		{
			name: "returns correct average when window is not full",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			expected: 20,
		},
		{
			name: "returns correct average when window is filled exactly once",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			expected: 20,
		},
		{
			name: "returns correct average when window is filled and beyond",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			expected: 20,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.ring.Average())
		})
	}
}

func Test_MinMaxAvg(t *testing.T) {
	tests := []struct {
		name        string
		ring        *ring[int]
		expectedMin float64
		expectedMax float64
		expectedAvg float64
	}{
		{
			name: "returns math.MaxFloat64, -math.MaxFloat64, and 0 when there are no elements",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{},
				window:  5,
				puts:    0,
			},
			expectedMin: math.MaxFloat64,
			expectedMax: -math.MaxFloat64,
			expectedAvg: 0,
		},
		{
			name: "returns correct values when window is not full",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			expectedMin: 10,
			expectedMax: 30,
			expectedAvg: 20,
		},
		{
			name: "returns correct values when window is filled exactly once",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			expectedMin: 10,
			expectedMax: 30,
			expectedAvg: 20,
		},
		{
			name: "returns correct values when window is filled and beyond",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			expectedMin: 10,
			expectedMax: 30,
			expectedAvg: 20,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			minimum, maximum, average := test.ring.MinMaxAvg()
			assert.Equal(t, test.expectedMin, minimum)
			assert.Equal(t, test.expectedMax, maximum)
			assert.Equal(t, test.expectedAvg, average)
		})
	}
}

func Test_Percentile(t *testing.T) {
	tests := []struct {
		name       string
		ring       *ring[int]
		percentile float32
		expected   float64
	}{
		{
			name: "returns zero when percentile value is under 0.0",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: -0.1,
			expected:   0,
		},
		{
			name: "returns zero when percentile value is over 1.0",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: 1.1,
			expected:   0,
		},
		{
			name: "returns zero when there are no items in ring",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{},
				window:  5,
				puts:    0,
			},
			percentile: 0.5,
			expected:   0,
		},
		{
			name: "returns correct percentile when window is not full (p0)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			percentile: 0.0,
			expected:   10,
		},
		{
			name: "returns correct percentile when window is not full (p50)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			percentile: 0.5,
			expected:   20,
		},
		{
			name: "returns correct percentile when window is not full (p90)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			percentile: 0.9,
			expected:   20,
		},
		{
			name: "returns correct percentile when window is not full (p100)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 20, 30},
				window:  5,
				puts:    3,
			},
			percentile: 1.0,
			expected:   30,
		},
		{
			name: "returns correct percentile when window is exactly full (p0)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: 0.0,
			expected:   10,
		},
		{
			name: "returns correct percentile when window is excatly full (p50)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: 0.5,
			expected:   20,
		},
		{
			name: "returns correct percentile when window is exactly full (p90)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: 0.9,
			expected:   25,
		},
		{
			name: "returns correct percentile when window is not full (p100)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    5,
			},
			percentile: 1.0,
			expected:   30,
		},
		{
			name: "returns correct percentile when window is full and beyond (p0)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			percentile: 0.0,
			expected:   10,
		},
		{
			name: "returns correct percentile when window is full and beyond (p50)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			percentile: 0.5,
			expected:   20,
		},
		{
			name: "returns correct percentile when window is full and beyond (p90)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			percentile: 0.9,
			expected:   25,
		},
		{
			name: "returns correct percentile when window is full and beyond (p100)",
			ring: &ring[int]{
				mu:      sync.RWMutex{},
				entries: []int{10, 15, 20, 25, 30},
				window:  5,
				puts:    10,
			},
			percentile: 1.0,
			expected:   30,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, test.ring.Percentile(test.percentile))
		})
	}
}

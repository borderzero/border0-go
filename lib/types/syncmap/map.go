package syncmap

import "sync"

// Map is a generics typed wrapper around sync.Map. This is useful because it avoids
// requiring repetitive type assertions and elliminates the chances of a panic.
type Map[K comparable, V any] struct {
	inner *sync.Map
}

// New returns a newly initialized Map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{inner: new(sync.Map)}
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	if v, ok := m.inner.Load(key); ok {
		return v.(V), true
	}
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if v, loaded := m.inner.LoadAndDelete(key); loaded {
		return v.(V), true
	}
	return value, false
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it stores and
// returns the given value. The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.inner.LoadOrStore(key, value)
	return a.(V), loaded
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's contents:
// no key will be visited more than once, but if the value for any key is stored or deleted
// concurrently (including by f), Range may reflect any mapping for that key from any point
// during the Range call. Range does not block other methods on the receiver; even f itself
// may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns false after a
// constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

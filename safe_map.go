package maps

import (
	"sync"
)

// SafeMap is a go map that is safe for concurrent use and that uses a standard set of functions
// shared with other Map-like types.
//
// The recommended way to create a SafeMap is to first declare a concrete type alias, and then call
// new on it, like this:
//   type MyMap = SafeMap[string,int]
//
//   m := new(MyMap)
//
// This will allow you to swap in a different kind of Map just by changing the type.
type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	items StdMap[K, V]
}

// Clear resets the map to an empty map.
func (m *SafeMap[K, V]) Clear() {
	if m.items == nil {
		return
	}
	m.Lock()
	m.items = nil
	m.Unlock()
}

// Set sets the key to the given value.
func (m *SafeMap[K, V]) Set(k K, v V) {
	m.Lock()
	if m.items == nil {
		m.items = map[K]V{k: v}
	} else {
		m.items[k] = v
	}
	m.Unlock()
}

// Get returns the value based on its key. If it does not exist, an empty string will be returned.
func (m *SafeMap[K, V]) Get(k K) (v V) {
	v, _ = m.Load(k)
	return
}

// Has returns true if the given key exists in the map.
func (m *SafeMap[K, V]) Has(k K) (exists bool) {
	_, exists = m.Load(k)
	return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load().
func (m *SafeMap[K, V]) Load(k K) (v V, ok bool) {
	if m.items == nil {
		return
	}
	m.RLock()
	if m.items != nil {
		v, ok = m.items[k]
	}
	m.RUnlock()
	return
}

// Delete removes the key from the map. If the key does not exist, nothing happens.
func (m *SafeMap[K, V]) Delete(k K) {
	m.Lock()
	m.items.Delete(k)
	m.Unlock()
}

// Values returns a slice of the values. It will return a nil slice if the map is empty.
// Multiple calls to Values will result in the same list of values, but may be in a different order.
func (m *SafeMap[K, V]) Values() (v []V) {
	if m.items == nil {
		return
	}
	m.RLock()
	v = m.items.Values()
	m.RUnlock()
	return
}

// Keys returns a slice of the keys. It will return a nil slice if the map is empty.
// Multiple calls to Keys will result in the same list of keys, but may be in a different order.
func (m *SafeMap[K, V]) Keys() (keys []K) {
	if m.items == nil {
		return nil
	}
	m.RLock()
	keys = m.items.Keys()
	m.RUnlock()
	return
}

// Len returns the number of items in the map
func (m *SafeMap[K, V]) Len() (l int) {
	if m.items == nil {
		return
	}
	m.RLock()
	l = m.items.Len()
	m.RUnlock()
	return
}

// Range will call the given function with every key and value in the map.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
// During this process, the map will be locked, so do not pass a function that will take
// significant amounts of time, nor will call into other methods of the SafeMap which might also need a lock.
func (m *SafeMap[K, V]) Range(f func(k K, v V) bool) {
	if m.items == nil {
		return
	}
	m.RLock()
	defer m.RUnlock()
	m.items.Range(f)
}

// Merge merges the given  map with the current one. The given one takes precedent on collisions.
func (m *SafeMap[K, V]) Merge(in MapI[K, V]) {
	if m.items == nil {
		m.items = make(map[K]V, in.Len())
	}
	m.Lock()
	defer m.Unlock()
	m.items.Merge(in)
}

// Equal returns true if all the keys in the given map exist in this map, and the values are the same
func (m *SafeMap[K, V]) Equal(m2 MapI[K, V]) bool {
	m.RLock()
	defer m.RUnlock()
	return m.items.Equal(m2)
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (m *SafeMap[K, V]) MarshalBinary() ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	return m.items.MarshalBinary()
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeMap.
func (m *SafeMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	m.Lock()
	defer m.Unlock()
	return m.items.UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *SafeMap[K, V]) MarshalJSON() (out []byte, err error) {
	m.RLock()
	defer m.RUnlock()
	return m.items.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SafeMap.
// The JSON must start with an object.
func (m *SafeMap[K, V]) UnmarshalJSON(in []byte) (err error) {
	m.Lock()
	defer m.Unlock()
	return m.items.UnmarshalJSON(in)
}

// String outputs the map as a string.
func (m *SafeMap[K, V]) String() string {
	m.RLock()
	defer m.RUnlock()
	return m.items.String()
}

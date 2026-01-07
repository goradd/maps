package maps

import (
	"fmt"
	"iter"
	"slices"
	"strings"
	"sync"
)

// SafeSliceMap is a go map that uses a slice to save the order of its keys so that the map can
// be ranged in a predictable order. SafeSliceMap is safe for concurrent use.
//
// By default, the order will be the same order that items were inserted,
// i.e. a FIFO list, which is similar to how PHP arrays work. You can also define a sort function on the list
// to keep it sorted.
//
// The recommended way to create a SliceMap is to first declare a concrete type alias, and then call
// new on it, like this:
//
//	type MyMap = SafeSliceMap[string,int]
//
//	m := new(MyMap)
//
// This will allow you to swap in a different kind of Map just by changing the type.
//
// Call SetSortFunc to give the map a function that will keep the keys sorted in a particular order.
//
// Do not make a copy of a SafeSliceMap using the equality operator. Use Clone() instead.
type SafeSliceMap[K comparable, V any] struct {
	mu sync.RWMutex
	sm SliceMap[K, V]
}

// NewSafeSliceMap creates a new SafeSliceMap.
// Pass in zero or more standard maps and the contents of those maps will be copied to the new SafeSliceMap.
func NewSafeSliceMap[K comparable, V any](sources ...map[K]V) *SafeSliceMap[K, V] {
	m := new(SafeSliceMap[K, V])
	for _, i := range sources {
		m.Copy(Cast(i))
	}
	return m
}

// SetSortFunc sets the sort function which will determine the order of the items in the map
// on an ongoing basis. Normally, items will iterate in the order they were added.
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
// The sort function receives both the keys and values, so it can use either or both to decide how to sort.
func (m *SafeSliceMap[K, V]) SetSortFunc(f func(key1, key2 K, val1, val2 V) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sm.SetSortFunc(f)
}

// Set sets the given key to the given value.
//
// If the key already exists, the range order will not change. If you want the order
// to change, call Delete first, and then Set.
func (m *SafeSliceMap[K, V]) Set(key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sm.Set(key, val)
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.
// If the index is bigger than
// the length, it puts it at the end. Negative indexes are backwards from the end.
func (m *SafeSliceMap[K, V]) SetAt(index int, key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sm.SetAt(index, key, val)
}

// Delete removes the item with the given key and returns the value.
func (m *SafeSliceMap[K, V]) Delete(key K) (val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sm.Delete(key)
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (m *SafeSliceMap[K, V]) Get(key K) (val V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Get(key)
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (m *SafeSliceMap[K, V]) Load(key K) (val V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Load(key)
}

// Has returns true if the given key exists in the map.
func (m *SafeSliceMap[K, V]) Has(key K) (ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Has(key)
}

// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (m *SafeSliceMap[K, V]) GetAt(position int) (val V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.GetAt(position)
}

// GetKeyAt returns the key based on its position. If the position is out of bounds, an empty value is returned.
func (m *SafeSliceMap[K, V]) GetKeyAt(position int) (key K) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.GetKeyAt(position)
}

// Values returns a slice of the values in the order they were added or sorted.
func (m *SafeSliceMap[K, V]) Values() (values []V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Values()
}

// Keys returns the keys of the map, in the order they were added or sorted.
func (m *SafeSliceMap[K, V]) Keys() (keys []K) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Keys()
}

// Len returns the number of items in the map.
func (m *SafeSliceMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Len()
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
// If you are using a sort function, you must save and restore the sort function in a separate operation
// since functions are not serializable.
func (m *SafeSliceMap[K, V]) MarshalBinary() (data []byte, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.MarshalBinary()
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeSliceMap.
func (m *SafeSliceMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sm.UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *SafeSliceMap[K, V]) MarshalJSON() (data []byte, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Json objects are unordered
	return m.sm.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Map.
// The JSON must start with an object.
func (m *SafeSliceMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sm.UnmarshalJSON(data)
}

// Merge the given map into the current one.
// Deprecated: Use copy instead.
func (m *SafeSliceMap[K, V]) Merge(in MapI[K, V]) {
	m.Copy(in)
}

// Copy will copy the given map into the current one.
func (m *SafeSliceMap[K, V]) Copy(in MapI[K, V]) {
	in.Range(func(k K, v V) bool {
		m.Set(k, v) // This will lock and unlock, making sure that a long operation does not deadlock another go routine.
		return true
	})
}

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
// During this process, the map will be locked, so do not pass a function that will take
// significant amounts of time, nor will call into other methods of the SafeSliceMap which might also need a lock.
// The workaround is to call Keys() and iterate over the returned copy of the keys, but making sure
// your function can handle the situation where the key no longer exists in the slice.
func (m *SafeSliceMap[K, V]) Range(f func(key K, value V) bool) {
	if m == nil || m.sm.items == nil { // prevent unnecessary lock
		return
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.sm.Range(f)
}

// Equal returns true if all the keys and values are equal, regardless of the order.
//
// If the values are not comparable, you should implement the Equaler interface on the values.
// Otherwise, you will get a runtime panic.
func (m *SafeSliceMap[K, V]) Equal(m2 MapI[K, V]) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sm.Equal(m2)
}

// Clear removes all the items in the map.
func (m *SafeSliceMap[K, V]) Clear() {
	m.mu.Lock()
	m.sm.Clear()
	m.mu.Unlock()
}

// String outputs the map as a string.
func (m *SafeSliceMap[K, V]) String() string {
	var s string

	s = "{"

	// Range will handle locking
	m.Range(func(k K, v V) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}

// All returns an iterator over all the items in the map in the order they were entered or sorted.
func (m *SafeSliceMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(yield)
	}
}

// KeysIter returns an iterator over all the keys in the map.
// During this process, the map will be locked, so do not pass a function that will take
// significant amounts of time, nor will call into other methods of the SafeSliceMap which might also need a lock.
func (m *SafeSliceMap[K, V]) KeysIter() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m == nil || m.sm.items == nil {
			return
		}
		m.mu.RLock()
		defer m.mu.RUnlock()
		m.sm.KeysIter()(yield)
	}
}

// ValuesIter returns an iterator over all the values in the map.
// During this process, the map will be locked, so do not pass a function that will take
// significant amounts of time, nor will call into other methods of the SafeSliceMap which might also need a lock.
func (m *SafeSliceMap[K, V]) ValuesIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m == nil || m.sm.items == nil {
			return
		}
		m.mu.RLock()
		defer m.mu.RUnlock()
		m.sm.ValuesIter()(yield)
	}
}

// Insert adds the values from seq to the end of the map.
// Duplicate keys are overridden but not moved.
// Will lock and unlock for each item in seq to give time to other go routines.
func (m *SafeSliceMap[K, V]) Insert(seq iter.Seq2[K, V]) {
	for k, v := range seq {
		m.Set(k, v)
	}
}

// CollectSafeSliceMap collects key-value pairs from seq into a new SafeSliceMap
// and returns it.
func CollectSafeSliceMap[K comparable, V any](seq iter.Seq2[K, V]) *SafeSliceMap[K, V] {
	m := new(SafeSliceMap[K, V])

	// no need to lock here since this is a private variable
	for k, v := range seq {
		m.sm.Set(k, v)
	}
	return m
}

// Clone returns a copy of the SafeSliceMap. This is a shallow clone of the keys and values:
// the new keys and values are set using ordinary assignment. The order is preserved.
func (m *SafeSliceMap[K, V]) Clone() *SafeSliceMap[K, V] {
	m1 := new(SafeSliceMap[K, V])
	m.mu.RLock()
	defer m.mu.RUnlock()
	m1.sm.items = m.sm.items.Clone()
	m1.sm.order = slices.Clone(m.sm.order)
	m1.sm.lessF = m.sm.lessF
	return m1
}

// DeleteFunc deletes any key/value pairs for which del returns true.
// Items are ranged in order.
// This function locks the entire slice structure for the entirety of the call,
// so be careful to avoid deadlocks when calling this on a very big structure.
func (m *SafeSliceMap[K, V]) DeleteFunc(del func(K, V) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sm.DeleteFunc(del)
}

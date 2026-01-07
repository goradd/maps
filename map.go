package maps

import (
	"iter"
)

// Map is a go map that uses a standard set of functions shared with other Map-like types.
//
// The recommended way to create a Map is to first declare a concrete type alias, and then call
// new on it, like this:
//
//	type MyMap = Map[string,int]
//
//	m := new(MyMap)
//
// This will allow you to swap in a different kind of Map just by changing the type.
type Map[K comparable, V any] struct {
	items StdMap[K, V]
}

// NewMap creates a new map that maps values of type K to values of type V.
// Pass in zero or more standard maps and the contents of those maps will be copied to the new Map.
func NewMap[K comparable, V any](sources ...map[K]V) *Map[K, V] {
	m := new(Map[K, V])
	for _, i := range sources {
		m.Copy(Cast(i))
	}
	return m
}

// Clear resets the map to an empty map
func (m *Map[K, V]) Clear() {
	m.items = nil
}

// Len returns the number of items in the map
func (m *Map[K, V]) Len() int {
	return m.items.Len()
}

// Range calls the given function for each key,value pair in the map.
// This is the same interface as sync.Map.Range().
// While its safe to call methods of the map from within the Range function, its discouraged.
// If you ever switch to one of the SafeMap maps, it will cause a deadlock.
func (m *Map[K, V]) Range(f func(k K, v V) bool) {
	m.items.Range(f)
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (m *Map[K, V]) Load(k K) (V, bool) {
	return m.items.Load(k)
}

// Get returns the value for the given key. If the key does not exist, the zero value will be returned.
func (m *Map[K, V]) Get(k K) V {
	return m.items.Get(k)
}

// Has returns true if the key exists.
func (m *Map[K, V]) Has(k K) bool {
	return m.items.Has(k)
}

// Delete removes the key from the map. If the key does not exist, nothing happens.
func (m *Map[K, V]) Delete(k K) V {
	return m.items.Delete(k)
}

// Keys returns a new slice containing the keys of the map.
func (m *Map[K, V]) Keys() []K {
	return m.items.Keys()
}

// Values returns a new slice containing the values of the map.
func (m *Map[K, V]) Values() []V {
	return m.items.Values()
}

// Set sets the key to the given value.
func (m *Map[K, V]) Set(k K, v V) {
	if m.items == nil {
		m.items = map[K]V{k: v}
	} else {
		m.items.Set(k, v)
	}
}

// Merge copies the items from in to the map, overwriting any conflicting keys.
// Deprecated: Call Copy instead.
func (m *Map[K, V]) Merge(in MapI[K, V]) {
	m.Copy(in)
}

// Copy copies the items from in to the map, overwriting any conflicting keys.
func (m *Map[K, V]) Copy(in MapI[K, V]) {
	if m.items == nil {
		m.items = make(map[K]V, in.Len())
	}
	m.items.Copy(in)
}

// Equal returns true if all the keys and values are equal.
//
// If the values are not comparable, you should implement the Equaler interface on the values.
// Otherwise, you will get a runtime panic.
func (m *Map[K, V]) Equal(m2 MapI[K, V]) bool {
	return m.items.Equal(m2)
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (m *Map[K, V]) MarshalBinary() ([]byte, error) {
	return m.items.MarshalBinary()
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a Map.
//
// Note that you may need to register the map at init time with gob like this:
//
//	func init() {
//	  gob.Register(new(Map[keytype,valuetype]))
//	}
func (m *Map[K, V]) UnmarshalBinary(data []byte) (err error) {
	return m.items.UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *Map[K, V]) MarshalJSON() (out []byte, err error) {
	return m.items.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Map.
// The JSON must start with an object.
func (m *Map[K, V]) UnmarshalJSON(in []byte) (err error) {
	return m.items.UnmarshalJSON(in)
}

// String returns the map as a string.
func (m *Map[K, V]) String() string {
	return m.items.String()
}

// All returns an iterator over all the items in the map.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return m.items.All()
}

// KeysIter returns an iterator over all the keys in the map.
func (m *Map[K, V]) KeysIter() iter.Seq[K] {
	return m.items.KeysIter()
}

// ValuesIter returns an iterator over all the values in the map.
func (m *Map[K, V]) ValuesIter() iter.Seq[V] {
	return m.items.ValuesIter()
}

// Insert adds the values from seq to the map.
// Duplicate keys are overridden.
func (m *Map[K, V]) Insert(seq iter.Seq2[K, V]) {
	if m.items == nil {
		m.items = map[K]V{}
	}

	m.items.Insert(seq)
}

// CollectMap collects key-value pairs from seq into a new Map
// and returns it.
func CollectMap[K comparable, V any](seq iter.Seq2[K, V]) *Map[K, V] {
	m := new(Map[K, V])
	m.Insert(seq)
	return m
}

// Clone returns a copy of the Map. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func (m *Map[K, V]) Clone() *Map[K, V] {
	m1 := new(Map[K, V])
	m1.items = m.items.Clone()
	return m1
}

// DeleteFunc deletes any key/value pairs for which del returns true.
func (m *Map[K, V]) DeleteFunc(del func(K, V) bool) {
	m.items.DeleteFunc(del)
}

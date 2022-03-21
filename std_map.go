package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
)

// StdMap wraps a standard go map with a standard set of functions shared with other MapI-like types.
//
// The zero value is NOT settable. Use NewStdMap to create a new StdMap object, or use standard
// map instantiation syntax like this:
//   m := StdMap[string, int]{"a":1}
//
// StdMap is mostly a convenience type for making a standard Go map into a MapI interface.
// Generally, you should use Map instead, as it presents a consistent interface that allows you
// to swap the underlying type without changing implemented code.
type StdMap[K comparable, V any] map[K]V

// NewStdMap creates a new map that maps values of type K to values of type V.
// Pass in zero or more standard maps and the contents of those maps will be copied to the new StdMap.
// You can also create a new StdMap like this:
//   m := StdMap[string, int]{"a":1}
func NewStdMap[K comparable, V any](sources ...map[K]V) StdMap[K, V] {
	m := StdMap[K, V]{}
	for _, i := range sources {
		m.Merge(Cast(i))
	}
	return m
}

// Cast is a convenience method for casting a standard Go map to a StdMap type.
// Note that this is a cast, so the return value is the equivalent map of what
// was past in. Use this primarily to make a standard map into a MapI object.
func Cast[M ~map[K]V, K comparable, V any](m M) StdMap[K, V] {
	return StdMap[K, V](m)
}

// Clear resets the map to an empty map
func (m StdMap[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
}

// Len returns the number of items in the map.
func (m StdMap[K, V]) Len() int {
	return len(m)
}

// Merge copies the items from in to the map, overwriting any conflicting keys.
func (m StdMap[K, V]) Merge(in MapI[K, V]) {
	if m == nil {
		panic("cannot merge into a nil map")
	}
	in.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})
}

// Range calls the given function for each key,value pair in the map.
// This is the same interface as sync.Map.Range().
// While its safe to call methods of the map from within the Range function, its discouraged.
// If you ever switch to one of the SafeMap maps, it will cause a deadlock.
func (m StdMap[K, V]) Range(f func(k K, v V) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (m StdMap[K, V]) Load(k K) (v V, ok bool) {
	if m == nil {
		return
	}
	v, ok = m[k]
	return
}

// Get returns the value for the given key. If the key does not exist, the zero value will be returned.
func (m StdMap[K, V]) Get(k K) (v V) {
	v, _ = m.Load(k)
	return
}

// Has returns true if the key exists.
func (m StdMap[K, V]) Has(k K) (exists bool) {
	_, exists = m.Load(k)
	return
}

// Set sets the given key to the given value.
func (m StdMap[K, V]) Set(k K, v V) {
	if m == nil {
		panic("cannot call Set() on a nil map")
	}
	m[k] = v
}

// Delete removes the key from the map. If the key does not exist, nothing happens.
func (m StdMap[K, V]) Delete(k K) {
	delete(m, k)
}

// Keys returns a new slice containing the keys of the map.
func (m StdMap[K, V]) Keys() (keys []K) {
	if m.Len() == 0 {
		return
	}

	keys = make([]K, m.Len())

	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns a new slice containing the values of the map.
func (m StdMap[K, V]) Values() (values []V) {
	if m.Len() == 0 {
		return
	}
	values = make([]V, m.Len())
	var i int
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Equal returns true if all the keys and values are equal.
//
// If the values are not comparable, you should implement the Equaler interface on the values.
// Otherwise you will get a runtime panic.
func (m StdMap[K, V]) Equal(m2 MapI[K, V]) bool {
	if m.Len() != m2.Len() {
		return false
	}
	ret := true
	m2.Range(func(k K, v V) bool {
		if v2, ok := m[k]; !ok || !equalValues(v, v2) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

// String returns a string representation of the map.
func (m StdMap[K, V]) String() string {
	s := fmt.Sprintf("%#v", m)
	loc := strings.IndexRune(s, '{')
	return s[loc:]
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (m StdMap[K, V]) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	err := enc.Encode(map[K]V(m))
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a Map.
//
// Note that you will likely need to register the unmarshaller at init time with gob like this:
//    func init() {
//      gob.Register(new(Map[K,V]))
//    }
func (m *StdMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var v map[K]V
	err = dec.Decode(&v)
	*m = v
	return
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m StdMap[K, V]) MarshalJSON() (out []byte, err error) {
	v := map[K]V(m)
	return json.Marshal(v)
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a StdMap.
// The JSON must start with an object.
func (m *StdMap[K, V]) UnmarshalJSON(in []byte) (err error) {
	var v map[K]V

	err = json.Unmarshal(in, &v)
	*m = v
	return
}

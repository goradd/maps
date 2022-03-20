package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// Map is a go map that uses a standard set of functions shared with other Map-like types.
//
// The recommended way to create a Map is to first declare a concrete type alias, and then call
// new on it, like this:
//   type MyMap = Map[string,int]
//
//   m := new(MyMap)
//
//   This will allow you to swap in a different kind of Map just by changing the type.
type Map[K comparable, V any] struct {
	items StdMap[K, V]
}

func (m Map[K, V]) Clear() {
	m.items.Clear()
}

func (m Map[K, V]) Len() int {
	return m.items.Len()
}

func (m Map[K, V]) Range(f func(k K, v V) bool) {
	m.items.Range(f)
}

func (m Map[K, V]) Load(k K) (V, bool) {
	return m.items.Load(k)
}

func (m Map[K, V]) Get(k K) V {
	return m.items.Get(k)
}

func (m Map[K, V]) Has(k K) bool {
	return m.items.Has(k)
}

func (m Map[K, V]) Keys() []K {
	return m.items.Keys()
}

func (m Map[K, V]) Values() []V {
	return m.items.Values()
}

func (m *Map[K, V]) Set(k K, v V) {
	if m.items == nil {
		m.items = map[K]V{k: v}
	} else {
		m.items.Set(k, v)
	}
}

func (m *Map[K, V]) Merge(in MapI[K, V]) {
	if m.items == nil {
		m.items = make(map[K]V, in.Len())
	}
	m.items.Merge(in)
}

// Equal returns true if all the keys and values are equal.
//
// You will get a runtime panic if your values are not comparable, or your values do
// not satisfy the Equaler interface.
func (m Map[K, V]) Equal(m2 MapI[K, V]) bool {
	return m.items.Equal(m2)
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (m Map[K, V]) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	err := enc.Encode(map[K]V(m.items))
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a Map.
//
// Note that you will likely need to register the unmarshaller at init time with gob like this:
//    func init() {
//      gob.Register(new(Map[K,V]))
//    }
func (m *Map[K, V]) UnmarshalBinary(data []byte) (err error) {
	var v map[K]V

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err = dec.Decode(&v); err == nil {
		m.items = Cast(v)
	}
	return
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m Map[K, V]) MarshalJSON() (out []byte, err error) {
	out, err = json.Marshal(map[K]V(m.items))
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Map.
// The JSON must start with an object.
func (m *Map[K, V]) UnmarshalJSON(in []byte) (err error) {
	var v map[K]V
	if err = json.Unmarshal(in, &v); err == nil {
		m.items = Cast(v)
	}
	return
}

func (m Map[K, V]) String() string {
	return m.items.String()
}

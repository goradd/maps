package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"iter"
)

// Set is a collection that keeps track of membership.
//
// The recommended way to create a Set is to first declare a concrete type alias, and then call
// new on it, like this:
//
//	type MySet = Set[string]
//
//	s := new(MySet)
//
// This will allow you to swap in a different kind of Set just by changing the type.
type Set[K comparable] struct {
	items StdMap[K, struct{}]
}

func NewSet[K comparable](values ...K) *Set[K] {
	s := new(Set[K])
	for _, k := range values {
		s.Add(k)
	}
	return s
}

// Clear resets the set to an empty set
func (m *Set[K]) Clear() {
	if m == nil {
		return
	}
	m.items = nil
}

// Len returns the number of items in the set
func (m *Set[K]) Len() int {
	if m == nil || m.items == nil {
		return 0
	}
	return m.items.Len()
}

// Range calls the given function for each member in the set.
// The function should return true to continue ranging, or false to stop.
// While its safe to call methods of the set from within the Range function, its discouraged.
// If you ever switch to one of the SafeSet sets, it will cause a deadlock.
func (m *Set[K]) Range(f func(k K) bool) {
	if m.Len() == 0 {
		return
	}
	for k := range m.items {
		if !f(k) {
			break
		}
	}
}

// Has returns true if the value exists in the set.
func (m *Set[K]) Has(k K) bool {
	if m.Len() == 0 {
		return false
	}
	return m.items.Has(k)
}

// Delete removes the value from the set. If the value does not exist, nothing happens.
func (m *Set[K]) Delete(k K) {
	if m.Len() == 0 {
		return
	}
	m.items.Delete(k)
}

// Values returns a new slice containing the values of the set.
func (m *Set[K]) Values() []K {
	if m.Len() == 0 {
		return nil
	}
	return m.items.Keys()
}

// Add adds the value to the set.
// If the value already exists, nothing changes.
func (m *Set[K]) Add(k ...K) SetI[K] {
	if m == nil {
		panic("cannot add values to a nil Set")
	}
	if m.items == nil {
		m.items = make(map[K]struct{})
	}
	for _, i := range k {
		m.items.Set(i, struct{}{})
	}
	return m
}

// Merge adds the values from the given set to the set.
// Deprecated: Call Copy instead.
func (m *Set[K]) Merge(in SetI[K]) {
	m.Copy(in)
}

// Copy adds the values from in to the set.
func (m *Set[K]) Copy(in SetI[K]) {
	if m == nil {
		panic("cannot copy to a nil Set")
	}
	if in == nil || in.Len() == 0 {
		return
	}
	if m.items == nil {
		m.items = make(map[K]struct{}, in.Len())
	}
	in.Range(func(k K) bool {
		m.items[k] = struct{}{}
		return true
	})
}

// Equal returns true if the two sets are the same length and contain the same values.
func (m *Set[K]) Equal(m2 SetI[K]) bool {
	if m == nil {
		return m2.Len() == 0
	}
	if m.Len() != m2.Len() {
		return false
	}
	ret := true
	m2.Range(func(k K) bool {
		if !m.Has(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

// MarshalBinary implements the BinaryMarshaler interface to convert the set to a byte stream.
func (m *Set[K]) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	err := enc.Encode(m.Values())
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a Set.
//
// Note that you may need to register the set at init time with gob like this:
//
//	func init() {
//	  gob.Register(new(Set[keytype]))
//	}
func (m *Set[K]) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var v []K
	err = dec.Decode(&v)
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *Set[K]) MarshalJSON() (out []byte, err error) {
	if m.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(m.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Set.
// The JSON must start with a list.
func (m *Set[K]) UnmarshalJSON(in []byte) (err error) {
	var v []K

	err = json.Unmarshal(in, &v)
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

// String returns the set as a string.
func (m *Set[K]) String() string {
	ret := "{"
	if m.Len() != 0 {
		for i, v := range m.Values() {
			ret += fmt.Sprintf("%#v", v)
			if i < m.Len()-1 {
				ret += ","
			}
		}
	}
	ret += "}"
	return ret
}

// All returns an iterator over all the items in the set. Order is not determinate.
func (m *Set[K]) All() iter.Seq[K] {
	if m.Len() == 0 {
		return func(yield func(K) bool) {
			return
		}
	}
	return m.items.KeysIter()
}

// Insert adds the values from seq to the map.
// Duplicates are overridden.
func (m *Set[K]) Insert(seq iter.Seq[K]) {
	if m == nil {
		panic("cannot insert into a nil Set")
	}
	if m.items == nil {
		m.items = NewStdMap[K, struct{}]()
	}

	for k := range seq {
		m.Add(k)
	}
}

// CollectSet collects values from seq into a new Set
// and returns it.
func CollectSet[K comparable](seq iter.Seq[K]) *Set[K] {
	m := NewSet[K]()
	m.Insert(seq)
	return m
}

// Clone returns a copy of the Set. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func (m *Set[K]) Clone() *Set[K] {
	m1 := NewSet[K]()
	if m.Len() != 0 {
		m1.items = m.items.Clone()
	}
	return m1
}

// DeleteFunc deletes any values for which del returns true.
func (m *Set[K]) DeleteFunc(del func(K) bool) {
	if m.Len() == 0 {
		return
	}
	del2 := func(k K, s struct{}) bool {
		return del(k)
	}
	m.items.DeleteFunc(del2)
}

package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
)

// SliceSet implements a set of values that will be returned in the order added,
// or based on a sorting function.
//
// The differences between SliceSet and OrderedSet are:
//   - A SliceSet can hold comparable types, vs. OrderedSet can only hold cmp.Ordered types.
//   - A SliceSet sorts whenever items are added, while OrderedSet only sorts when the order is asked for.
//
// SliceSet is built on top of SliceMap.
type SliceSet[K comparable] struct {
	sm SliceMap[K, struct{}]
}

func NewSliceSet[K comparable](values ...K) *SliceSet[K] {
	s := new(SliceSet[K])
	for _, k := range values {
		s.Add(k)
	}
	return s
}

// SetSortFunc sets the sort function which will determine the order of the items in the set
// on an ongoing basis. Normally, items will iterate in the order they were added.
//
// When you call SetSortFunc, the values will be sorted. To turn off sorting, set the sort function to nil.
//
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
func (m *SliceSet[K]) SetSortFunc(f func(val1, val2 K) bool) {
	if m == nil {
		panic("cannot set a sort function on a nil SliceSet")
	}
	if f == nil {
		m.sm.SetSortFunc(nil)
		return
	}
	f2 := func(key1, key2 K, val1, val2 struct{}) bool {
		return f(key1, key2)
	}
	m.sm.SetSortFunc(f2)
}

// Clear resets the set to an empty set
func (m *SliceSet[K]) Clear() {
	if m == nil {
		return
	}
	m.sm.Clear()
}

// Len returns the number of items in the set
func (m *SliceSet[K]) Len() int {
	if m == nil {
		return 0
	}
	return m.sm.Len()
}

// Range will range over the values in order.
func (m *SliceSet[K]) Range(f func(k K) bool) {
	if m.Len() == 0 {
		return
	}
	values := m.Values()
	for _, k := range values {
		if !f(k) {
			break
		}
	}
}

// Has returns true if the value exists in the set.
func (m *SliceSet[K]) Has(k K) bool {
	if m.Len() == 0 {
		return false
	}
	return m.sm.Has(k)
}

// Delete removes the value from the set. If the value does not exist, nothing happens.
func (m *SliceSet[K]) Delete(k K) {
	if m.Len() == 0 {
		return
	}
	m.sm.Delete(k)
}

// Equal returns true if the two sets are the same length and contain the same values.
func (m *SliceSet[K]) Equal(m2 SetI[K]) bool {
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

// Values returns a new slice containing the values of the set in order.
func (m *SliceSet[K]) Values() []K {
	if m.Len() == 0 {
		return nil
	}
	v := m.sm.Keys()
	return v
}

// Add adds the value to the set.
// If the value already exists, nothing changes.
func (m *SliceSet[K]) Add(k ...K) SetI[K] {
	if m == nil {
		panic("cannot add values to a nil Set")
	}
	for _, i := range k {
		m.sm.Set(i, struct{}{})
	}
	return m
}

// Copy adds the values from in to the set.
func (m *SliceSet[K]) Copy(in SetI[K]) {
	if m == nil {
		panic("cannot copy to a nil Set")
	}
	if in == nil || in.Len() == 0 {
		return
	}
	for i := range in.All() {
		m.Add(i)
	}
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *SliceSet[K]) MarshalJSON() (out []byte, err error) {
	if m.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(m.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Set.
// The JSON must start with a list.
func (m *SliceSet[K]) UnmarshalJSON(in []byte) (err error) {
	var v []K

	err = json.Unmarshal(in, &v)
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

// MarshalBinary implements the BinaryMarshaler interface to convert the set to a byte stream.
func (m *SliceSet[K]) MarshalBinary() ([]byte, error) {
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
func (m *SliceSet[K]) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var v []K
	err = dec.Decode(&v)
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

// All returns an iterator over all the items in the set. Order is determinate.
func (m *SliceSet[K]) All() iter.Seq[K] {
	if m.Len() == 0 {
		return func(yield func(K) bool) {
			return
		}
	}
	v := m.Values()
	return slices.Values(v)
}

// Insert adds the values from seq to the map.
// Duplicates are overridden.
func (m *SliceSet[K]) Insert(seq iter.Seq[K]) {
	if m == nil {
		panic("cannot insert into a nil Set")
	}
	for i := range seq {
		m.Add(i)
	}
}

// Clone returns a copy of the Set. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func (m *SliceSet[K]) Clone() *SliceSet[K] {
	m1 := NewSliceSet[K]()
	if m != nil {
		for i := range m.All() {
			m1.Add(i)
		}
	}
	return m1
}

// DeleteFunc deletes any values for which del returns true.
func (m *SliceSet[K]) DeleteFunc(del func(K) bool) {
	if m.Len() == 0 {
		return
	}
	del2 := func(k K, s struct{}) bool {
		return del(k)
	}
	m.sm.DeleteFunc(del2)
}

// String returns the set as a string.
func (m *SliceSet[K]) String() string {
	if m == nil {
		return "{}"
	}
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

// Merge adds the values from the given set to the set.
// Deprecated: Call Copy instead.
func (m *SliceSet[K]) Merge(in SetI[K]) {
	m.Copy(in)
}

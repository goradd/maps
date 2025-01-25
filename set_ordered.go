package maps

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
)

// OrderedSet implements a set of values that will be returned sorted.
//
// Ordered sets are useful when in general you don't care about ordering, but
// you would still like the same values to be presented in the same order when
// they are asked for. Examples include test code, iterators, values stored in a database,
// or values that will be presented to a user.
type OrderedSet[K cmp.Ordered] struct {
	Set[K]
}

func NewOrderedSet[K cmp.Ordered](values ...K) *OrderedSet[K] {
	s := new(OrderedSet[K])
	for _, k := range values {
		s.Add(k)
	}
	return s
}

// Clear resets the set to an empty set
func (m *OrderedSet[K]) Clear() {
	if m == nil {
		return
	}
	m.Set.Clear()
}

// Len returns the number of items in the set
func (m *OrderedSet[K]) Len() int {
	if m == nil || m.items == nil {
		return 0
	}
	return m.Set.Len()
}

// Range will range over the values in order.
func (m *OrderedSet[K]) Range(f func(k K) bool) {
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
func (m *OrderedSet[K]) Has(k K) bool {
	if m.Len() == 0 {
		return false
	}
	return m.Set.Has(k)
}

// Delete removes the value from the set. If the value does not exist, nothing happens.
func (m *OrderedSet[K]) Delete(k K) {
	if m.Len() == 0 {
		return
	}
	m.Set.Delete(k)
}

// Equal returns true if the two sets are the same length and contain the same values.
func (m *OrderedSet[K]) Equal(m2 SetI[K]) bool {
	if m == nil {
		return m2.Len() == 0
	}
	return m.Set.Equal(m2)
}

// Values returns a new slice containing the values of the set.
func (m *OrderedSet[K]) Values() []K {
	if m.Len() == 0 {
		return nil
	}
	v := m.items.Keys()
	slices.Sort(v)
	return v
}

// Add adds the value to the set.
// If the value already exists, nothing changes.
func (m *OrderedSet[K]) Add(k ...K) SetI[K] {
	if m == nil {
		panic("cannot add values to a nil Set")
	}
	m.Set.Add(k...)
	return m
}

// Copy adds the values from in to the set.
func (m *OrderedSet[K]) Copy(in SetI[K]) {
	if m == nil {
		panic("cannot copy to a nil Set")
	}
	m.Set.Copy(in)
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *OrderedSet[K]) MarshalJSON() (out []byte, err error) {
	if m.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(m.Values())
}

// All returns an iterator over all the items in the set. Order is determinate.
func (m *OrderedSet[K]) All() iter.Seq[K] {
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
func (m *OrderedSet[K]) Insert(seq iter.Seq[K]) {
	if m == nil {
		panic("cannot insert into a nil Set")
	}
	m.Set.Insert(seq)
}

// Clone returns a copy of the Set. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func (m *OrderedSet[K]) Clone() *OrderedSet[K] {
	m1 := NewOrderedSet[K]()
	if m != nil {
		m1.items = m.items.Clone()
	}
	return m1
}

// DeleteFunc deletes any values for which del returns true.
func (m *OrderedSet[K]) DeleteFunc(del func(K) bool) {
	if m.Len() == 0 {
		return
	}
	m.Set.DeleteFunc(del)
}

// String returns the set as a string.
func (m *OrderedSet[K]) String() string {
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

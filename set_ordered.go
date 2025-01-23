package maps

import (
	"cmp"
	"encoding/json"
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

// Range will range over the values in order.
func (m *OrderedSet[K]) Range(f func(k K) bool) {
	if m == nil || m.items == nil {
		return
	}
	values := m.Values()
	for _, k := range values {
		if !f(k) {
			break
		}
	}
}

// Values returns the values as a slice, in order.
func (m *OrderedSet[K]) Values() []K {
	v := m.items.Keys()
	slices.Sort(v)
	return v
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
	v := m.Values()
	return slices.Values(v)
}

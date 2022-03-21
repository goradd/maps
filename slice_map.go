package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SliceMap is a go map that uses a slice to save the order of its keys so that the map can
// be ranged in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list, which is similar to how PHP arrays work. You can also define a sort function on the list
// to keep it sorted.
//
// The recommended way to create a SliceMap is to first declare a concrete type alias, and then call
// new on it, like this:
//   type MyMap = SliceMap[string,int]
//
//   m := new(MyMap)
//
// This will allow you to swap in a different kind of Map just by changing the type.
//
// Call SetSortFunc to give the map a function that will keep the keys sorted in a particular order.
type SliceMap[K comparable, V any] struct {
	items StdMap[K, V]
	order []K
	lessF func(key1, key2 K, val1, val2 V) bool
}

// SetSortFunc sets the sort function which will determine the order of the items in the map
// on an ongoing basis. Normally, items will iterate in the order they were added.
//
// When you call SetSortFunc, the map keys will be sorted. To turn off sorting, set the sort function to nil.
//
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
// The sort function receives both the keys and values, so it can use either or both to decide how to sort.
func (m *SliceMap[K, V]) SetSortFunc(f func(key1, key2 K, val1, val2 V) bool) {
	m.lessF = f
	if f != nil && len(m.order) > 0 {
		sort.Slice(m.order, func(i, j int) bool {
			return f(m.order[i], m.order[j], m.items[m.order[i]], m.items[m.order[j]])
		})
	}
}

// Set sets the given key to the given value.
//
// If the key already exists, the range order will not change. If you want the order
// to change, call Delete first, and then Set.
func (m *SliceMap[K, V]) Set(key K, val V) {
	var ok bool
	var oldVal V

	if m.items == nil {
		m.items = make(map[K]V)
	}

	_, ok = m.items[key]
	if m.lessF != nil {
		if ok {
			// delete old key location
			loc := sort.Search(len(m.items), func(n int) bool {
				return !m.lessF(m.order[n], key, m.items[m.order[n]], oldVal)
			})
			m.order = append(m.order[:loc], m.order[loc+1:]...)
		}

		loc := sort.Search(len(m.order), func(n int) bool {
			return m.lessF(key, m.order[n], val, m.items[m.order[n]])
		})
		// insert
		m.order = append(m.order, key)
		copy(m.order[loc+1:], m.order[loc:])
		m.order[loc] = key
	} else {
		if !ok {
			m.order = append(m.order, key)
		}
	}
	m.items[key] = val
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.
// If the index is bigger than
// the length, it puts it at the end. Negative indexes are backwards from the end.
func (m *SliceMap[K, V]) SetAt(index int, key K, val V) {
	if m.lessF != nil {
		panic("cannot use SetAt if you are also using a sort function")
	}

	if index >= len(m.order) {
		m.Set(key, val)
		return
	}

	var ok bool
	var emptyKey K

	if _, ok = m.items[key]; ok {
		m.Delete(key)
	}
	if index <= -len(m.items) {
		index = 0
	}
	if index < 0 {
		index = len(m.items) + index
	}

	m.order = append(m.order, emptyKey)
	copy(m.order[index+1:], m.order[index:])
	m.order[index] = key

	m.items[key] = val
}

// Delete removes the item with the given key.
func (m *SliceMap[K, V]) Delete(key K) {
	if _, ok := m.items[key]; ok {
		if m.lessF != nil {
			oldVal := m.items[key]
			loc := sort.Search(len(m.items), func(n int) bool {
				return !m.lessF(m.order[n], key, m.items[m.order[n]], oldVal)
			})
			m.order = append(m.order[:loc], m.order[loc+1:]...)
		} else {
			for i, v := range m.order {
				if v == key {
					m.order = append(m.order[:i], m.order[i+1:]...)
					break
				}
			}
		}
		delete(m.items, key)
	}
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (m *SliceMap[K, V]) Get(key K) (val V) {
	return m.items.Get(key)
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.StdMap.Load()
func (m *SliceMap[K, V]) Load(key K) (val V, ok bool) {
	return m.items.Load(key)
}

// Has returns true if the given key exists in the map.
func (m *SliceMap[K, V]) Has(key K) (ok bool) {
	return m.items.Has(key)
}

// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (m *SliceMap[K, V]) GetAt(position int) (val V) {
	if position < len(m.order) && position >= 0 {
		val, _ = m.items[m.order[position]]
	}
	return
}

// GetKeyAt returns the key based on its position. If the position is out of bounds, an empty value is returned.
func (m *SliceMap[K, V]) GetKeyAt(position int) (key K) {
	if position < len(m.order) && position >= 0 {
		key = m.order[position]
	}
	return
}

// Values returns a slice of the values in the order they were added or sorted.
func (m *SliceMap[K, V]) Values() (vals []V) {
	return m.items.Values()
}

// Keys returns the keys of the map, in the order they were added or sorted
func (m *SliceMap[K, V]) Keys() (keys []K) {
	return m.items.Keys()
}

// Len returns the number of items in the map
func (m *SliceMap[K, V]) Len() int {
	return m.items.Len()
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
// If you are using a sort function, you must save and restore the sort function in a separate operation
// since functions are not serializable.
func (m *SliceMap[K, V]) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err = encoder.Encode(map[K]V(m.items))
	if err == nil {
		err = encoder.Encode(m.order)
	}
	data = buf.Bytes()
	return
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SliceMap.
func (m *SliceMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	var items map[K]V
	var order []K

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&items); err == nil {
		err = dec.Decode(&order)
	}

	if err == nil {
		m.items = items
		m.order = order
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (m *SliceMap[K, V]) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	return m.items.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SliceMap.
// The JSON must start with an object.
func (m *SliceMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	var items map[K]V

	if err = json.Unmarshal(data, &items); err == nil {
		m.items = items
		// Create a default order, since these are inherently unordered
		m.order = make([]K, len(m.items))
		i := 0
		for k := range m.items {
			m.order[i] = k
			i++
		}
	}
	return
}

// Merge the given map into the current one.
func (m *SliceMap[K, V]) Merge(in MapI[K, V]) {
	if in != nil {
		in.Range(func(k K, v V) bool {
			m.Set(k, v)
			return true
		})
	}
}

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (m *SliceMap[K, V]) Range(f func(key K, value V) bool) {
	if m.items != nil {
		for _, k := range m.order {
			if !f(k, m.items[k]) {
				break
			}
		}
	}
}

// Equal returns true if all the keys and values are equal, regardless of the order.
//
// If the values are not comparable, you should implement the Equaler interface on the values.
// Otherwise, you will get a runtime panic.
func (m *SliceMap[K, V]) Equal(m2 MapI[K, V]) bool {
	return m.items.Equal(m2)
}

// Clear removes all the items in the map.
func (m *SliceMap[K, V]) Clear() {
	m.items = nil
	m.order = nil
}

// String outputs the map as a string.
func (m *SliceMap[K, V]) String() string {
	var s string

	s = "{"
	m.Range(func(k K, v V) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}

// Equaler is the interface that implements an Equal function and that provides a way for the
// various MapI like objects to determine if they are equal.
//
// In particular, if your Map has
// non-comparible values, like a slice, but you would still like to call Equal() on that
// map, define an Equal function on the values to do the comparison. For example:
//
//   type mySlice []int
//
//   func (s mySlice) Equal(b any) bool {
//   	if s2, ok := b.(mySlice); ok {
//   		if len(s) == len(s2) {
//   			for i, v := range s2 {
//   				if s[i] != v {
//   					return false
//   				}
//   			}
//   			return true
//   		}
//   	}
//   	return false
//   }
type Equaler interface {
	Equal(a any) bool
}

func equalValues(a, b any) bool {
	if e, ok := a.(Equaler); ok {
		return e.Equal(b)
	}

	return a == b
}

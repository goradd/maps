package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
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
//
//	type MyMap = SliceMap[string,int]
//
//	m := new(MyMap)
//
// This will allow you to swap in a different kind of Map just by changing the type.
//
// Call SetSortFunc to give the map a function that will keep the keys sorted in a particular order.
type SliceMap[K comparable, V any] struct {
	items StdMap[K, V]
	order []K
	lessF func(key1, key2 K, val1, val2 V) bool
}

// NewSliceMap creates a new SliceMap.
// Pass in zero or more standard maps and the contents of those maps will be copied to the new SafeMap.
func NewSliceMap[K comparable, V any](sources ...map[K]V) *SliceMap[K, V] {
	m := new(SliceMap[K, V])
	for _, i := range sources {
		m.Copy(Cast(i))
	}
	return m
}

// SetSortFunc sets the sort function which will determine the order of the items in the map
// on an ongoing basis. Normally, items will iterate in the order they were added.
//
// When you call SetSortFunc, the map keys will be sorted. To turn off sorting, set the sort function to nil.
//
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
// The sort function receives both the keys and values, so it can use either or both to decide how to sort.
func (m *SliceMap[K, V]) SetSortFunc(f func(key1, key2 K, val1, val2 V) bool) {
	if m == nil {
		panic("cannot set a sort function on a nil SliceMap")
	}
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

	if m == nil {
		panic("cannot set a value on a nil SliceMap")
	}

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
		// will handle m.items == nil
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

// Delete removes the key from the map and returns the value. If the key does not exist, the zero value will be returned.
func (m *SliceMap[K, V]) Delete(key K) (val V) {
	if m == nil {
		return
	}

	if _, ok := m.items[key]; ok {
		val = m.items[key]
		if m.lessF != nil {
			loc := sort.Search(len(m.items), func(n int) bool {
				return !m.lessF(m.order[n], key, m.items[m.order[n]], val)
			})
			m.order = slices.Delete(m.order, loc, loc+1)
		} else {
			for i, v := range m.order {
				if v == key {
					m.order = slices.Delete(m.order, i, i+1)
					break
				}
			}
		}
		delete(m.items, key)
	}
	return
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (m *SliceMap[K, V]) Get(key K) (val V) {
	if m == nil {
		return
	}
	return m.items.Get(key)
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.StdMap.Load()
func (m *SliceMap[K, V]) Load(key K) (val V, ok bool) {
	if m == nil {
		return
	}
	return m.items.Load(key)
}

// Has returns true if the given key exists in the map.
func (m *SliceMap[K, V]) Has(key K) (ok bool) {
	if m == nil {
		return
	}
	return m.items.Has(key)
}

// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (m *SliceMap[K, V]) GetAt(position int) (val V) {
	if m == nil {
		return
	}
	if position < len(m.order) && position >= 0 {
		val, _ = m.items[m.order[position]]
	}
	return
}

// GetKeyAt returns the key based on its position. If the position is out of bounds, an empty value is returned.
func (m *SliceMap[K, V]) GetKeyAt(position int) (key K) {
	if m == nil {
		return
	}
	if position < len(m.order) && position >= 0 {
		key = m.order[position]
	}
	return
}

// Values returns a slice of the values in the order they were added or sorted.
func (m *SliceMap[K, V]) Values() (values []V) {
	if m == nil {
		return
	}
	for _, k := range m.order {
		values = append(values, m.items[k])
	}
	return values
}

// Keys returns a new slice of the keys of the map, in the order they were added or sorted
func (m *SliceMap[K, V]) Keys() (keys []K) {
	if m == nil {
		return
	}
	return slices.Clone(m.order)
}

// Len returns the number of items in the map
func (m *SliceMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return m.items.Len()
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
// If you are using a sort function, you must save and restore the sort function in a separate operation
// since functions are not serializable.
func (m *SliceMap[K, V]) MarshalBinary() (data []byte, err error) {
	if m == nil {
		return
	}
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

	if m == nil {
		panic("cannot Unmarshal into a nil SliceMap")
	}

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
	if m == nil {
		return
	}
	return m.items.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SliceMap.
// The JSON must start with an object.
func (m *SliceMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	var items map[K]V

	if m == nil {
		panic("cannot unmarshall into a nil SliceMap")
	}
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
// Deprecated: use Copy instead.
func (m *SliceMap[K, V]) Merge(in MapI[K, V]) {
	in.Range(func(k K, v V) bool {
		m.Set(k, v)
		return true
	})
}

// Copy copies the keys and values of in into the current one.
// Duplicate keys will have the values replaced, but not the order.
func (m *SliceMap[K, V]) Copy(in MapI[K, V]) {
	in.Range(func(k K, v V) bool {
		m.Set(k, v)
		return true
	})
}

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (m *SliceMap[K, V]) Range(f func(key K, value V) bool) {
	if m != nil && m.items != nil {
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
	if m == nil {
		return m2 == nil || m2.Len() == 0
	}
	return m.items.Equal(m2)
}

// Clear removes all the items in the map.
func (m *SliceMap[K, V]) Clear() {
	if m == nil {
		return
	}
	m.items = nil
	m.order = nil
}

// String outputs the map as a string.
func (m *SliceMap[K, V]) String() string {
	var s string

	if m == nil {
		return s
	}

	s = "{"
	m.Range(func(k K, v V) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}

// All returns an iterator over all the items in the map in the order they were entered or sorted.
func (m *SliceMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(yield)
	}
}

// KeysIter returns an iterator over all the keys in the map.
func (m *SliceMap[K, V]) KeysIter() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m == nil || m.items == nil {
			return
		}
		for _, k := range m.order {
			if !yield(k) {
				break
			}
		}
	}
}

// ValuesIter returns an iterator over all the values in the map.
func (m *SliceMap[K, V]) ValuesIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m == nil || m.items == nil {
			return
		}
		for _, k := range m.order {
			if !yield(m.items[k]) {
				break
			}
		}
	}
}

// Insert adds the values from seq to the end of the map.
// Duplicate keys are overridden but not moved.
func (m *SliceMap[K, V]) Insert(seq iter.Seq2[K, V]) {
	for k, v := range seq {
		m.Set(k, v)
	}
}

// CollectSliceMap collects key-value pairs from seq into a new SliceMap
// and returns it.
func CollectSliceMap[K comparable, V any](seq iter.Seq2[K, V]) *SliceMap[K, V] {
	m := new(SliceMap[K, V])
	m.Insert(seq)
	return m
}

// Clone returns a copy of the SliceMap. This is a shallow clone of the keys and values:
// the new keys and values are set using ordinary assignment. The order is preserved.
func (m *SliceMap[K, V]) Clone() *SliceMap[K, V] {
	m1 := new(SliceMap[K, V])
	m1.items = m.items.Clone()
	m1.order = slices.Clone(m.order)
	m1.lessF = m.lessF
	return m1
}

// DeleteFunc deletes any key/value pairs for which del returns true.
// Items are ranged in order.
func (m *SliceMap[K, V]) DeleteFunc(del func(K, V) bool) {
	for i, k := range slices.Backward(m.order) {
		if del(k, m.items[k]) {
			m.items.Delete(k)
			m.order = slices.Delete(m.order, i, i+1)
		}
	}
}

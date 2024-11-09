package maps

import "iter"

// MapI is the interface used by all the Map types.
type MapI[K comparable, V any] interface {
	Setter[K, V]
	Getter[K, V]
	Loader[K, V]
	Clear()
	Len() int
	Range(func(k K, v V) bool)
	Has(k K) bool
	Keys() []K
	Values() []V
	Merge(MapI[K, V])
	Equal(MapI[K, V]) bool
	Delete(k K)
	All() iter.Seq2[K, V]
	KeysIter() iter.Seq[K]
	ValuesIter() iter.Seq[V]
	Insert(seq iter.Seq2[K, V])
	DeleteFunc(del func(K, V) bool)
	String() string
}

// Setter sets a value in a map.
type Setter[K comparable, V any] interface {
	Set(K, V)
}

// Getter gets a value from a map.
type Getter[K comparable, V any] interface {
	Get(k K) (v V)
}

// Loader loads a value from a map.
type Loader[K comparable, V any] interface {
	Load(k K) (v V, ok bool)
}

// EqualFunc returns true if all the keys and values of the m1 and m2 are equal.
//
// The function eq is called on the values to determine equality. Keys are compared using ==.
// If one of the maps is a "safe" map, its more efficient to pass that map as m2.
func EqualFunc[K comparable, V1, V2 any](m1 MapI[K, V1], m2 MapI[K, V2], eq func(V1, V2) bool) bool {
	if m1.Len() != m2.Len() {
		return false
	}
	ret := true
	m2.Range(func(k K, v V2) bool {
		if !m1.Has(k) || !eq(m1.Get(k), v) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

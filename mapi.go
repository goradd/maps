package maps

// MapI is the interface used by all the Map types.
type MapI[K comparable, V any] interface {
	Clear()
	Len() int
	Range(func(k K, v V) bool)
	Load(k K) (v V, ok bool)
	Get(k K) (v V)
	Has(k K) bool
	Keys() []K
	Values() []V
	Set(K, V)
	Merge(MapI[K, V])
	Equal(MapI[K, V]) bool
}

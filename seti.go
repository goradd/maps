package maps

import "iter"

// SetI is the interface used by all the Set types.
type SetI[K comparable] interface {
	Add(k ...K) SetI[K]
	Clear()
	Len() int
	Range(func(k K) bool)
	Has(k K) bool
	Values() []K
	Merge(SetI[K])
	Equal(SetI[K]) bool
	Delete(k K)
	All() iter.Seq[K]
	Insert(seq iter.Seq[K])
	Clone() *Set[K]
	DeleteFunc(del func(K) bool)
}

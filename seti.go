package maps

import "iter"

// SetI is the interface used by all the Set types.
type SetI[K comparable] interface {
	Add(k ...K) SetI[K]
	Clear()
	Len() int
	Copy(in SetI[K])
	Range(func(k K) bool)
	Has(k K) bool
	Values() []K
	// Deprecated: use Copy instead
	Merge(SetI[K])
	Equal(SetI[K]) bool
	Delete(k K)
	All() iter.Seq[K]
	Insert(seq iter.Seq[K])
	DeleteFunc(del func(K) bool)
}

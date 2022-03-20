package maps

// Equaler is the interface that implements an Equal function. If your Map has
// non-comparible values, like a slice, but you would still like to call Equal() on that
// map, define an Equal function to do the comparison.
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

func maker[M any, K comparable, V any]() MapI[K, V] {
	var i any
	i = new(M)
	return i.(MapI[K, V])
}

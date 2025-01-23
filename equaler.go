package maps

// Equaler is the interface that implements an Equal function and that provides a way for the
// various MapI like objects to determine if they are equal.
//
// In particular, if your Map has
// non-comparable values, like a slice, but you would still like to call Equal() on that
// map, define an Equal function on the values to do the comparison. For example:
//
//	type mySlice []int
//
//	func (s mySlice) Equal(b any) bool {
//		if s2, ok := b.(mySlice); ok {
//			if len(s) == len(s2) {
//				for i, v := range s2 {
//					if s[i] != v {
//						return false
//					}
//				}
//				return true
//			}
//		}
//		return false
//	}
type Equaler interface {
	Equal(a any) bool
}

func equalValues(a, b any) bool {
	if e, ok := a.(Equaler); ok {
		return e.Equal(b)
	}

	return a == b
}

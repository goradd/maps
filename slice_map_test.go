package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceMap_Mapi(t *testing.T) {
	runMapiTests[SliceMap[string, int]](t, makeMapi[SliceMap[string, int]])
}

func init() {
	gob.Register(new(SliceMap[string, int]))
}

func TestSliceMap_MapiWithSortFunction(t *testing.T) {
	runMapiTests[SliceMap[string, int]](t,
		func(sources ...mapT) MapI[string, int] {
			m := new(SliceMap[string, int])
			for _, s := range sources {
				m.Merge(s)
			}
			m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
				return k1 < k2 // mimic natural sort for test result comparisons
			})
			return m
		},
	)
}

func ExampleSliceMap_SetSortFunc() {
	m := new(SliceMap[string, int])

	m.Set("b", 2)
	m.Set("a", 1)

	// This will print in the order items were assigned
	fmt.Println(m)

	// sort by keys
	m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
		return k1 < k2
	})
	fmt.Println(m)

	// Output: {"b":2,"a":1}
	// {"a":1,"b":2}
}

func ExampleSliceMap_SetAt() {
	m := new(SliceMap[string, int])
	m.Set("b", 2)
	m.Set("a", 1)
	m.SetAt(1, "c", 3)
	fmt.Println(m)
	// Output: {"b":2,"c":3,"a":1}
}

func ExampleSliceMap_GetAt() {
	m := new(SliceMap[string, int])
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 1)
	v := m.GetAt(1)
	fmt.Print(v)
	// Output: 3
}

func ExampleSliceMap_GetKeyAt() {
	m := new(SliceMap[string, int])
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 1)
	v := m.GetKeyAt(1)
	fmt.Print(v)
	// Output: c
}

func TestSliceMap_SetAt(t *testing.T) {
	m := new(SliceMap[string, int])
	m.Set("b", 2)
	m.Set("a", 1)
	m.SetAt(5, "c", 3)
	assert.Equal(t, 3, m.GetAt(2))

	// backwards from end
	m.SetAt(-1, "d", 4)
	assert.Equal(t, 4, m.GetAt(2))

	// past start, so at beginning
	m.SetAt(-7, "e", 5)
	assert.Equal(t, 5, m.GetAt(0))

	// set same item does not cause reshuffle
	m.Set("e", 6)
	assert.Equal(t, 6, m.GetAt(0))

	// delete and set will put to end
	m.Delete("e")
	m.Set("e", 6)
	assert.Equal(t, 6, m.GetAt(4))

	// Or force it to new location
	m.SetAt(3, "e", 6)
	assert.Equal(t, 6, m.GetAt(3))

	// can't call SetAt when there is a sort function
	m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
		return k1 < k2
	})
	assert.Panics(t, func() {
		m.SetAt(3, "f", 4)
	})
}

func TestSliceMap_GetAt(t *testing.T) {
	// Just need to check operation of empty maps. Other GetAt uses are checked elsewhere.
	m := new(SliceMap[string, int])
	assert.Equal(t, 0, m.GetAt(0))
	assert.Equal(t, "", m.GetKeyAt(0))
}

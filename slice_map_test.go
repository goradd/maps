package maps

import (
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 6, m.GetAt(m.Len()-1))

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

func TestSliceMap_NilMap(t *testing.T) {
	var m *SliceMap[string, int]
	assert.Equal(t, 0, m.GetAt(0))
	assert.Equal(t, "", m.GetKeyAt(0))
	assert.Equal(t, 0, m.Get("m"))
	v, ok := m.Load("m")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)
	assert.False(t, m.Has("n"))
	assert.Nil(t, m.Values())
	assert.Nil(t, m.Keys())

	assert.Equal(t, m.Len(), 0)
	assert.Panics(t, func() {
		m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool { return false })
	})
	assert.Panics(t, func() {
		m.Set("m", 1)
	})
	assert.NotPanics(t, func() {
		m.Delete("m")
	})
	assert.NotPanics(t, func() {
		m.Delete("m")
	})

	b, err := m.MarshalBinary()
	assert.Nil(t, b)
	assert.Nil(t, err)

	b2, err2 := m.MarshalJSON()
	assert.Nil(t, b2)
	assert.Nil(t, err2)

	assert.Panics(t, func() {
		_ = m.UnmarshalBinary(nil)
	})
	assert.Panics(t, func() {
		_ = m.UnmarshalJSON(nil)
	})

	assert.True(t, m.Equal(nil))
	assert.NotPanics(t, func() {
		m.Clear()
	})
	assert.Equal(t, "", m.String())
}

func TestSliceMap_Clone(t *testing.T) {
	// Create a new SafeSliceMap and populate it
	originalMap := NewSliceMap[string, int]()
	originalMap.Set("b", 2)
	originalMap.Set("a", 1)
	originalMap.Set("c", 3)

	// Clone the original map
	clonedMap := originalMap.Clone()

	assert.True(t, clonedMap.Equal(originalMap))

	// Verify that modifying the cloned map does not affect the original map
	clonedMap.Set("a", 100)
	if originalMap.Get("a") == 100 {
		t.Error("Modification in cloned map affected the original map")
	}

	// Verify that the original map remains unchanged
	if originalMap.Get("a") != 1 {
		t.Errorf("Expected value for key 'a' in original map to be %d, got %d", 1, originalMap.Get("a"))
	}

	// Verify order is the same
	values := clonedMap.Values()
	expectedValues := []int{2, 100, 3}
	assert.Equal(t, expectedValues, values)
}

func TestCollectSliceMap(t *testing.T) {
	// Create a sequence of key-value pairs
	s := NewSliceMap[string, int]()
	s.Set("b", 2)
	s.Set("a", 1)
	s.Set("c", 3)
	seq := s.All()
	// Use CollectSafeSliceMap to create a new SafeSliceMap
	collectedMap := CollectSafeSliceMap(seq)

	assert.True(t, s.Equal(collectedMap))

	// Ensure the order of keys follows the insertion order
	keys := collectedMap.Keys()
	expectedKeys := []string{"b", "a", "c"}
	assert.Equal(t, keys, expectedKeys)
}

package maps

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

type mapT = StdMap[string, int]
type mapTI = MapI[string, int]

func TestNewStdMap(t *testing.T) {
	m := NewStdMap(map[string]int{"a": 1})
	assert.Equal(t, 1, m.Get("a"))
}

func init() {
	gob.Register(StdMap[string, int]{})
}

func TestMap_Clear(t *testing.T) {
	var mNil mapT
	tests := []struct {
		name string
		m    mapTI
	}{
		{"zero", mNil},
		{"empty", NewStdMap[string, int]()},
		{"1 item", mapT{"a": 1}},
		{"2 items", mapT{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
			if tt.m.Len() != 0 {
				t.Errorf("StdMap not cleared: %q", tt.m)
			}
		})
	}
}

func TestMap_Merge(t *testing.T) {
	tests := []struct {
		name     string
		m1       mapT
		m2       mapT
		expected mapT
	}{
		{"1 to 1", mapT{"a": 1}, mapT{"b": 2}, mapT{"a": 1, "b": 2}},
		{"overwrite", mapT{"a": 1}, mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"to empty", mapT{}, mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"from empty", mapT{"a": 1, "b": 2}, mapT{}, mapT{"a": 1, "b": 2}},
		{"from cast map", mapT{"a": 1}, Cast(map[string]int{"b": 2}), mapT{"a": 1, "b": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m1.Merge(tt.m2)
			if !tt.m1.Equal(tt.expected) {
				t.Errorf("Merge error. Expected: %q, got %q", tt.expected, tt.m1)
			}
		})
	}
}

func TestMap_MergePanic(t *testing.T) {
	assert.Panics(t, func() {
		var m mapT
		m.Merge(mapT{"a": 1})
	})
}

func TestMap_GetHasLoad(t *testing.T) {
	m := mapT{"a": 1, "b": 2}
	if m.Has("c") {
		t.Errorf("Expected false, got true")
	}
	if !m.Has("b") {
		t.Errorf("Expected true, got false")
	}
	if v := m.Get("b"); v != 2 {
		t.Errorf("Expected 2, got %q", v)
	}
	if v := m.Get("c"); v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}

	v, ok := m.Load("a")
	if v != 1 {
		t.Errorf("Expected 1, got %q", v)
	}
	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestMap_NilGetHasLoad(t *testing.T) {
	var m mapT
	if m.Has("c") {
		t.Errorf("Expected false, got true")
	}
	if v := m.Get("c"); v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}

	v, ok := m.Load("a")
	if v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestMap_Range(t *testing.T) {
	var mNil mapT

	tests := []struct {
		name     string
		m        mapTI
		expected int
	}{
		{"nil", mNil, 0},
		{"1", mapT{"a": 1}, 1},
		{"2", mapT{"a": 1, "b": 2}, 2},
		{"3", mapT{"a": 1, "b": 2, "c": 3}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := 0
			tt.m.Range(func(k string, v int) bool {
				count++
				if count > 1 {
					return false
				}
				return true
			})
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func ExampleCast() {
	m := map[string]int{"a": 1}
	b := Cast(m)
	fmt.Print(b.Len())
	//Output: 1
}

func TestStdMap_Set(t *testing.T) {
	var a mapT

	assert.Panics(t, func() {
		a.Set("a", 1)
	})

	a = mapT{}
	a.Set("a", 1)
	a.Set("b", 2)
	assert.Equal(t, 2, a.Get("b"))
}

func TestStdMap_Keys(t *testing.T) {
	m := NewStdMap(mapT{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "c")

	m = NewStdMap(mapT{})
	assert.Nil(t, m.Keys())

}

func TestStdMap_Values(t *testing.T) {
	m := NewStdMap(mapT{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	assert.Len(t, values, 3)
	assert.Contains(t, values, 3)

	m = NewStdMap(mapT{})
	assert.Nil(t, m.Values())
}

func TestStdMap_Equal(t *testing.T) {
	tests := []struct {
		name string
		m    StdMap[string, int]
		m2   mapT
		want bool
	}{
		{"equal", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 1}, true},
		{"empty", StdMap[string, int]{}, StdMap[string, int]{}, true},
		{"dif len", StdMap[string, int]{"a": 1}, StdMap[string, int]{}, false},
		{"dif len 1", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 1, "b": 2}, false},
		{"dif value", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 2}, false},
		{"dif key", StdMap[string, int]{"a": 1}, StdMap[string, int]{"b": 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Equal(tt.m2), "Equal(%v)", tt.m2)
		})
	}
}

type mySlice []int

func (s mySlice) Equal(b any) bool {
	if s2, ok := b.(mySlice); ok {
		if len(s) == len(s2) {
			for i, v := range s2 {
				if s[i] != v {
					return false
				}
			}
			return true
		}
	}
	return false
}

func TestEqualValues(t *testing.T) {
	a := 1
	b := 1
	assert.True(t, equalValues(a, b))
	b = 2
	assert.False(t, equalValues(a, b))

	c := mySlice{1, 2}
	d := mySlice{1, 2}
	assert.True(t, equalValues(c, d))

	e := []float32{1, 2}
	f := []float32{1, 2}
	assert.Panics(t, func() { equalValues(e, f) })
}

func TestMarshalBinary(t *testing.T) {
	m := StdMap[string, int]{"a": 1, "b": 2}

	// Marshal the map
	data, err := m.MarshalBinary()
	if err != nil {
		t.Fatalf("Error marshalling: %v", err)
	}

	// Unmarshal the data
	var m2 StdMap[string, int]
	err = m2.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("Error unmarshalling: %v", err)
	}

	// Compare the original and unmarshalled maps
	assert.Equal(t, m, m2)
}

func TestMarshalJSON(t *testing.T) {
	m := StdMap[string, int]{"a": 1, "b": 2}

	// Marshal the map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Error marshalling to JSON: %v", err)
	}

	// Assert the JSON output
	expectedJSON := `{"a":1,"b":2}`
	assert.Equal(t, expectedJSON, string(jsonData))
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{"a":1,"b":2}`)

	// Unmarshal the JSON into a StdMap
	var m StdMap[string, int]
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		t.Fatalf("Error unmarshalling from JSON: %v", err)
	}

	// Assert the unmarshalled map
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, 2, m["b"])
}

func TestUnmarshalJSONInvalidInput(t *testing.T) {
	invalidJSON := []byte(`invalid json`)

	// Unmarshal the invalid JSON
	var m StdMap[string, int]
	err := json.Unmarshal(invalidJSON, &m)
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	m := StdMap[string, int]{"a": 1, "b": 2}

	// Delete an existing key
	m.Delete("a")
	_, ok := m["a"]
	assert.False(t, ok)

	// Delete a non-existent key
	m.Delete("c")
	// No error should occur, and the map should remain unchanged
	assert.Equal(t, 2, m["b"])
}

func TestString(t *testing.T) {
	m := StdMap[string, int]{"a": 1, "b": 2}

	// Get the string representation
	str := m.String()

	// Check the string representation
	expected := `{"a":1, "b":2}`
	assert.Equal(t, expected, str)
}

func ExampleStdMap_All() {
	m := StdMap[string, int]{"a": 1, "b": 2, "c": 3}

	var actualKeys []string
	var actualValues []int

	for k, v := range m.All() {
		actualKeys = append(actualKeys, k)
		actualValues = append(actualValues, v)
	}
	slices.Sort(actualKeys)
	slices.Sort(actualValues)
	fmt.Println(actualKeys)
	fmt.Println(actualValues)

	// Output: [a b c]
	// [1 2 3]
}

func ExampleStdMap_KeysIter() {
	m := StdMap[string, int]{"a": 1, "b": 2, "c": 3}

	var actualKeys []string

	for k := range m.KeysIter() {
		actualKeys = append(actualKeys, k)
	}
	slices.Sort(actualKeys)
	fmt.Println(actualKeys)

	// Output: [a b c]
}

func ExampleStdMap_ValuesIter() {
	m := StdMap[string, int]{"a": 1, "b": 2, "c": 3}

	var actualValues []int

	for v := range m.ValuesIter() {
		actualValues = append(actualValues, v)
	}
	slices.Sort(actualValues)
	fmt.Println(actualValues)

	// Output: [1 2 3]
}

func TestStdMap_Insert(t *testing.T) {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m2 := StdMap[string, int]{"a": 1}
	m2.Insert(m1.All())
	assert.True(t, m1.Equal(m2))
}

func TestStdMap_Collect(t *testing.T) {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m2 := CollectStdMap(m1.All())
	assert.True(t, m1.Equal(m2))
}

func TestStdMap_Clone(t *testing.T) {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m2 := m1.Clone()
	assert.True(t, m1.Equal(m2))
}

func ExampleStdMap_DeleteFunc() {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m1.DeleteFunc(func(k string, v int) bool {
		return v != 2
	})
	fmt.Println(m1.String())
	// Output: {"b":2}
}

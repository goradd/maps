package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

type sliceSetT = SliceSet[string]
type sliceSetTI = SetI[string]

func TestSliceSet_SetI(t *testing.T) {
	runSetITests[sliceSetT](t, makeSetI[sliceSetT])
}

func init() {
	gob.Register(new(sliceSetT))
}

func TestSliceSet_Values(t *testing.T) {
	type testCase[K comparable] struct {
		name string
		m    *SliceSet[K]
		want []K
	}
	tests := []testCase[int]{
		{"none", NewSliceSet[int](), []int(nil)},
		{"one", NewSliceSet[int](1), []int{1}},
		{"three", NewSliceSet[int](1, 2, 3), []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Values(), "Values()")
		})
	}
}

func TestSliceSet_MarshalJSON(t *testing.T) {
	type testCase[K comparable] struct {
		name    string
		m       *SliceSet[K]
		wantOut string
		wantErr bool
	}
	tests := []testCase[string]{
		{"zero", NewSliceSet[string](), `[]`, false},
		{"one", NewSliceSet("a"), `["a"]`, false},
		{"three", NewSliceSet("a", "c", "b"), `["a","c","b"]`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.m.MarshalJSON()
			gotOut := string(b)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.wantOut, gotOut, "MarshalJSON()")
		})
	}
}

func TestSliceSet_All(t *testing.T) {
	set := NewSliceSet[int]()
	set.Add(5)
	set.Add(3)
	set.Add(8)
	set.Add(1)

	iterator := set.All()
	var result []int

	for v := range iterator {
		result = append(result, v)
	}

	expected := []int{5, 3, 8, 1}
	assert.Equal(t, expected, result)
}

func TestSliceSet_Range(t *testing.T) {
	type testCase[K comparable] struct {
		name     string
		m        *SliceSet[K]
		expected []int
	}
	tests := []testCase[int]{
		{"none", NewSliceSet[int](), []int(nil)},
		{"one", NewSliceSet[int](1), []int{1}},
		{"three", NewSliceSet[int](1, 2, 3), []int{1, 2, 3}},
		{"four", NewSliceSet[int](4, 3, 2, 1), []int{4, 3, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var values []int
			var count int
			tt.m.Range(func(i int) bool {
				values = append(values, i)
				count++
				return count < 3
			})
			assert.Equal(t, tt.expected, values)
		})
	}
}

func TestSliceSet_Clone(t *testing.T) {
	m1 := NewSliceSet("a", "b", "c")
	m2 := m1.Clone()
	assert.True(t, m1.Equal(m2))

	var m3 *SliceSet[string]
	m4 := m3.Clone()
	m3.Equal(m4)
	assert.True(t, m3.Equal(m4))

	m2.Add("d")
	assert.False(t, m1.Equal(m2))
}

func TestSliceSet_Nil(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var m1, m2 *SliceSet[string]

		assert.Equal(t, 0, m1.Len())
		m1.Clear()
		assert.True(t, m1.Equal(m2))
		m3 := m2.Clone()
		assert.True(t, m1.Equal(m3))
		m3.Add("a")
		assert.False(t, m1.Equal(m3))
		m1.Range(func(k string) bool {
			assert.Fail(t, "no range should happen")
			return false
		})
		assert.False(t, m1.Has("b"))
		m1.Delete("a")
		assert.Empty(t, m1.Values())
		assert.Equal(t, "{}", m1.String())
		m1.DeleteFunc(func(k string) bool {
			return false
		})
		for _ = range m1.All() {
			assert.Fail(t, "no range should happen")
		}
		assert.Panics(t, func() {
			m1.Insert(slices.Values([]string{"a"}))
		})
		assert.Panics(t, func() {
			m1.Add("a")
		})
		assert.Panics(t, func() {
			m1.Copy(m2)
		})
	})
}

func ExampleSliceSet_String() {
	m := NewSliceSet("a", "c", "a", "b")
	fmt.Print(m.String())
	// Output: {"a","c","b"}
}

func TestSliceSet_SetSortFunc(t *testing.T) {
	m := NewSliceSet[string]("b", "d", "c")
	assert.Equal(t, []string{"b", "d", "c"}, m.Values())
	m.SetSortFunc(func(k1, k2 string) bool {
		return k1 < k2
	})
	assert.Equal(t, []string{"b", "c", "d"}, m.Values())
	m.Add("a")
	assert.Equal(t, []string{"a", "b", "c", "d"}, m.Values())

	m.SetSortFunc(nil)
	m.Add("aa")
	assert.Equal(t, []string{"a", "b", "c", "d", "aa"}, m.Values())

	var m2 *SliceSet[string]
	assert.Panics(t, func() {
		m2.SetSortFunc(func(k1, k2 string) bool {
			return k1 < k2
		})
	})
}

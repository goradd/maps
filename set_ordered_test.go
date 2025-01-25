package maps

import (
	"cmp"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

type orderedSetT = OrderedSet[string]
type orderedSetTI = SetI[string]

func TestOrderedSet_SetI(t *testing.T) {
	runSetITests[orderedSetT](t, makeSetI[orderedSetT])
}

func init() {
	gob.Register(new(orderedSetT))
}

func TestOrderedSet_Values(t *testing.T) {
	type testCase[K cmp.Ordered] struct {
		name string
		m    *OrderedSet[K]
		want []K
	}
	tests := []testCase[int]{
		{"none", NewOrderedSet[int](), []int(nil)},
		{"one", NewOrderedSet[int](1), []int{1}},
		{"three", NewOrderedSet[int](1, 2, 3), []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Values(), "Values()")
		})
	}
}

func TestOrderedSet_MarshalJSON(t *testing.T) {
	type testCase[K cmp.Ordered] struct {
		name    string
		m       *OrderedSet[K]
		wantOut string
		wantErr bool
	}
	tests := []testCase[string]{
		{"zero", NewOrderedSet[string](), `[]`, false},
		{"one", NewOrderedSet("a"), `["a"]`, false},
		{"three", NewOrderedSet("a", "c", "b"), `["a","b","c"]`, false},
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

func TestOrderedSet_All(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Add(5)
	set.Add(3)
	set.Add(8)
	set.Add(1)

	iterator := set.All()
	var result []int

	for v := range iterator {
		result = append(result, v)
	}

	expected := []int{1, 3, 5, 8}
	assert.Equal(t, expected, result)
}

func TestOrderedSet_Range(t *testing.T) {
	type testCase[K cmp.Ordered] struct {
		name     string
		m        *OrderedSet[K]
		expected []int
	}
	tests := []testCase[int]{
		{"none", NewOrderedSet[int](), []int(nil)},
		{"one", NewOrderedSet[int](1), []int{1}},
		{"three", NewOrderedSet[int](1, 2, 3), []int{1, 2, 3}},
		{"four", NewOrderedSet[int](4, 3, 2, 1), []int{1, 2, 3}},
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

func TestOrderedSet_Clone(t *testing.T) {
	m1 := NewOrderedSet("a", "b", "c")
	m2 := m1.Clone()
	assert.True(t, m1.Equal(m2))

	var m3 *OrderedSet[string]
	m4 := m3.Clone()
	m3.Equal(m4)
	assert.True(t, m3.Equal(m4))

	m2.Add("d")
	assert.False(t, m1.Equal(m2))
}

func TestOrderedSet_Nil(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var m1, m2 *OrderedSet[string]

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

func ExampleOrderedSet_String() {
	m := NewOrderedSet("a", "c", "a", "b")
	fmt.Print(m.String())
	// Output: {"a","b","c"}
}

package maps

import (
	"cmp"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"testing"
)

type orderedSetT = OrderedSet[string]
type orderedSetTI = SetI[string]

func TestOrderedSet_SetI(t *testing.T) {
	runSetITests[OrderedSet[string]](t, makeSetI[OrderedSet[string]])
}

func init() {
	gob.Register(new(OrderedSet[string]))
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

func TestOrderedSetAll(t *testing.T) {
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

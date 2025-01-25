package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

type setT = Set[string]
type setTI = SetI[string]

func TestSet_SetI(t *testing.T) {
	runSetITests[Set[string]](t, makeSetI[Set[string]])
}

func init() {
	gob.Register(new(Set[string]))
}

func ExampleSet_String() {
	m := new(Set[string])
	m.Add("a")
	fmt.Print(m.String())
	// Output: {"a"}
}

func TestCollectSet(t *testing.T) {
	m1 := NewSet("a", "b", "c")
	m2 := CollectSet(m1.All())
	fmt.Println(m2.String())
	assert.True(t, m1.Equal(m2))
}

func TestSet_Clone(t *testing.T) {
	m1 := NewSet("a", "b", "c")
	m2 := CollectSet(m1.All())
	m3 := m2.Clone()
	assert.True(t, m1.Equal(m3))
}

func TestSet_Nil(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var m1, m2 *Set[string]

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

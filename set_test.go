package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
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
	m.Add("b")
	m.Add("a")
	v := m.Values()
	sort.Strings(v)
	fmt.Print(v)
	// Output: [a b]
}

func ExampleCollectSet() {
	m1 := NewSet("a", "b", "c")
	m2 := CollectSet(m1.All())
	fmt.Println(m2.String())
	// Output: {"a","b","c"}
}

func TestSet_Clone(t *testing.T) {
	m1 := NewSet("a", "b", "c")
	m2 := CollectSet(m1.All())
	m3 := m2.Clone()
	assert.True(t, m1.Equal(m3))
}

package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap_Mapi(t *testing.T) {
	runMapiTests[Map[string, int]](t, makeMapi[Map[string, int]])
}

func init() {
	gob.Register(new(Map[string, int]))
}

func ExampleMap_String() {
	m := new(Map[string, int])
	m.Set("a", 1)
	m.Set("b", 2)
	fmt.Print(m)
	// Output: {"a":1, "b":2}
}

func ExampleCollectMap() {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m2 := CollectMap(m1.All())
	fmt.Println(m2.String())
	// Output: {"a":1, "b":2, "c":3}
}

func TestMap_Clone(t *testing.T) {
	m1 := StdMap[string, int]{"a": 1, "b": 2, "c": 3}
	m2 := CollectMap(m1.All())
	m3 := m2.Clone()
	assert.True(t, m1.Equal(m3))
}

package maps

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafeMap_Mapi(t *testing.T) {
	runMapiTests[SafeMap[string, int]](t, makeMapi[SafeMap[string, int]])
}

func init() {
	gob.Register(new(SafeMap[string, int]))
}

func TestNil(t *testing.T) {
	var m SafeMap[string, int]

	assert.False(t, m.Has("z"))

	a, ok := m.Load("a")
	assert.Empty(t, a)
	assert.False(t, ok)

	m.Delete("a")

	assert.Nil(t, m.Values())
	assert.Nil(t, m.Keys())

}

func ExampleSafeMap_String() {
	m := new(SafeMap[string, int])
	m.Set("a", 1)
	m.Set("b", 2)
	fmt.Print(m)
	// Output: {"a":1, "b":2}
}

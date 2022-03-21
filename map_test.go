package maps

import (
	"fmt"
	"testing"
)

func TestMap_Mapi(t *testing.T) {
	runMapiTests[Map[string, int]](t, makeMapi[Map[string, int]])
}

func ExampleMap_String() {
	m := new(Map[string, int])
	m.Set("a", 1)
	m.Set("b", 2)
	fmt.Print(m)
	// Output: {"a":1, "b":2}
}

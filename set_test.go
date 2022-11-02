package maps

import (
	"encoding/gob"
	"fmt"
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
	fmt.Print(m)
	// Output: {"a","b"}
}

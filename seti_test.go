package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

type makeSetF func(sources ...string) SetI[string]

func makeSetI[M any](sources ...string) SetI[string] {
	var m any
	m = new(M)
	i := m.(SetI[string])
	for _, s := range sources {
		i.Add(s)
	}
	return i
}

func runSetITests[M any](t *testing.T, f makeSetF) {
	testSetClear(t, f)
	testSetLen(t, f)
	testSetMerge(t, f)
	testSetHas(t, f)
	testSetRange(t, f)
	testSetAdd(t, f)
	testSetValues(t, f)
	testSetEqual(t, f)
	testSetBinaryMarshal[M](t, f)
	testSetMarshalJSON(t, f)
	testSetUnmarshalJSON[M](t, f)
	testSetDelete(t, f)
	testSetAll(t, f)
	testSetInsert(t, f)
	testSetDeleteFunc(t, f)
	testSetCopy(t, f)
}

func testSetClear(t *testing.T, f makeSetF) {
	tests := []struct {
		name string
		m    setTI
	}{
		{"empty", f()},
		{"1 item", f("a")},
		{"2 items", f("a", "b")},
	}

	for _, tt := range tests {
		t.Run("Clear "+tt.name, func(t *testing.T) {
			tt.m.Clear()
			if tt.m.Len() != 0 {
				t.Errorf("SetI not cleared")
			}
		})
	}
}

func testSetLen(t *testing.T, f makeSetF) {
	assert.Equal(t, 0, f().Len())
	assert.Equal(t, 2, f("a", "b").Len())
}

func testSetMerge(t *testing.T, f makeSetF) {
	tests := []struct {
		name     string
		m1       setTI
		m2       setTI
		expected setTI
	}{
		{"1 to 1", f("a"), f("b"), f("a", "b")},
		{"overwrite", f("a"), f("a", "b"), f("a", "b")},
		{"to empty", f(), f("a", "b"), f("a", "b")},
		{"from nil", f("a"), nil, f("a")},
		{"from empty", f("a"), f(), f("a")},
	}
	for _, tt := range tests {
		t.Run("Merge "+tt.name, func(t *testing.T) {
			tt.m1.Merge(tt.m2)
			if !tt.m1.Equal(tt.expected) {
				t.Errorf("Merge error. Expected: %q, got %q", tt.expected, tt.m1)
			}
		})
	}
}

func testSetHas(t *testing.T, f makeSetF) {

	m := f("a", "b")

	t.Run("Has", func(t *testing.T) {
		if m.Has("c") {
			t.Errorf("Expected false, got true")
		}
		if !m.Has("b") {
			t.Errorf("Expected true, got false")
		}
	})
}

func testSetRange(t *testing.T, f makeSetF) {
	tests := []struct {
		name        string
		m           setTI
		expectedLen int
	}{
		{"none", f(), 0},
		{"one", f("a"), 1},
		{"three", f("b", "a", "c"), 3},
		{"four", f("d", "a", "c", "b"), 3},
	}
	for _, tt := range tests {
		t.Run("Range "+tt.name, func(t *testing.T) {
			var values []string
			var count int
			tt.m.Range(func(i string) bool {
				values = append(values, i)
				count++
				return count < 3
			})
			assert.Equal(t, tt.expectedLen, len(values))
		})
	}
}

func testSetAdd(t *testing.T, f makeSetF) {
	t.Run("Set", func(t *testing.T) {
		a := f()
		a.Add("a")
		a.Add("b")
		a.Add("b")
		assert.True(t, a.Has("b"))
	})
}

func testSetValues(t *testing.T, f makeSetF) {
	t.Run("Values", func(t *testing.T) {
		m := f("a", "b", "c")
		values := m.Values()
		assert.Len(t, values, 3)
		assert.Contains(t, values, "c")
	})
}

func testSetEqual(t *testing.T, f makeSetF) {
	tests := []struct {
		name string
		m    setTI
		m2   setTI
		want bool
	}{
		{"equal", f("a"), f("a"), true},
		{"empty", f(), f(), true},
		{"dif len", f("a"), f(), false},
		{"dif len 1", f("a"), f("a", "b"), false},
		{"dif value", f("a"), f("b"), false},
	}
	for _, tt := range tests {
		t.Run("Equal "+tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Equal(tt.m2), "Equal(%v)", tt.m2)
		})
	}
}

func testSetBinaryMarshal[M any](t *testing.T, f makeSetF) {
	t.Run("BinaryMarshal", func(t *testing.T) {
		// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding
		m := f("a", "b", "c")
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf) // Will write
		dec := gob.NewDecoder(&buf) // Will read

		var i any
		i = m

		err := enc.Encode(&i)
		assert.NoError(t, err)

		var i2 any
		err = dec.Decode(&i2)
		assert.NoError(t, err)

		m2 := i2.(SetI[string])
		assert.True(t, m2.Has("a"))
		assert.True(t, m2.Has("c"))
	})
}

func testSetMarshalJSON(t *testing.T, f makeSetF) {
	t.Run("MarshalJSON", func(t *testing.T) {
		m := f("a", "b", "c")
		s, err := json.Marshal(m)
		assert.NoError(t, err)
		// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
		assert.Contains(t, string(s), `"a"`)

		m = f()
		s, err = json.Marshal(m)
		assert.NoError(t, err)
		// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
		assert.Equal(t, "[]", string(s))
	})
}

func testSetUnmarshalJSON[M any](t *testing.T, f makeSetF) {
	b := []byte(`["a","b","c"]`)
	var m M

	json.Unmarshal(b, &m)
	var i interface{}
	i = &m
	m2 := i.(SetI[string])

	assert.True(t, m2.Has("c"))

	b = []byte(`[]`)

	var m3 M

	json.Unmarshal(b, &m3)
	i = &m3
	m4 := i.(SetI[string])

	assert.Equal(t, 0, m4.Len())

	b = []byte(`["d"]`)

	// Unmarshalling into an existing set should add values
	json.Unmarshal(b, &m)
	i = &m
	m5 := i.(SetI[string])

	assert.Equal(t, 4, m5.Len())

}

func testSetDelete(t *testing.T, f makeSetF) {
	t.Run("Delete", func(t *testing.T) {
		m := f("a", "b")
		m.Delete("a")

		assert.False(t, m.Has("a"))
		assert.True(t, m.Has("b"))

		m.Delete("b")
		assert.False(t, m.Has("b"))

		m.Delete("b") // make sure deleting from an empty map is a no-op
	})
}

func testSetAll(t *testing.T, f makeSetF) {
	t.Run("All", func(t *testing.T) {
		m := f("a", "b", "c")

		var actualValues []string

		for k := range m.All() {
			actualValues = append(actualValues, k)
		}
		slices.Sort(actualValues)

		assert.Equal(t, []string{"a", "b", "c"}, actualValues)
	})
}

func testSetInsert(t *testing.T, f makeSetF) {
	t.Run("Insert", func(t *testing.T) {
		m1 := f("a", "b", "c")
		m2 := f("a")
		m2.Insert(m1.All())
		assert.True(t, m1.Equal(m2))
	})
}

func testSetDeleteFunc(t *testing.T, f makeSetF) {
	t.Run("DeleteFunc", func(t *testing.T) {
		m1 := f("a", "b", "c")
		m1.DeleteFunc(func(k string) bool {
			return k != "b"
		})
		assert.Equal(t, 1, m1.Len())
	})
}

func testSetCopy(t *testing.T, f makeSetF) {
	t.Run("DeleteFunc", func(t *testing.T) {
		m1 := f("a", "b", "c")
		m2 := f()
		m2.Copy(m1)
		assert.True(t, m1.Equal(m2))
	})
}

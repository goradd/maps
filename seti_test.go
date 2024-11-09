package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
		name     string
		m        setTI
		expected int
	}{
		{"0", f(), 0},
		{"1", f("a"), 1},
		{"2", f("a", "b"), 2},
		{"3", f("a", "b", "c"), 2},
	}
	for _, tt := range tests {
		t.Run("Range "+tt.name, func(t *testing.T) {
			count := 0
			tt.m.Range(func(k string) bool {
				count++
				if count > 1 {
					return false
				}
				return true
			})
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
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

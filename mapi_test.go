package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type makeF func(sources ...mapT) MapI[string, int]

func makeMapi[M any](sources ...mapT) MapI[string, int] {
	var m any
	m = new(M)
	i := m.(MapI[string, int])
	for _, s := range sources {
		i.Merge(s)
	}
	return i
}

func runMapiTests[M any](t *testing.T, f makeF) {
	testClear(t, f)
	testLen(t, f)
	testMerge(t, f)
	testGetHasLoad(t, f)
	testRange(t, f)
	testSet(t, f)
	testKeys(t, f)
	testValues(t, f)
	testEqual(t, f)
	testBinaryMarshal[M](t, f)
	testMarshalJSON(t, f)
	testUnmarshalJSON[M](t, f)
	testDelete(t, f)
}

func testClear(t *testing.T, f makeF) {
	tests := []struct {
		name string
		m    mapTI
	}{
		{"empty", f()},
		{"1 item", f(mapT{"a": 1})},
		{"2 items", f(mapT{"a": 1, "b": 2})},
	}

	for _, tt := range tests {
		t.Run("Clear "+tt.name, func(t *testing.T) {
			tt.m.Clear()
			if tt.m.Len() != 0 {
				t.Errorf("MapI not cleared")
			}
		})
	}
}

func testLen(t *testing.T, f makeF) {
	assert.Equal(t, 0, f().Len())
	assert.Equal(t, 2, f(mapT{"a": 1, "b": 2}).Len())
}

func testMerge(t *testing.T, f makeF) {
	tests := []struct {
		name     string
		m1       mapTI
		m2       mapT
		expected mapT
	}{
		{"1 to 1", f(mapT{"a": 1}), mapT{"b": 2}, mapT{"a": 1, "b": 2}},
		{"overwrite", f(mapT{"a": 1}), mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"to empty", f(), mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"from nil", f(mapT{"a": 1}), nil, mapT{"a": 1}},
		{"from empty", f(mapT{"a": 1}), mapT{}, mapT{"a": 1}},
		{"from cast map", f(mapT{"a": 1}), Cast(map[string]int{"b": 2}), mapT{"a": 1, "b": 2}},
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

func testGetHasLoad(t *testing.T, f makeF) {

	m := f(mapT{"a": 1, "b": 2})

	t.Run("Has", func(t *testing.T) {
		if m.Has("c") {
			t.Errorf("Expected false, got true")
		}
		if !m.Has("b") {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Get", func(t *testing.T) {
		if v := m.Get("b"); v != 2 {
			t.Errorf("Expected 2, got %q", v)
		}
		if v := m.Get("c"); v != 0 {
			t.Errorf("Expected 0, got %q", v)
		}
	})

	t.Run("Load", func(t *testing.T) {
		v, ok := m.Load("a")
		if v != 1 {
			t.Errorf("Expected 1, got %q", v)
		}
		if !ok {
			t.Errorf("Expected true, got false")
		}
	})
}

func testRange(t *testing.T, f makeF) {
	tests := []struct {
		name     string
		m        mapTI
		expected int
	}{
		{"0", f(), 0},
		{"1", f(mapT{"a": 1}), 1},
		{"2", f(mapT{"a": 1, "b": 2}), 2},
		{"3", f(mapT{"a": 1, "b": 2, "c": 3}), 2},
	}
	for _, tt := range tests {
		t.Run("Range "+tt.name, func(t *testing.T) {
			count := 0
			tt.m.Range(func(k string, v int) bool {
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

func testSet(t *testing.T, f makeF) {
	t.Run("Set", func(t *testing.T) {
		a := f()
		a.Set("a", 1)
		a.Set("b", 2)
		assert.Equal(t, 2, a.Get("b"))
	})
}

func testKeys(t *testing.T, f makeF) {
	t.Run("Keys", func(t *testing.T) {
		m := f(mapT{"a": 1, "b": 2, "c": 3})
		keys := m.Keys()
		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "c")
	})
}

func testValues(t *testing.T, f makeF) {
	t.Run("Values", func(t *testing.T) {
		m := f(mapT{"a": 1, "b": 2, "c": 3})
		values := m.Values()
		assert.Len(t, values, 3)
		assert.Contains(t, values, 3)
	})
}

func testEqual(t *testing.T, f makeF) {
	tests := []struct {
		name string
		m    mapTI
		m2   mapT
		want bool
	}{
		{"equal", f(mapT{"a": 1}), mapT{"a": 1}, true},
		{"empty", f(), mapT{}, true},
		{"dif len", f(mapT{"a": 1}), mapT{}, false},
		{"dif len 1", f(mapT{"a": 1}), mapT{"a": 1, "b": 2}, false},
		{"dif value", f(mapT{"a": 1}), mapT{"a": 2}, false},
		{"dif key", f(mapT{"a": 1}), mapT{"b": 1}, false},
	}
	for _, tt := range tests {
		t.Run("Equal "+tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Equal(tt.m2), "Equal(%v)", tt.m2)
		})
	}
}

func testBinaryMarshal[M any](t *testing.T, f makeF) {
	t.Run("BinaryMarshal", func(t *testing.T) {
		// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding
		m := f(mapT{"a": 1, "b": 2, "c": 3})
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

		m2 := i2.(MapI[string, int])
		assert.Equal(t, 1, m2.Get("a"))
		assert.Equal(t, 3, m2.Get("c"))
	})
}

func testMarshalJSON(t *testing.T, f makeF) {
	t.Run("MarshalJSON", func(t *testing.T) {
		m := f(mapT{"a": 1, "b": 2, "c": 3})
		s, err := json.Marshal(m)
		assert.NoError(t, err)
		// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
		assert.Equal(t, `{"a":1,"b":2,"c":3}`, string(s))
	})
}

func testUnmarshalJSON[M any](t *testing.T, f makeF) {
	b := []byte(`{"a":1,"b":2,"c":3}`)
	var m M

	json.Unmarshal(b, &m)
	var i interface{}
	i = &m
	m2 := i.(MapI[string, int])

	assert.Equal(t, 3, m2.Get("c"))
}

func testDelete(t *testing.T, f makeF) {
	t.Run("Delete", func(t *testing.T) {
		m := f(mapT{"a": 1, "b": 2})
		v := m.Delete("a")

		assert.Equal(t, 1, v)
		assert.False(t, m.Has("a"))
		assert.True(t, m.Has("b"))

		v = m.Delete("b")
		assert.Equal(t, 2, v)
		assert.False(t, m.Has("b"))

		v = m.Delete("b") // make sure deleting from an empty map is a no-op
		assert.Equal(t, 0, v)
	})
}

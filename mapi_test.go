package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeMapi[M any](sources ...mapT) MapI[string, int] {
	var m any
	m = new(M)
	i := m.(MapI[string, int])
	for _, s := range sources {
		i.Merge(s)
	}
	return i
}

func runMapiTests[M any](t *testing.T) {
	testClear[M](t)
	testMerge[M](t)
	testGetHasLoad[M](t)
	testRange[M](t)
	testSet[M](t)
	testKeys[M](t)
	testValues[M](t)
	testEqual[M](t)
	testBinaryMarshal[M](t)
	testMarshalJSON[M](t)
	testUnmarshalJSON[M](t)
}

func testClear[M any](t *testing.T) {
	tests := []struct {
		name string
		m    mapTI
	}{
		{"empty", makeMapi[M]()},
		{"1 item", makeMapi[M](mapT{"a": 1})},
		{"2 items", makeMapi[M](mapT{"a": 1, "b": 2})},
	}

	for _, tt := range tests {
		t.Run("Clear "+tt.name, func(t *testing.T) {
			tt.m.Clear()
			if tt.m.Len() != 0 {
				t.Errorf("MapI not cleared: %q", tt.m)
			}
		})
	}
}

func testMerge[M any](t *testing.T) {
	tests := []struct {
		name     string
		m1       mapTI
		m2       mapT
		expected mapT
	}{
		{"1 to 1", makeMapi[M](mapT{"a": 1}), mapT{"b": 2}, mapT{"a": 1, "b": 2}},
		{"overwrite", makeMapi[M](mapT{"a": 1}), mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"to empty", makeMapi[M](), mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"from empty", makeMapi[M](mapT{"a": 1}), mapT{}, mapT{"a": 1}},
		{"from cast map", makeMapi[M](mapT{"a": 1}), Cast(map[string]int{"b": 2}), mapT{"a": 1, "b": 2}},
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

func testGetHasLoad[M any](t *testing.T) {

	m := makeMapi[M](mapT{"a": 1, "b": 2})

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

func testRange[M any](t *testing.T) {
	tests := []struct {
		name     string
		m        mapTI
		expected int
	}{
		{"1", makeMapi[M](mapT{"a": 1}), 1},
		{"2", makeMapi[M](mapT{"a": 1, "b": 2}), 2},
		{"3", makeMapi[M](mapT{"a": 1, "b": 2, "c": 3}), 2},
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

func testSet[M any](t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		a := makeMapi[M]()
		a.Set("a", 1)
		a.Set("b", 2)
		assert.Equal(t, 2, a.Get("b"))
	})
}

func testKeys[M any](t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
		m := makeMapi[M](mapT{"a": 1, "b": 2, "c": 3})
		keys := m.Keys()
		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "c")
	})
}

func testValues[M any](t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		m := makeMapi[M](mapT{"a": 1, "b": 2, "c": 3})
		values := m.Values()
		assert.Len(t, values, 3)
		assert.Contains(t, values, 3)
	})
}

func testEqual[M any](t *testing.T) {
	tests := []struct {
		name string
		m    mapTI
		m2   mapT
		want bool
	}{
		{"equal", makeMapi[M](mapT{"a": 1}), mapT{"a": 1}, true},
		{"empty", makeMapi[M](), mapT{}, true},
		{"dif len", makeMapi[M](mapT{"a": 1}), mapT{}, false},
		{"dif len 1", makeMapi[M](mapT{"a": 1}), mapT{"a": 1, "b": 2}, false},
		{"dif value", makeMapi[M](mapT{"a": 1}), mapT{"a": 2}, false},
		{"dif key", makeMapi[M](mapT{"a": 1}), mapT{"b": 1}, false},
	}
	for _, tt := range tests {
		t.Run("Equal "+tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Equal(tt.m2), "Equal(%v)", tt.m2)
		})
	}
}

func testBinaryMarshal[M any](t *testing.T) {
	t.Run("BinaryMarshal", func(t *testing.T) {
		// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding
		m := makeMapi[M](mapT{"a": 1, "b": 2, "c": 3})
		var m2 M
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf) // Will write
		dec := gob.NewDecoder(&buf) // Will read

		err := enc.Encode(m)
		assert.NoError(t, err)
		err = dec.Decode(&m2)
		assert.NoError(t, err)
		var i interface{}
		i = &m2
		m3 := i.(MapI[string, int])
		assert.Equal(t, 1, m3.Get("a"))
		assert.Equal(t, 3, m3.Get("c"))
	})
}

func testMarshalJSON[M any](t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		m := makeMapi[M](mapT{"a": 1, "b": 2, "c": 3})
		s, err := json.Marshal(m)
		assert.NoError(t, err)
		// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
		assert.Equal(t, `{"a":1,"b":2,"c":3}`, string(s))
	})
}

func testUnmarshalJSON[M any](t *testing.T) {
	b := []byte(`{"a":1,"b":2,"c":3}`)
	var m M

	json.Unmarshal(b, &m)
	var i interface{}
	i = &m
	m2 := i.(MapI[string, int])

	assert.Equal(t, 3, m2.Get("c"))
}

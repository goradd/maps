package maps

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
func TestEqual(t *testing.T) {
	var mNil1 mapT
	var mNil2 mapT

	tests := []struct {
		name string
		m1   mapTI
		m2   mapTI
		want bool
	}{
		{"nilAll", mNil1, mNil2, true},
		{"nilFirst", mNil1, mapT{"a": 1}, false},
		{"nilSecond", mapT{"a": 1}, mNil1, false},
		{"sameSizeTrue", mapT{"a": 1}, mapT{"a": 1}, true},
		{"sameSizeFalseKey", mapT{"a": 1}, mapT{"b": 1}, false},
		{"sameSizeFalseValue", mapT{"a": 1}, mapT{"a": 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.m1, tt.m2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

type mySlice []int

func (s mySlice) Equal(b any) bool {
	if s2, ok := b.(mySlice); ok {
		if len(s) == len(s2) {
			for i, v := range s2 {
				if s[i] != v {
					return false
				}
			}
			return true
		}
	}
	return false
}

func TestEqualValues(t *testing.T) {
	a := 1
	b := 1
	assert.True(t, equalValues(a, b))
	b = 2
	assert.False(t, equalValues(a, b))

	c := mySlice{1, 2}
	d := mySlice{1, 2}
	assert.True(t, equalValues(c, d))

	e := []float32{1, 2}
	f := []float32{1, 2}
	assert.Panics(t, func() { equalValues(e, f) })
}

func TestT(t *testing.T) {
	a := maker[Map[int, int], int, int]()
	_ = a
}

package sdc

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGet(t *testing.T) {
	array := NewArray()
	values := make([]int, 1e5)
	for i := range values {
		v := rand.Int()

		values[i] = v
		array.Add(v)
	}

	for i, v := range values {
		if !assert.Equal(t, v, array.Get(i)) {
			break
		}
	}

	assert.Equal(t, len(values), array.Len())
}

func TestEncodeDecode(t *testing.T) {
	array := &Array{}
	values := make([]int, 1e5)
	for i := range values {
		v := rand.Int()

		values[i] = v
		array.Add(v)
	}

	data, _ := array.GobEncode()
	narray := &Array{}
	narray.GobDecode(data)

	for i, v := range values {
		if !assert.Equal(t, v, narray.Get(i)) {
			break
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	array := NewArray()
	values := make([]int, b.N)
	for i := range values {
		values[i] = rand.Int()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.Add(values[i])
	}
}

func BenchmarkGet(b *testing.B) {
	array := NewArray()
	for i := 0; i < 1e5; i++ {
		array.Add(rand.Int())
	}

	idx := make([]int, b.N)
	for i := range idx {
		idx[i] = rand.Intn(array.Len())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.Get(idx[i])
	}
}

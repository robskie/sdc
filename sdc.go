// Package sdc implements simple dense coding which compresses an array of
// integers and allows random access in constant time.
//
// For more details see Simple Random Access Compression:
// http://www.cs.uku.fi/~fredriks/pub/papers/fi09.pdf
package sdc

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/robskie/bit"
	"github.com/robskie/ranksel"
)

// These are the rank and select
// sampling block sizes used by
// ranksel.BitVector.
const (
	sr = 1024
	ss = 1792
)

// Array represents an array of integers.
type Array struct {
	bits     *bit.Array
	selector *ranksel.BitVector

	length      int
	initialized bool
}

func (a *Array) init() {
	a.bits = bit.NewArray(0)

	selOpts := &ranksel.Options{Sr: sr, Ss: ss}
	a.selector = ranksel.NewBitVector(selOpts)
	a.selector.Add(1, 1)

	a.initialized = true
}

// NewArray returns an empty array.
func NewArray() *Array {
	array := &Array{}
	array.init()
	return array
}

func encode(value int) (uint64, int) {
	v := uint64(value)
	length := bit.MSBIndex(v + 2)
	code := v + 2 - (1 << uint(length))

	return code, length
}

// Add adds an integer to the array.
func (a *Array) Add(v int) {
	if !a.initialized {
		a.init()
	}

	a.length++
	code, length := encode(v)

	a.bits.Add(code, length)
	a.selector.Add(1<<uint(length-1), length)
}

func decode(value uint64, length int) int {
	return int(value - 2 + (1 << uint(length)))
}

// Get returns the value at index i.
func (a *Array) Get(i int) int {
	start := a.selector.Select1(i + 1)
	bits := a.selector.Get(start, min(64, a.selector.Len()-start))

	length := bit.Select(bits, 2)
	code := a.bits.Get(start, length)

	return decode(code, length)
}

// Len returns the number of values stored.
func (a *Array) Len() int {
	return a.length
}

// Size returns the array size in bytes.
func (a *Array) Size() int {
	size := a.bits.Size()
	size += a.selector.Size()

	return size
}

// GobEncode encodes this array into gob streams.
func (a *Array) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	err := checkErr(
		enc.Encode(a.bits),
		enc.Encode(a.selector),
		enc.Encode(a.length),
		enc.Encode(a.initialized),
	)

	if err != nil {
		err = fmt.Errorf("sdc: encode failed (%v)", err)
	}

	return buf.Bytes(), err
}

// GobDecode populates this array from gob streams.
func (a *Array) GobDecode(data []byte) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)

	a.bits = bit.NewArray(0)
	a.selector = ranksel.NewBitVector(nil)
	err := checkErr(
		dec.Decode(a.bits),
		dec.Decode(a.selector),
		dec.Decode(&a.length),
		dec.Decode(&a.initialized),
	)

	if err != nil {
		err = fmt.Errorf("sdc: decode failed (%v)", err)
	}

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func checkErr(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}

	return nil
}

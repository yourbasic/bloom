// Package bloom provides a Bloom filter implementation.
//
// Bloom filters
//
// A Bloom filter is a space-efficient probabilistic data structure
// used to test set membership. A member query returns either
// ”likely in set” or ”definitely not in set”. Elements can be added,
// but not removed. With more elements in the set, the probability of
// false positives increases.
//
// Implementation
//
// A  full filter with a false-positives rate of 1/p uses roughly
// 0.26ln(p) bytes per element and performs ⌈1.4ln(p)⌉ bit array lookups
// per query:
//
//	    p     bytes   lookups
//	-------------------------
//	    4      0.4       2
//	    8      0.5       3
//	   16      0.7       4
//	   32      0.9       5
//	   64      1.1       6
//	  128      1.3       7
//	  256      1.5       8
//	  512      1.6       9
//	 1024      1.8      10
//
// This implementation is not intended for cryptographic use.
// Each membership query makes a single call to a 128-bit MurmurHash3 function.
// This saves on hashing without increasing the false-positives
// probability as shown by Kirsch and Mitzenmacher.
//
package bloom

import (
	"math"
)

const (
	shift = 6
	mask  = 0x3f
)

// Filter represents a Bloom filter.
type Filter struct {
	data    []uint64 // Bit array, the length is a power of 2.
	lookups int      // Lookups per query
	count   int64    // Estimated number of unique elements
}

var murmur = new(digest)

// New creates an empty Bloom filter with room for n elements
// at a false-positives rate less than 1/p.
func New(n int, p int) *Filter {
	f := &Filter{}
	minWords := int(0.0325 * math.Log(float64(p)) * float64(n))
	words := 1
	for words < minWords {
		words *= 2
	}
	f.data = make([]uint64, words)
	f.lookups = int(1.4*math.Log(float64(p)) + 1)
	return f
}

// AddByte adds b to the filter and tells if b was already a likely member.
func (f *Filter) AddByte(b []byte) bool {
	h1, h2 := murmur.hash(b)
	trunc := uint64(len(f.data))<<shift - 1
	member := true
	for i := f.lookups; i > 0; i, h1 = i-1, h1+h2 {
		n := h1 & trunc
		k, b := n>>shift, uint64(1<<uint(n&mask))
		if f.data[k]&b == 0 {
			member = false
			f.data[k] |= b
		}
	}
	if !member {
		f.count++
	}
	return member
}

// Add adds s to the filter and tells if s was already a likely member.
func (f *Filter) Add(s string) bool {
	return f.AddByte([]byte(s))
}

// LikelyByte tells if b is a likely member of this filter.
func (f *Filter) LikelyByte(b []byte) bool {
	h1, h2 := murmur.hash(b)
	trunc := uint64(len(f.data))<<shift - 1
	for i := f.lookups; i > 0; i, h1 = i-1, h1+h2 {
		n := h1 & trunc
		k, b := n>>shift, uint64(1<<uint(n&mask))
		if f.data[k]&b == 0 {
			return false
		}
	}
	return true
}

// Likely tells if s is a likely member of this filter.
func (f *Filter) Likely(s string) bool {
	return f.LikelyByte([]byte(s))
}

// Count returns an estimate of the number of unique elements added to this filter.
func (f *Filter) Count() int64 {
	return f.count
}

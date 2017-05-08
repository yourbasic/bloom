// Package bloom provides a Bloom filter implementation.
//
// Bloom filters
//
// A Bloom filter is a space-efficient probabilistic data structure
// used to test set membership. A member test returns either
// ”likely member” or ”definitely not a member”. Only false positives
// can occur: an element that has been added to the filter
// will be identified as ”likely member”.
//
// Elements can be added, but not removed. With more elements in the filter,
// the probability of false positives increases.
//
// Implementation
//
// A  full filter with a false-positives rate of 1/p uses roughly
// 0.26ln(p) bytes per element and performs ⌈1.4ln(p)⌉ bit array lookups
// per test:
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
// Each membership test makes a single call to a 128-bit MurmurHash3 function.
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
	for i := f.lookups; i > 0; i-- {
		h1 += h2
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
	b := make([]byte, len(s))
	copy(b, s)
	return f.AddByte(b)
}

// TestByte tells if b is a likely member of this filter.
func (f *Filter) TestByte(b []byte) bool {
	h1, h2 := murmur.hash(b)
	trunc := uint64(len(f.data))<<shift - 1
	for i := f.lookups; i > 0; i-- {
		h1 += h2
		n := h1 & trunc
		k, b := n>>shift, uint64(1<<uint(n&mask))
		if f.data[k]&b == 0 {
			return false
		}
	}
	return true
}

// Test tells if s is a likely member of this filter.
func (f *Filter) Test(s string) bool {
	b := make([]byte, len(s))
	copy(b, s)
	return f.TestByte(b)
}

// Count returns an estimate of the number of elements in this filter.
func (f *Filter) Count() int64 {
	return f.count
}

// Union returns a new Bloom filter that consists of all elements
// that belong to either f1 or f2. The two filters must be of
// the same size n and have the same false-positives rate p.
//
// The resulting filter is the same as the filter created
// from scratch using the union of the two sets.
func (f1 *Filter) Union(f2 *Filter) *Filter {
	if len(f1.data) != len(f2.data) || f1.lookups != f2.lookups {
		panic("operation requires filters of the same type")
	}
	len := len(f1.data)
	res := &Filter{
		data:    make([]uint64, len),
		lookups: f1.lookups,
	}
	bitCount := 0
	for i := 0; i < len; i++ {
		w := f1.data[i] | f2.data[i]
		res.data[i] = w
		bitCount += count(w)
	}
	// Estimate the number of elements from the bitCount.
	m := 64 * float64(len)
	n := m / float64(f1.lookups) * math.Log(m/(m-float64(bitCount)))
	res.count = int64(n)
	return res
}

// count returns the number of nonzero bits in w.
func count(w uint64) int {
	// Adapted from github.com/yourbasic/bit/funcs.go.
	const maxw = 1<<64 - 1
	const bpw = 64
	w -= (w >> 1) & (maxw / 3)
	w = w&(maxw/15*3) + (w>>2)&(maxw/15*3)
	w += w >> 4
	w &= maxw / 255 * 15
	w *= maxw / 255
	w >>= (bpw/8 - 1) * 8
	return int(w)
}

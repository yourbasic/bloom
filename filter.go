// Package bloom provides a Bloom filter implementation.
//
// Bloom filters
//
// A Bloom filter is a fast and space-efficient probabilistic data structure
// used to test set membership.
//
// A membership test returns either ”likely member” or ”definitely not
// a member”. Only false positives can occur: an element that has been added
// to the filter will always be identified as ”likely member”.
//
// The probabilities of different outcomes of a membership test at
// a false-positives rate of 1/100 are:
//
//	Test(s)                 true     false
//	--------------------------------------
//	s has been added        1        0
//	s has not been added    0.01     0.99
//
// Elements can be added, but not removed. With more elements in the filter,
// the probability of false positives increases.
//
// Performance
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
// Each membership test makes a single call to a 128-bit hash function.
// This improves speed without increasing the false-positives rate
// as shown by Kirsch and Mitzenmacher.
//
// Limitations
//
// This implementation is not intended for cryptographic use.
//
// The internal data representation is different for big-endian
// and little-endian machines.
//
// Typical use case
//
// The Basics example contains a typcial use case:
// a blacklist of shady websites.
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
	count   int64    // Estimate number of elements
}

// New creates an empty Bloom filter with room for n elements
// at a false-positives rate less than 1/p.
func New(n int, p int) *Filter {
	minWords := int(0.0325 * math.Log(float64(p)) * float64(n))
	words := 1
	for words < minWords {
		words *= 2
	}
	return &Filter{
		data:    make([]uint64, words),
		lookups: int(1.4*math.Log(float64(p)) + 1),
	}
}

// AddByte adds b to the filter and tells if b was already a likely member.
func (f *Filter) AddByte(b []byte) bool {
	return f.add(hash(b))
}

// Add adds s to the filter and tells if s was already a likely member.
func (f *Filter) Add(s string) bool {
	return f.add(hashString(s))
}

func (f *Filter) add(h1, h2 uint64) bool {
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

// TestByte tells if b is a likely member of the filter.
// If true, b is probably a member; if false, b is definitely not a member.
func (f *Filter) TestByte(b []byte) bool {
	return f.test(hash(b))
}

// Test tells if s is a likely member of the filter.
// If true, s is probably a member; if false, s is definitely not a member.
func (f *Filter) Test(s string) bool {
	return f.test(hashString(s))
}

func (f *Filter) test(h1, h2 uint64) bool {
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

// Count returns an estimate of the number of elements in the filter.
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

package bloom_test

import (
	"fmt"
	"github.com/yourbasic/bloom"
	"math/rand"
	"strconv"
)

// Create and use a Bloom filter.
func Example_basics() {
	// Create a Bloom filter with room for 10000 elements
	// at a false-positives rate less than 0.5 percent.
	blacklist := bloom.New(10000, 200)

	// Add an element to the filter.
	url := "https://rascal.com"
	blacklist.Add(url)

	// Test for membership.
	if blacklist.Test(url) {
		fmt.Println(url, "seems to be shady.")
	} else {
		fmt.Println(url, "has not yet been added to our blacklist.")
	}
	// Output: https://rascal.com seems to be shady.
}

// Count the number of false positives.
func Example_falsePositives() {
	// Create a Bloom filter with room for n elements
	// at a false-positives rate less than 1/p.
	n := 1000
	p := 100
	filter := bloom.New(n, p)

	// Add n random strings.
	for i := 0; i < n; i++ {
		filter.Add(strconv.FormatUint(rand.Uint64(), 10))
	}

	// Do n random lookups and count the (mostly accidental) hits.
	// It shouldn't be much more than n/p, and hopefully less.
	count := 0
	for i := 0; i < n; i++ {
		if filter.Test(strconv.FormatUint(rand.Uint64(), 10)) {
			count++
		}
	}
	fmt.Println(count, "mistakes were made.")
	// Output: 1 mistakes were made.
}

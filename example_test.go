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
	// Create a Bloom filter with room for 10000 elements
	// at a false-positives rate less than 1/100.
	n := 10000
	p := 100
	filter := bloom.New(n, p)

	// Add n random strings.
	for i := 0; i < n; i++ {
		filter.Add(strconv.Itoa(rand.Int()))
	}

	// Do n random lookups and count the (mostly accidental) hits.
	// It shouldn't be much more than n/p, and hopefully less.
	count := 0
	for i := 0; i < n; i++ {
		if filter.Test(strconv.Itoa(rand.Int())) {
			count++
		}
	}
	fmt.Println(count, "mistakes were made.")
	// Output: 26 mistakes were made.
}

// Compute the intersection and union of two filters.
func ExampleFilter_And() {
	// Create two Bloom filter with room for 1000 elements
	// at a false-positives rate less than 1/100.
	n := 1000
	p := 100
	f1, f2 := bloom.New(n, p), bloom.New(n, p)

	// Add "0", "1", …, "499" to f1
	for i := 0; i < n/2; i++ {
		f1.Add(strconv.Itoa(i))
	}

	// Add "250", "251", …, "749" to f2
	for i := n / 4; i < 3*n/4; i++ {
		f2.Add(strconv.Itoa(i))
	}

	// Compute the approximate size of f1 ∩ f2 and f1 ∪ f2.
	fmt.Println("f1 ∩ f2:", f1.And(f2).Count())
	fmt.Println("f1 ∪ f2:", f1.Or(f2).Count())
	// Output:
	// f1 ∩ f2: 276
	// f1 ∪ f2: 758

}

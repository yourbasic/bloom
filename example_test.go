package bloom_test

import (
	"fmt"
	"github.com/yourbasic/bloom"
	"strconv"
)

// Build a blacklist of shady websites.
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

// Compute the union of two filters.
func ExampleFilter_Union() {
	// Create two Bloom filters, each with room for 1000 elements
	// at a false-positives rate less than 1/100.
	n, p := 1000, 100
	f1 := bloom.New(n, p)
	f2 := bloom.New(n, p)

	// Add "0", "2", …, "498" to f1
	for i := 0; i < n/2; i += 2 {
		f1.Add(strconv.Itoa(i))
	}

	// Add "1", "3", …, "499" to f2
	for i := 1; i < n/2; i += 2 {
		f2.Add(strconv.Itoa(i))
	}

	// Compute the approximate size of f1 ∪ f2.
	fmt.Println("f1 ∪ f2:", f1.Union(f2).Count())
	// Output: f1 ∪ f2: 505
}

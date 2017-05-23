package bloom_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/yourbasic/bloom"
	"log"
	"math/rand"
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

// Estimate the number of false positives.
func Example_falsePositives() {
	// Create a Bloom filter with room for n elements
	// at a false-positives rate less than 1/p.
	n, p := 10000, 100
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

// Send a filter over a network using the encoding/gob package.
func ExampleFilter_MarshalBinary_network() {
	// Create a mock network and a new Filter.
	var network bytes.Buffer
	f1 := bloom.New(1000, 100)
	f1.Add("Hello, filter!")

	// Create an encoder and send the filter to the network.
	enc := gob.NewEncoder(&network)
	if err := enc.Encode(f1); err != nil {
		log.Fatal("encode error:", err)
	}

	// Create a decoder and receive the filter from the network.
	dec := gob.NewDecoder(&network)
	var f2 bloom.Filter
	if err := dec.Decode(&f2); err != nil {
		log.Fatal("decode error:", err)
	}

	// Check that we got the same filter back.
	if f2.Test("Hello, filter!") {
		fmt.Println("Filter arrived safely.")
	}
	// Output: Filter arrived safely.
}

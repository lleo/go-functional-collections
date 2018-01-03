package set_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
	"github.com/lleo/go-functional-collections/set"
)

func ExampleSet_Range() {
	var s = set.New().
		Set(hash.StringKey("a")).
		Set(hash.StringKey("b")).
		Set(hash.StringKey("c"))

	s.Range(func(k hash.Key) bool {
		// Does not provide hash.Key entries in string order.
		fmt.Println(k)
		return true
	})

	// Output:
	// c
	// b
	// a
}

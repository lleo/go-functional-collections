package set_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
	"github.com/lleo/go-functional-collections/set"
)

func ExampleSet_Iter() {
	var s = set.New().
		Set(hash.StringKey("a")).
		Set(hash.StringKey("b")).
		Set(hash.StringKey("c"))

	var it = s.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		fmt.Println(k)
	}

	// Output:
	// c
	// b
	// a
}

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

package set_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/set"
)

func ExampleSet_Iter() {
	var s = set.New().
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

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
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	s.Range(func(k key.Hash) bool {
		// Does not provide key.Hash entries in string order.
		fmt.Println(k)
		return true
	})

	// Output:
	// c
	// b
	// a
}

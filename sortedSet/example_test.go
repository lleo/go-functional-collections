package sortedSet_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/sortedSet"
)

func ExampleSet_Iter() {
	var s = sortedSet.New().
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	var it = s.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		fmt.Println(k)
	}

	// Output:
	// a
	// b
	// c
}

func ExampleRange() {
	var s = sortedSet.New().
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c"))

	s.Range(func(k key.Sort) bool {
		// Provides key.Sort entries in string order.
		fmt.Println(k)
		return true
	})

	// Output:
	// a
	// b
	// c
}

func ExampleRangeLimit() {
	var s = sortedSet.New().
		Set(key.Str("a")).
		Set(key.Str("b")).
		Set(key.Str("c")).
		Set(key.Str("d")).
		Set(key.Str("e"))

	s.RangeLimit(key.Str("b"), key.Str("d"),
		func(k key.Sort) bool {
			// Provides key.Sort entries in string order.
			fmt.Println(k)
			return true
		})

	// Output:
	// b
	// c
	// d
}

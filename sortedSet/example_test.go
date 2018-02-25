package sortedSet_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/sorted"
	"github.com/lleo/go-functional-collections/sortedSet"
)

func ExampleSet_Iter() {
	var s = sortedSet.New().
		Set(sorted.StringKey("a")).
		Set(sorted.StringKey("b")).
		Set(sorted.StringKey("c"))

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
		Set(sorted.StringKey("a")).
		Set(sorted.StringKey("b")).
		Set(sorted.StringKey("c"))

	s.Range(func(k sorted.Key) bool {
		// Provides sorted.Key entries in string order.
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
		Set(sorted.StringKey("a")).
		Set(sorted.StringKey("b")).
		Set(sorted.StringKey("c")).
		Set(sorted.StringKey("d")).
		Set(sorted.StringKey("e"))

	s.RangeLimit(sorted.StringKey("b"), sorted.StringKey("d"),
		func(k sorted.Key) bool {
			// Provides sorted.Key entries in string order.
			fmt.Println(k)
			return true
		})

	// Output:
	// b
	// c
	// d
}

package sorted_map_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/sorted"
	"github.com/lleo/go-functional-collections/sorted_map"
)

func ExampleRange() {
	var s = sorted_map.New().
		Put(sorted.StringKey("a"), 1).
		Put(sorted.StringKey("b"), 2).
		Put(sorted.StringKey("c"), 3)

	s.Range(func(k sorted.Key, v interface{}) bool {
		// Provides sorted.Key entries in string order.
		fmt.Println(k, v)
		return true
	})

	// Output:
	// a 1
	// b 2
	// c 3
}

func ExampleRangeLimit() {
	var s = sorted_map.New().
		Put(sorted.StringKey("a"), 1).
		Put(sorted.StringKey("b"), 2).
		Put(sorted.StringKey("c"), 3).
		Put(sorted.StringKey("d"), 4).
		Put(sorted.StringKey("e"), 5)

	s.RangeLimit(sorted.StringKey("b"), sorted.StringKey("d"),
		func(k sorted.Key, v interface{}) bool {
			// Provides sorted.Key entries in string order.
			fmt.Println(k, v)
			return true
		})

	// Output:
	// b 2
	// c 3
	// d 4
}

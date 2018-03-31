package sortedMap_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/sortedMap"
)

func ExampleRange() {
	var s = sortedMap.New().
		Put(key.Str("a"), 1).
		Put(key.Str("b"), 2).
		Put(key.Str("c"), 3)

	s.Range(func(k key.Sort, v interface{}) bool {
		// Provides key.Sort entries in string order.
		fmt.Println(k, v)
		return true
	})

	// Output:
	// a 1
	// b 2
	// c 3
}

func ExampleRangeLimit() {
	var s = sortedMap.New().
		Put(key.Str("a"), 1).
		Put(key.Str("b"), 2).
		Put(key.Str("c"), 3).
		Put(key.Str("d"), 4).
		Put(key.Str("e"), 5)

	s.RangeLimit(key.Str("b"), key.Str("d"),
		func(k key.Sort, v interface{}) bool {
			// Provides key.Sort entries in string order.
			fmt.Println(k, v)
			return true
		})

	// Output:
	// b 2
	// c 3
	// d 4
}

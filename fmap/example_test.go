package fmap_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/hash"
)

func ExampleMap_Range() {
	var m = fmap.New().
		Put(hash.StringKey("a"), 1).
		Put(hash.StringKey("b"), 2).
		Put(hash.StringKey("c"), 3)

	m.Range(func(k hash.Key, v interface{}) bool {
		// Does not provide hash.Key entries in string order.
		fmt.Println(k, v)
		return true
	})

	// Output:
	// c 3
	// b 2
	// a 1
}

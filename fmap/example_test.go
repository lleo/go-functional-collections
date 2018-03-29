package fmap_test

import (
	"fmt"

	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/key"
)

func ExampleMap_Range() {
	var m = fmap.New().
		Put(key.Str("a"), 1).
		Put(key.Str("b"), 2).
		Put(key.Str("c"), 3)

	m.Range(func(k key.Hash, v interface{}) bool {
		// Does not provide key.Hash entries in string order.
		fmt.Println(k, v)
		return true
	})

	// Output:
	// c 3
	// b 2
	// a 1
}

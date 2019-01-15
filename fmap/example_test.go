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

	m.Range(func(kv KeyVal) bool {
		fmt.Println(kv.Key, kv.Val)
		return true
	})

	// Output:
	// c 3
	// b 2
	// a 1
}

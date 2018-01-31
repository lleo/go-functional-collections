package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

// KeyVal is a simple struct used to transfer lists ([]KeyVal) from one
// function to another.
type KeyVal struct {
	Key hash.Key
	Val interface{}
}

func (kv KeyVal) String() string {
	return fmt.Sprintf("{%q, %v}", kv.Key, kv.Val)
}

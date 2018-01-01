package set

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

// keyVal is a simple struct used to transfer lists ([]keyVal) from one
// function to another.
type keyVal struct {
	Key hash.Key
	Val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("{%q, %v}", kv.Key, kv.Val)
}

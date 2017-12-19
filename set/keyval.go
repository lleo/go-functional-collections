package set

import (
	"fmt"
)

// keyVal is a simple struct used to transfer lists ([]keyVal) from one
// function to another.
type keyVal struct {
	Key SetKey
	Val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("{%q, %v}", kv.Key, kv.Val)
}

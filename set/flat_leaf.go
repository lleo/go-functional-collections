package set

import (
	"fmt"

	"github.com/lleo/go-functional-collections/set/hash"
)

type flatLeaf struct {
	key SetKey
}

func newFlatLeaf(key SetKey) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	return fl
}

func (l *flatLeaf) hash() hash.HashVal {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s}", l.key)
}

func (l *flatLeaf) get(key SetKey) bool {
	if l.key.Equals(key) {
		return true
	}
	return false
}

// put() maintains the functional behavior that any modification returns a new
// leaf and the original remains unaltered. It returns the new leafI and a bool
// indicating if the key was added ontop of the current leaf key or if
// the val mearly replaced the current key's val (either way a new leafI is
// allocated and returned).
func (l *flatLeaf) put(key SetKey) (leafI, bool) {
	var nl leafI

	if l.key.Equals(key) {
		// maintain functional behavior of flatLeaf
		//nl = newFlatLeaf(key)
		//return nl, false //replaced
		return l, false
	}

	nl = newCollisionLeaf([]SetKey{l.key, key})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key SetKey) (leafI, bool) {
	if l.key.Equals(key) {
		return nil, true //found
	}
	return l, false //not found
}

func (l *flatLeaf) keys() []SetKey {
	return []SetKey{l.key}
}

func (l *flatLeaf) visit(fn visitFn, depth uint) (error, bool) {
	return nil, fn(l, depth)
}

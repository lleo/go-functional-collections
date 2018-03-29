package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
)

type flatLeaf struct {
	key key.Hash
	val interface{}
}

func newFlatLeaf(key key.Hash, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	fl.val = val
	return fl
}

func (l *flatLeaf) copy() leafI {
	var nl = new(flatLeaf)
	nl.key = l.key
	nl.val = l.val
	return nl
}

func (l *flatLeaf) hash() hash.Val {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l *flatLeaf) get(key key.Hash) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

// putResolve maintains the functional behavior that any modification returns a
// new leaf and the original remains unaltered.
//
// If the given current key Equals() the given key, then the resolve function
// is used to generate a new value that is used to generate a new flatLeaf.
//
// If the current key does not equal the given key, then a new collisionLeaf is
// generated which adds the current flatLeaf's key,val pair to the given key,val
// pair in the returned collisionLeaf.
func (l *flatLeaf) putResolve(
	key key.Hash,
	val interface{},
	resolve ResolveConflictFunc,
) (leafI, bool) {
	var nl leafI

	if l.key.Equals(key) {
		// maintain functional behavior of flatLeaf
		var newVal = resolve(l.key, l.val, val)
		nl = newFlatLeaf(l.key, newVal)
		return nl, false // replaced
	}

	nl = newCollisionLeaf([]KeyVal{{l.key, l.val}, {key, val}})
	return nl, true // key,val was added
}

// put maintains the functional behavior that any modification returns a new
// leaf and the original remains unaltered.
//
// If the current key Equals() the given key, then a new flatLeaf is generate
// to replace the current flatLeaf's value with the given value.
//
// If the current key does not equal the given key, then a new collisionLeaf is
// generated which adds the current flatLeaf's key,val pair to the given key,val
// pair in the returned collisionLeaf.
func (l *flatLeaf) put(key key.Hash, val interface{}) (leafI, bool) {
	var nl leafI

	if l.key.Equals(key) {
		// maintain functional behavior of flatLeaf
		nl = newFlatLeaf(l.key, val)
		return nl, false // replaced
	}

	nl = newCollisionLeaf([]KeyVal{{l.key, l.val}, {key, val}})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key key.Hash) (leafI, interface{}, bool) {
	if l.key.Equals(key) {
		return nil, l.val, true // found
	}
	return l, nil, false // not found
}

func (l *flatLeaf) keyVals() []KeyVal {
	return []KeyVal{{l.key, l.val}}
}

func (l *flatLeaf) walkPreOrder(fn visitFunc, depth uint) bool {
	return fn(l, depth)
}

// equiv comparse this *flatLeaf against another node by value.
func (l *flatLeaf) equiv(other nodeI) bool {
	var ol, ok = other.(*flatLeaf)
	if !ok {
		return false
	}
	if !l.key.Equals(ol.key) {
		return false
	}
	if l.val != ol.val {
		return false
	}
	return true
}

func (l *flatLeaf) count() int {
	return 1
}

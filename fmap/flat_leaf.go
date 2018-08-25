package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
)

type flatLeaf KeyVal

func newFlatLeaf(key key.Hash, val interface{}) *flatLeaf {
	return &flatLeaf{Key: key, Val: val}
}

func (l *flatLeaf) copy() leafI {
	return &flatLeaf{Key: l.Key, Val: l.Val}
}

func (l *flatLeaf) hash() hash.Val {
	return l.Key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.Key, l.Val)
}

func (l *flatLeaf) get(key key.Hash) (interface{}, bool) {
	if l.Key.Equals(key) {
		return l.Val, true
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

	if l.Key.Equals(key) {
		// maintain functional behavior of flatLeaf
		var newVal = resolve(l.Key, l.Val, val)
		nl = newFlatLeaf(l.Key, newVal)
		return nl, false // replaced
	}

	nl = newCollisionLeaf([]KeyVal{{l.Key, l.Val}, {key, val}})
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

	if l.Key.Equals(key) {
		// maintain functional behavior of flatLeaf
		nl = newFlatLeaf(l.Key, val)
		return nl, false // replaced
	}

	nl = newCollisionLeaf([]KeyVal{{l.Key, l.Val}, {key, val}})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key key.Hash) (leafI, interface{}, bool) {
	if l.Key.Equals(key) {
		return nil, l.Val, true // found
	}
	return l, nil, false // not found
}

func (l *flatLeaf) keyVals() []KeyVal {
	return []KeyVal{{Key: l.Key, Val: l.Val}}
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
	if !l.Key.Equals(ol.Key) {
		return false
	}
	if l.Val != ol.Val {
		return false
	}
	return true
}

func (l *flatLeaf) count() int {
	return 1
}

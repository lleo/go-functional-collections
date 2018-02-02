package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

type flatLeaf struct {
	key hash.Key
	val interface{}
}

func newFlatLeaf(key hash.Key, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	fl.val = val
	return fl
}

func (l *flatLeaf) hash() hash.Val {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l *flatLeaf) get(key hash.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

// put() maintains the functional behavior that any modification returns a new
// leaf and the original remains unaltered. It returns the new leafI and a bool
// indicating if the key,val was added ontop of the current leaf key,val or if
// the val mearly replaced the current key's val (either way a new leafI is
// allocated and returned).
func (l *flatLeaf) put(key hash.Key, val interface{}) (leafI, bool) {
	var nl leafI

	if l.key.Equals(key) {
		// maintain functional behavior of flatLeaf
		nl = newFlatLeaf(l.key, val)
		return nl, false // replaced
	}

	nl = newCollisionLeaf([]KeyVal{{l.key, l.val}, {key, val}})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key hash.Key) (leafI, interface{}, bool) {
	if l.key.Equals(key) {
		return nil, l.val, true // found
	}
	return l, nil, false // not found
}

func (l *flatLeaf) keyVals() []KeyVal {
	return []KeyVal{{l.key, l.val}}
}

func (l *flatLeaf) walkPreOrder(fn visitFn, depth uint) (bool, error) {
	return fn(l, depth), nil
}

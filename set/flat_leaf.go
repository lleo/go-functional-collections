package set

import (
	"fmt"
	"log"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
)

type flatLeaf struct {
	key key.Hash
}

func newFlatLeaf(key key.Hash) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	return fl
}

func (l *flatLeaf) hash() hash.Val {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s}", l.key)
}

func (l *flatLeaf) get(key key.Hash) bool {
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
func (l *flatLeaf) put(k key.Hash) (leafI, bool) {
	var nl leafI

	if l.key.Equals(k) {
		// maintain functional behavior of flatLeaf
		//nl = newFlatLeaf(k)
		//return nl, false //replaced
		return l, false
	}

	nl = newCollisionLeaf([]key.Hash{l.key, k})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key key.Hash) (leafI, bool) {
	if l.key.Equals(key) {
		return nil, true //found
	}
	return l, false //not found
}

func (l *flatLeaf) keys() []key.Hash {
	return []key.Hash{l.key}
}

func (l *flatLeaf) walkPreOrder(fn visitFunc, depth uint) bool {
	return fn(l, depth)
}

// equiv comparse this *flatLeaf against another node by value.
func (l *flatLeaf) equiv(other nodeI) bool {
	var ol, ok = other.(*flatLeaf)
	if !ok {
		log.Println("other is not a *flatLeaf")
		return false
	}
	if !l.key.Equals(ol.key) {
		log.Println("l.key != ol.key")
		return false
	}
	return true
}

func (l *flatLeaf) count() int {
	return 1
}

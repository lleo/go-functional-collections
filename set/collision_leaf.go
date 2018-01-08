package set

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

// implements nodeI
// implements leafI
type collisionLeaf struct {
	ks []hash.Key
}

func newCollisionLeaf(keys []hash.Key) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.ks = append(leaf.ks, keys...)

	//log.Println("newCollisionLeaf:", leaf)

	return leaf
}

func (l *collisionLeaf) copy() *collisionLeaf {
	var nl = new(collisionLeaf)
	nl.ks = append(nl.ks, l.ks...)
	return nl
}

func (l *collisionLeaf) hash() hash.Val {
	return l.ks[0].Hash()
}

func (l *collisionLeaf) String() string {
	var keystrs = make([]string, len(l.ks))
	for i := 0; i < len(l.ks); i++ {
		keystrs[i] = l.ks[i].String()
	}
	var jkeystr = strings.Join(keystrs, ",")

	return fmt.Sprintf("collisionLeaf{hash:%s, keys:[]hash.Key{%s}}",
		l.ks[0].Hash(), jkeystr)
}

func (l *collisionLeaf) get(key hash.Key) bool {
	for _, keyN := range l.ks {
		if keyN.Equals(key) {
			return true
		}
	}
	return false
}

func (l *collisionLeaf) put(key hash.Key) (leafI, bool) {
	for _, keyN := range l.ks {
		if keyN.Equals(key) {
			//var nl = l.copy()
			//return nl, false //replaced
			return l, false
		}
	}
	var nl = new(collisionLeaf)
	nl.ks = make([]hash.Key, len(l.ks)+1)
	copy(nl.ks, l.ks)
	nl.ks[len(l.ks)] = key
	// v-- this, instead of that --^ make&copy&assign
	//nl.ks = append(nl.ks, append(l.ks, k)...)

	//log.Printf("%s : %d\n", l.hash(), len(l.ks))

	return nl, true // k,v was added
}

func (l *collisionLeaf) del(key hash.Key) (leafI, bool) {
	for i, lkey := range l.ks {
		if lkey.Equals(key) {
			var nl leafI
			if len(l.ks) == 2 {
				// think about the index... it works, really :)
				nl = newFlatLeaf(l.ks[1-i])
			} else {
				var cl = l.copy()
				cl.ks = append(cl.ks[:i], cl.ks[i+1:]...)
				nl = cl // needed access to cl.ks; nl is type leafI
			}
			//log.Printf("l.del(); kv=%s removed; returning %s", kv, nl)
			return nl, true
		}
	}
	//log.Printf("cl.del(%s) removed nothing.", k)
	return l, false
}

func (l *collisionLeaf) keys() []hash.Key {
	var r = make([]hash.Key, 0, len(l.ks))
	r = append(r, l.ks...)
	return r
	//return l.ks
}

func (l *collisionLeaf) visit(fn visitFn, depth uint) (bool, error) {
	return fn(l, depth), nil
}

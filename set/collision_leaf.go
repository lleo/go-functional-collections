package set

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

// implements nodeI
// implements leafI
type collisionLeaf struct {
	keys_ []hash.Key
}

func newCollisionLeaf(keys []hash.Key) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.keys_ = append(leaf.keys_, keys...)

	//log.Println("newCollisionLeaf:", leaf)

	return leaf
}

func (l *collisionLeaf) copy() *collisionLeaf {
	var nl = new(collisionLeaf)
	nl.keys_ = append(nl.keys_, l.keys_...)
	return nl
}

func (l *collisionLeaf) hash() hash.Val {
	return l.keys_[0].Hash()
}

func (l *collisionLeaf) String() string {
	var keystrs = make([]string, len(l.keys_))
	for i := 0; i < len(l.keys_); i++ {
		keystrs[i] = l.keys_[i].String()
	}
	var jkeystr = strings.Join(keystrs, ",")

	return fmt.Sprintf("collisionLeaf{hash:%s, keys:[]hash.Key{%s}}",
		l.keys_[0].Hash(), jkeystr)
}

func (l *collisionLeaf) get(key hash.Key) bool {
	for _, keyN := range l.keys_ {
		if keyN.Equals(key) {
			return true
		}
	}
	return false
}

func (l *collisionLeaf) put(key hash.Key) (leafI, bool) {
	for _, keyN := range l.keys_ {
		if keyN.Equals(key) {
			//var nl = l.copy()
			//return nl, false //replaced
			return l, false
		}
	}
	var nl = new(collisionLeaf)
	nl.keys_ = make([]hash.Key, len(l.keys_)+1)
	copy(nl.keys_, l.keys_)
	nl.keys_[len(l.keys_)] = key
	// v-- this, instead of that --^ make&copy&assign
	//nl.keys_ = append(nl.keys_, append(l.keys_, k)...)

	//log.Printf("%s : %d\n", l.hash(), len(l.keys_))

	return nl, true // k,v was added
}

func (l *collisionLeaf) del(key hash.Key) (leafI, bool) {
	for i, lkey := range l.keys_ {
		if lkey.Equals(key) {
			var nl leafI
			if len(l.keys_) == 2 {
				// think about the index... it works, really :)
				nl = newFlatLeaf(l.keys_[1-i])
			} else {
				var cl = l.copy()
				cl.keys_ = append(cl.keys_[:i], cl.keys_[i+1:]...)
				nl = cl // needed access to cl.keys_; nl is type leafI
			}
			//log.Printf("l.del(); kv=%s removed; returning %s", kv, nl)
			return nl, true
		}
	}
	//log.Printf("cl.del(%s) removed nothing.", k)
	return l, false
}

func (l *collisionLeaf) keys() []hash.Key {
	var r = make([]hash.Key, 0, len(l.keys_))
	r = append(r, l.keys_...)
	return r
	//return l.keys_
}

func (l *collisionLeaf) visit(fn visitFn, depth uint) (error, bool) {
	return nil, fn(l, depth)
}

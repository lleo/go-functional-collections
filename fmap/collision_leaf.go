package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/fmap/hash"
)

// implements nodeI
// implements leafI
type collisionLeaf struct {
	kvs []keyVal
}

func newCollisionLeaf(kvs []keyVal) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.kvs = append(leaf.kvs, kvs...)

	//log.Println("newCollisionLeaf:", leaf)

	return leaf
}

func (l *collisionLeaf) copy() *collisionLeaf {
	var nl = new(collisionLeaf)
	nl.kvs = append(nl.kvs, l.kvs...)
	return nl
}

func (l *collisionLeaf) hash() hash.HashVal {
	return l.kvs[0].Key.Hash()
}

func (l *collisionLeaf) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash:%s, kvs:[]keyVal{%s}}",
		l.kvs[0].Key.Hash(), jkvstr)
}

func (l *collisionLeaf) get(key MapKey) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.Key.Equals(key) {
			return kv.Val, true
		}
	}
	return nil, false
}

func (l *collisionLeaf) put(key MapKey, val interface{}) (leafI, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			var nl = l.copy()
			nl.kvs[i].Val = val
			return nl, false //replaced
		}
	}
	var nl = new(collisionLeaf)
	nl.kvs = make([]keyVal, len(l.kvs)+1)
	copy(nl.kvs, l.kvs)
	nl.kvs[len(l.kvs)] = keyVal{key, val}
	//nl.kvs = append(nl.kvs, append(l.kvs, keyVal{k, v})...)

	//log.Printf("%s : %d\n", l.hash(), len(l.kvs))

	return nl, true // k,v was added
}

func (l *collisionLeaf) del(key MapKey) (leafI, interface{}, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			var nl leafI
			if len(l.kvs) == 2 {
				// think about the index... it works, really :)
				nl = newFlatLeaf(l.kvs[1-i].Key, l.kvs[1-i].Val)
			} else {
				var cl = l.copy()
				cl.kvs = append(cl.kvs[:i], cl.kvs[i+1:]...)
				nl = cl // needed access to cl.kvs; nl is type leafI
			}
			//log.Printf("l.del(); kv=%s removed; returning %s", kv, nl)
			return nl, kv.Val, true
		}
	}
	//log.Printf("cl.del(%s) removed nothing.", k)
	return l, nil, false
}

func (l *collisionLeaf) keyVals() []keyVal {
	var r = make([]keyVal, 0, len(l.kvs))
	r = append(r, l.kvs...)
	return r
	//return l.kvs
}

func (l *collisionLeaf) visit(fn visitFn, depth uint) (error, bool) {
	return nil, fn(l, depth)
}

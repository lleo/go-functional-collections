package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
)

// implements nodeI
// implements leafI
type collisionLeaf struct {
	kvs []KeyVal
}

func newCollisionLeaf(kvs []KeyVal) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.kvs = append(leaf.kvs, kvs...)

	//log.Println("newCollisionLeaf:", leaf)

	return leaf
}

func (l *collisionLeaf) copy() leafI {
	var nl = new(collisionLeaf)
	nl.kvs = append(nl.kvs, l.kvs...)
	return nl
}

func (l *collisionLeaf) hash() hash.Val {
	return l.kvs[0].Key.Hash()
}

func (l *collisionLeaf) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash:%s, kvs:[]KeyVal{%s}}",
		l.kvs[0].Key.Hash(), jkvstr)
}

func (l *collisionLeaf) get(key key.Hash) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.Key.Equals(key) {
			return kv.Val, true
		}
	}
	return nil, false
}

func (l *collisionLeaf) putResolve(
	key key.Hash,
	val interface{},
	resolve ResolveConflictFunc,
) (leafI, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			var nl = l.copy().(*collisionLeaf)
			var newVal = resolve(kv.Key, kv.Val, val)
			nl.kvs[i].Val = newVal
			return nl, false // replaced
		}
	}
	var nl = new(collisionLeaf)
	nl.kvs = make([]KeyVal, len(l.kvs)+1)
	copy(nl.kvs, l.kvs)
	nl.kvs[len(l.kvs)] = KeyVal{key, val}
	return nl, true // k,v was added
}

func (l *collisionLeaf) put(key key.Hash, val interface{}) (leafI, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			var nl = l.copy().(*collisionLeaf)
			nl.kvs[i].Val = val
			return nl, false // replaced
		}
	}
	var nl = new(collisionLeaf)
	nl.kvs = make([]KeyVal, len(l.kvs)+1)
	copy(nl.kvs, l.kvs)
	nl.kvs[len(l.kvs)] = KeyVal{key, val}
	//nl.kvs = append(nl.kvs, append(l.kvs, KeyVal{k, v})...)

	//log.Printf("%s : %d\n", l.hash(), len(l.kvs))

	return nl, true // k,v was added
}

func (l *collisionLeaf) del(key key.Hash) (leafI, interface{}, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			var nl leafI
			if len(l.kvs) == 2 {
				// think about the index... it works, really :)
				nl = newFlatLeaf(l.kvs[1-i].Key, l.kvs[1-i].Val)
			} else {
				var cl = l.copy().(*collisionLeaf)
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

func (l *collisionLeaf) keyVals() []KeyVal {
	var r = make([]KeyVal, 0, len(l.kvs))
	r = append(r, l.kvs...)
	return r
	//return l.kvs
}

func (l *collisionLeaf) walkPreOrder(fn visitFunc, depth uint) bool {
	return fn(l, depth)
}

// equiv comparse this *collisionLeaf against another node by value.
func (l *collisionLeaf) equiv(other nodeI) bool {
	var ol, ok = other.(*collisionLeaf)
	if !ok {
		return false
	}
	if len(l.kvs) != len(ol.kvs) {
		return false
	}
	// This assumes the kvs are in the same order.
	for i, kv := range l.kvs {
		if !kv.Key.Equals(ol.kvs[i].Key) {
			return false
		}
		if kv.Val != ol.kvs[i].Val {
			return false
		}
	}
	return true
}

func (l *collisionLeaf) count() int {
	return len(l.kvs)
}

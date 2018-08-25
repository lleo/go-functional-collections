package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
)

// implements nodeI
// implements leafI
type collisionLeaf []KeyVal

func newCollisionLeaf(kvs []KeyVal) *collisionLeaf {
	var lKvs collisionLeaf = make([]KeyVal, len(kvs))
	copy(lKvs, kvs)
	return &lKvs
}

func (l *collisionLeaf) copy() leafI {
	//return newCollisionLeaf([]KeyVal(*l))
	return newCollisionLeaf(*l)
}

func (l *collisionLeaf) hash() hash.Val {
	return (*l)[0].Key.Hash()
}

func (l *collisionLeaf) String() string {
	var kvstrs = make([]string, len(*l))
	for i := 0; i < len(*l); i++ {
		kvstrs[i] = (*l)[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash:%s, kvs:[]KeyVal{%s}}",
		(*l)[0].Key.Hash(), jkvstr)
}

func (l *collisionLeaf) get(key key.Hash) (interface{}, bool) {
	for _, kv := range *l {
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
	for i, kv := range *l {
		if kv.Key.Equals(key) {
			var nl = l.copy().(*collisionLeaf)
			var newVal = resolve(kv.Key, kv.Val, val)
			(*nl)[i].Val = newVal
			return nl, false // replaced
		}
	}
	var nl collisionLeaf = make([]KeyVal, len(*l)+1)
	nl[len(*l)] = KeyVal{Key: key, Val: val}
	return &nl, true
}

func (l *collisionLeaf) put(key key.Hash, val interface{}) (leafI, bool) {
	for i, kv := range *l {
		if kv.Key.Equals(key) {
			var nl = l.copy().(*collisionLeaf)
			(*nl)[i].Val = val
			return nl, false // replaced
		}
	}
	var nl collisionLeaf = make([]KeyVal, len(*l)+1)
	nl[len(*l)] = KeyVal{Key: key, Val: val}
	return &nl, true
}

func (l *collisionLeaf) del(key key.Hash) (leafI, interface{}, bool) {
	for i, kv := range *l {
		if kv.Key.Equals(key) {
			var nl leafI
			if len(*l) == 2 {
				// think about the index... it works, really :)
				nl = newFlatLeaf((*l)[1-i].Key, (*l)[1-i].Val)
			} else {
				var cl = l.copy().(*collisionLeaf)
				*cl = append((*cl)[:i], (*cl)[i+1:]...)
				nl = cl
			}
			return nl, kv.Val, true
		}
	}
	return l, nil, false
}

func (l *collisionLeaf) keyVals() []KeyVal {
	var r = make([]KeyVal, 0, len(*l))
	r = append(r, *l...)
	return r
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
	if len(*l) != len(*ol) {
		return false
	}

	for _, kv := range *l {
		var keyFound = false
		for _, okv := range *ol {
			if kv.Key.Equals(okv.Key) {
				keyFound = true
				if kv.Val != okv.Val {
					return false
				}
			}
		}
		if !keyFound {
			return false
		}
	}

	//for i, kv := range *l {
	//// This assumes the kvs are in the same order.
	//	if !kv.Key.Equals((*ol)[i].Key) {
	//		return false
	//	}
	//	if kv.Val != (*ol)[i].Val {
	//		return false
	//	}
	//}

	return true
}

func (l *collisionLeaf) count() int {
	return len(*l)
}

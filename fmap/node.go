package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

// visitFn will be passed a value for every slot in the Hamt; this includes
// leafs, tables, and nil.
//
// If the visitFn returns false then the tree walk should stop.
//
type visitFn func(nodeI, uint) bool

type nodeI interface {
	hash() hash.Val
	visit(fn visitFn, depth uint) (error, bool)
	String() string
}

type leafI interface {
	nodeI

	get(key hash.Key) (interface{}, bool)
	put(key hash.Key, val interface{}) (leafI, bool)
	del(key hash.Key) (leafI, interface{}, bool)
	keyVals() []keyVal
}

type tableIterFunc func() nodeI

type tableI interface {
	nodeI

	copy() tableI
	deepCopy() tableI

	numEntries() uint
	entries() []tableEntry

	get(idx uint) nodeI

	insert(idx uint, n nodeI)
	replace(idx uint, n nodeI)
	remove(idx uint)

	iter() tableIterFunc

	treeString(string, uint) string
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}

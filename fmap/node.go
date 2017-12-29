package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/fmap/hash"
)

// visitFn will be passed a value for every slot in the Hamt; this includes
// leafs, tables, and nil.
//
// If the visitFn returns false then the tree walk should stop.
//
type visitFn func(nodeI, uint) bool

type nodeI interface {
	hash() hash.HashVal
	visit(fn visitFn, depth uint) (error, bool)
	String() string
}

type leafI interface {
	nodeI

	get(key MapKey) (interface{}, bool)
	put(key MapKey, val interface{}) (leafI, bool)
	del(key MapKey) (leafI, interface{}, bool)
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

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
	walkInOrder(fn visitFn, depth uint) (bool, error)
	String() string
}

type leafI interface {
	nodeI

	get(key hash.Key) (interface{}, bool)
	put(key hash.Key, val interface{}) (leafI, bool)
	del(key hash.Key) (leafI, interface{}, bool)
	keyVals() []KeyVal
}

type tableIterFunc func() nodeI

type tableI interface {
	nodeI

	copy() tableI
	deepCopy() tableI

	slotsUsed() uint //numEntries() uint
	entries() []tableEntry

	get(idx uint) nodeI

	insertInplace(idx uint, n nodeI)
	replaceInplace(idx uint, n nodeI)
	removeInplace(idx uint)

	insert(idx uint, n nodeI) tableI
	replace(idx uint, n nodeI) tableI
	remove(idx uint) tableI

	needsUpgrade() bool
	needsDowngrade() bool
	upgrade() tableI
	downgrade() tableI

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

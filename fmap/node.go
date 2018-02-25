package fmap

import (
	"fmt"

	"github.com/lleo/go-functional-collections/hash"
)

// visitFunc will be passed a value for every slot in the Hamt; this includes
// leafs, tables, and nil.
//
// If the visitFunc returns false then the tree walk should stop.
//
type visitFunc func(nodeI, uint) bool

type nodeI interface {
	hash() hash.Val
	walkPreOrder(fn visitFunc, depth uint) bool
	equiv(nodeI) bool
	count() int
	String() string
}

// ResolveConflictFunc is the signature of functions used to choose between, or
// create a new value from, two key,value pairs where the keys are equal (this
// is defined by k0.Equal(k1), hence only the Map key is passed in).
type ResolveConflictFunc func(
	key hash.Key,
	origVal, newVal interface{},
) interface{}

// KeepOrigVal is an implementation of ResolveConflictFunc type which returns
// the first (origVal) value.
func KeepOrigVal(key hash.Key, origVal, newVal interface{}) interface{} {
	return origVal
}

// TakeNewVal is an implementation of ResolveConflictFunc type which returns
// the second (newVal) value.
func TakeNewVal(key hash.Key, origVal, newVal interface{}) interface{} {
	return newVal
}

type leafI interface {
	nodeI

	get(key hash.Key) (interface{}, bool)
	putResolve(key hash.Key, val interface{}, resolve ResolveConflictFunc) (leafI, bool)
	put(key hash.Key, val interface{}) (leafI, bool)
	del(key hash.Key) (leafI, interface{}, bool)

	copy() leafI
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

package set

import (
	"fmt"

	"github.com/lleo/go-functional-collections/set/hash"
)

// visitFn will be passed a value for every slot in the Hamt; this includes
// leafs, tables, and nil.
//
// If the visitFn returns false then the tree walk should stop.
//
type visitFn func(nodeI, uint) bool

type nodeI interface {
	String() string
	hash() hash.HashVal
	visit(fn visitFn, depth uint) (error, bool)
}

type leafI interface {
	nodeI

	get(key SetKey) bool
	put(key SetKey) (leafI, bool)
	del(key SetKey) (leafI, bool)
	keys() []SetKey
}

type tableIterFunc func() nodeI

type tableI interface {
	nodeI

	copy() tableI
	deepCopy() tableI

	LongString(indent string, depth uint) string

	numEntries() uint
	entries() []tableEntry

	get(idx uint) nodeI

	insert(idx uint, n nodeI)
	replace(idx uint, n nodeI)
	remove(idx uint)

	iter() tableIterFunc
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}

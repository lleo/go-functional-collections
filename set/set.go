// Package set implements a functional Set data structure. The internal data
// structure of set is a Hashed Array Mapped Trie (see https://en.wikipedia.org/wiki/Hash_array_mapped_trie).
//
// Functional means that each data structure is immutable and persistent.
// The Set is immutable because you never modify a Set in place, but rather
// every modification (like a Add or Remove) creates a new Set with that
// modification. This is not as inefficient as it sounds like it would be. Each
// modification only changes the smallest  branch of the data structure it needs
// to in order to effect the new set. Otherwise, the new data structure
// shares the majority of the previous data structure. That is the persistent
// property.
//
// Each method call that potentially modifies the Set, returns a new Set data
// structure in addition to the other pertinent return values.
//
// The unique values stored in a Set must implement the hash.Key interface.
package set

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

// downgradeThreshold is the constant that sets the threshold for the size of a
// table, such that when a table decreases to the threshold size, the table is
// converted from a fixedTable to a sparseTable.
// downgradeThreshold = 8 for hash.numIndexBits=4 aka hash.IndexLimit=16
// downgradeThreshold = 16 for hash.numIndexBits=5 aka hash.IndexLimit=32
const downgradeThreshold uint = hash.IndexLimit / 2

// upgradeThreshold is the constant that sets the threshold for the size of a
// table, such that when a table increases to the threshold size, the table is
// converted from a sparseTable to a fixedTable.
// upgradeThreshold = 10 for hash.numIndexBits=4 aka hash.IndexLimit=16
// upgradeThreshold = 20 for hash.numIndexBits=5 aka hash.IndexLimit=32
const upgradeThreshold uint = hash.IndexLimit * 5 / 8

// The Key struct mainains an immutable collection of hash.Key entries.
type Set struct {
	root    fixedTable
	numEnts uint
}

// New returns a properly initialized pointer to a Set struct.
func New() *Set {
	return new(Set)
}

// copy creates a shallow copy of the Set data structure and returns a pointer
// to that shallow copy.
func (s *Set) copy() *Set {
	var ns = new(Set)
	*ns = *s
	return s
}

// IsSet searches the Set for a hash.Key value where the given key (k) matches
// a key in the Set (k0) such that k.Equals(k0) returns true. If the given
// key is found IsSet return true, otherwise it returns false.
func (s *Set) IsSet(key hash.Key) bool {
	if s.NumEntries() == 0 {
		return false
	}

	var hv = key.Hash()
	var curTable tableI = &s.root

	var found bool

DepthIter:
	for depth := uint(0); depth <= hash.MaxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		switch n := curNode.(type) {
		case nil:
			found = false
			break DepthIter
		case leafI:
			found = n.get(key)
			break DepthIter
		case tableI:
			curTable = n
		}
	}

	return found
}

// find() traverses the path defined by the given hash.Val till it encounters
// a leafI, then it returns the table path leading to the current table (also
// returned) and the Index in the current table the leaf is at.
//func (m *Set) find(hv hash.Val) (*tableStack, tableI, uint) {
func (s *Set) find(hv hash.Val) (*tableStack, leafI, uint) {
	var curTable tableI = &s.root

	var path = newTableStack()
	var leaf leafI
	var idx uint

DepthIter:
	for depth := uint(0); depth <= hash.MaxDepth; depth++ {
		path.push(curTable)
		idx = hv.Index(depth)
		var curNode = curTable.get(idx)

		switch n := curNode.(type) {
		case nil:
			leaf = nil
			break DepthIter
		case leafI:
			leaf = n
			break DepthIter
		case tableI:
			curTable = n
		}
	}

	return path, leaf, idx
}

// persist() is ONLY called on a fresh copy of the current Hamt.
// Hence, modifying it is allowed.
func (s *Set) persist(oldTable, newTable tableI, path *tableStack) {
	// Removed the case where path.len() == 0 on the first call to nh.perist(),
	// because that case is handled in Put & Del now. It is handled in Put & Del
	// because otherwise we were allocating an extraneous fixedTable for the
	// old s.root.
	_ = assertOn && assert(path.len() != 0,
		"path.len()==0; This case should be handled directly in Put & Del.")

	var depth = uint(path.len()) //guaranteed depth > 0
	var parentDepth = depth - 1

	var parentIdx = oldTable.hash().Index(parentDepth)

	var oldParent = path.pop()

	var newParent tableI
	if path.len() == 0 {
		// This condition and the last if path.len() > 0; shaves off one call
		// to persist and one fixed table allocation (via oldParent.copy()).
		s.root = *oldParent.(*fixedTable)
		newParent = &s.root
	} else {
		newParent = oldParent.copy()
	}

	if newTable == nil {
		newParent.remove(parentIdx)
	} else {
		newParent.replace(parentIdx, newTable)
	}

	if path.len() > 0 {
		s.persist(oldParent, newParent, path)
	}

	return
}

// Set inserts the give hash.Key into the Set and returns a new *Set. If an
// equivalent key exists in the receiver *Set nothing is done to the *Set and
// the original receiver *Set is returned.
//
// Equivalentcy of keys is determined by k.Equals(k0) where k is the given
// hash.Key and k0 is the hash.Key already stored in the *Set.
func (s *Set) Set(key hash.Key) *Set {
	s, _ = s.Add(key)
	return s
}

// Add inserts a new key to the *Set data structure. It returns the
// new *Set data structure and a bool indicating if the given key was added. If
// the *Set already contains and equivalent key the *Set is not modified and
// the original *Set is returned along with a false value. The false value
// indicates that the given key was not added.
//
// Equivalentcy of keys is determined by k.Equals(k0) where k is the given
// hash.Key and k0 is the hash.Key already stored in the *Set.
func (s *Set) Add(key hash.Key) (*Set, bool) {
	var ns = s.copy()

	var hv = key.Hash()

	var path, leaf, idx = s.find(hv)
	var curTable = path.pop()

	var depth = uint(path.len())

	var added bool

	if curTable == &s.root {
		// Special handling of root table.
		// 1) It is never upgraded.
		// 2) It doesn't need to be copied, cuz the new Set has a fixed root.
		if leaf == nil {
			// if curTable.get(idx) is slot is empty
			ns.root.insert(idx, newFlatLeaf(key))
			added = true
		} else {
			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key))
				added = true
			}
			ns.root.replace(idx, node)
		}
	} else {
		var newTable tableI

		if leaf == nil {
			// if curTable.get(idx) slot is empty
			if (curTable.numEntries() + 1) == upgradeThreshold {
				newTable = upgradeToFixedTable(
					curTable.hash(), depth, curTable.entries())
			} else {
				newTable = curTable.copy()
			}

			newTable.insert(idx, newFlatLeaf(key))
			added = true
		} else {
			newTable = curTable.copy()

			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key))
				added = true
			}

			newTable.replace(idx, node)
		}

		ns.persist(curTable, newTable, path)
	}

	if added {
		ns.numEnts++
	}

	return ns, added
}

// Unset removes the any hash.Key that is equivalent to the given hash.Key and
// returns the new *Set. If the hash.Key does not exist in the *Set, then
// nothing will occur and the original *Set will be returned.
func (s *Set) Unset(key hash.Key) *Set {
	s, _ = s.Remove(key)
	return s
}

// Remove deletes the given hash.Key, if it exists and returns a new *Set
// reflecting that change and a true value indicating the hash.Key was found.
// If the hash.Key does not exist in the *Set then the original *Set is returned
// with a false value indicating that the hash.Key was not found.
func (s *Set) Remove(key hash.Key) (*Set, bool) {
	if s.numEnts == 0 {
		return s, false
	}

	var hv = key.Hash()
	var path, leaf, idx = s.find(hv)

	if leaf == nil {
		return s, false
	}

	var newLeaf, deleted = leaf.del(key)

	if !deleted {
		return s, false
	}

	var curTable = path.pop()
	var depth = uint(path.len())

	var ns = s.copy()

	ns.numEnts--

	if curTable == &s.root {
		//copying all s.root into ns.root already done in *ns = *s
		if newLeaf == nil { //leaf was a FlatLeaf
			ns.root.remove(idx)
		} else { //leaf was a CollisionLeaf
			ns.root.replace(idx, newLeaf)
		}
	} else {
		var newTable = curTable.copy()

		if newLeaf == nil { //leaf was a FlatLeaf
			newTable.remove(idx)

			// Side-Effects of removing a Key from the table
			var nents = newTable.numEntries()
			switch {
			case nents == 0:
				newTable = nil
			case nents == downgradeThreshold:
				newTable = downgradeToSparseTable(
					newTable.hash(), depth, newTable.entries())
			}
		} else { //leaf was a CollisionLeaf
			newTable.replace(idx, newLeaf)
		}

		ns.persist(curTable, newTable, path)
	}

	return ns, deleted
}

//func (s *Set) walk(fn visitFn) bool {
//	var err, keepOn = s.root.visit(fn, 0)
//	if err != nil {
//		panic(err)
//	}
//	return keepOn
//}

// Iter returns an *Iter structure. You can call the Next() method on the *Iter
// structure sucessively until it returns a nil key value to walk the keys in
// the Set data structure. This is safe under any usage of the *Set because the
// Set is immutable.
func (s *Set) Iter() *Iter {
	if s.NumEntries() == 0 {
		return nil
	}

	//log.Printf("s.Iter: s=\n%s", s.treeString(""))
	var it = newIter(&s.root)

	//find left-most leaf
LOOP:
	for {
		var curNode = it.tblNextNode()
		switch x := curNode.(type) {
		case nil:
			panic("finding first leaf; it.tblNextNode() returned nil")
		case tableI:
			it.stack.push(it.tblNextNode)
			it.tblNextNode = x.iter()
			assert(it.tblNextNode != nil, "it.tblNextNode==nil")
			break //switch
		case leafI:
			it.curLeaf = x
			break LOOP
		default:
			panic("finding first leaf; unknown type")
		}
	}

	return it
}

// Range applies the given function to every hash.Key in the *Set. If the
// function returns false the Range operation stops.
func (s *Set) Range(fn func(hash.Key) bool) {
	//var visitLeafs = func(n nodeI, depth uint) bool {
	//	if leaf, ok := n.(leafI); ok {
	//		for _, key := range leaf.keys() {
	//			if !fn(key) {
	//				return false
	//			}
	//		}
	//	}
	//	return true
	//} //end: visitLeafsFn = func(nodeI)
	//s.walk(visitLeafs)
	var it = s.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		if !fn(k) {
			break
		}
	}
}

// NumEntries returns the number of hash.Keys in the *Set. This operation is
// O(1) because the count is maintained at the top level for the *Set and does
// not require a walk of the *Set data structure to return the count.
func (s *Set) NumEntries() uint {
	return s.numEnts
}

// String prints a string representation of the Set. It is intended to be
// simmilar to fmt.Printf("%#v") of a golang set[].
func (s *Set) String() string {
	var ents = make([]string, s.NumEntries())
	var i int = 0
	s.Range(func(k hash.Key) bool {
		//log.Printf("i=%d, k=%#v\n", i, k)
		ents[i] = fmt.Sprintf("%#v", k)
		i++
		return true
	})
	return "Set{" + strings.Join(ents, ",") + "}"
}

// treeString returns a (potentially very large) string that represets the
// entire Set data structure.
func (s *Set) treeString(indent string) string {
	var str string

	str = indent +
		fmt.Sprintf("Set{ numEnts: %d, root:\n", s.numEnts)
	str += indent + s.root.treeString(indent, 0)
	str += indent + "}"

	return str
}

//type Stats struct {
//	DeepestKeys struct {
//		Keys  []hash.Key
//		Depth uint
//	}
//
//	// Depth of deepest table
//	MaxDepth uint
//
//	// TableCountsByNumEntries is a Hash table of the number of tables with each
//	// given number of entries in the tatble. There are slots for
//	// [0..IndexLimit] inclusive (so there are IndexLimit+1 slots). Technically,
//	// there should never be a table with zero entries, but I allow counting
//	// tables with zero entries just to catch those errors.
//	// [0..IndexLimit] inclusive
//	TableCountsByNumEntries [hash.IndexLimit + 1]uint
//
//	// TableCountsByDepth is a Hash table of the number of tables at a given
//	// depth. There are slots for [0..DepthLimit).
//	// [0..DepthLimit)
//	TableCountsByDepth [hash.DepthLimit]uint
//
//	// Nils is the total count of allocated slots that are unused in the Set.
//	Nils uint
//
//	// Nodes is the total count of nodeI capable structs in the Set.
//	Nodes uint
//
//	// Tables is the total count of tableI capable structs in the Set.
//	Tables uint
//
//	// Leafs is the total count of leafI capable structs in the Set.
//	Leafs uint
//
//	// FixedTables is the total count of fixedTable structs in the Set.
//	FixedTables uint
//
//	// SparseTables is the total count of sparseTable structs in the Set.
//	SparseTables uint
//
//	// FlatLeafs is the total count of flatLeaf structs in the Set.
//	FlatLeafs uint
//
//	// CollisionLeafs is the total count of collisionLeaf structs in the Set.
//	CollisionLeafs uint
//
//	// Keys is the total number of Keys in the Set.
//	Keys uint
//}
//
//// Stats walks the Hamt in a pre-order traversal and populates a Stats data
//// struture which it returns.
//func (s *Set) Stats() *Stats {
//	var stats = new(Stats)
//
//	// statFn closes over the stats variable
//	var statFn = func(n nodeI, depth uint) bool {
//		var keepOn = true
//		switch x := n.(type) {
//		case nil:
//			stats.Nils++
//			keepOn = false
//		case *fixedTable:
//			stats.Nodes++
//			stats.Tables++
//			stats.FixedTables++
//			stats.TableCountsByNumEntries[x.numEntries()]++
//			stats.TableCountsByDepth[x.depth]++
//			if x.depth > stats.MaxDepth {
//				stats.MaxDepth = x.depth
//			}
//		case *sparseTable:
//			stats.Nodes++
//			stats.Tables++
//			stats.SparseTables++
//			stats.TableCountsByNumEntries[x.numEntries()]++
//			stats.TableCountsByDepth[x.depth]++
//			if x.depth > stats.MaxDepth {
//				stats.MaxDepth = x.depth
//			}
//		case *flatLeaf:
//			stats.Nodes++
//			stats.Leafs++
//			stats.FlatLeafs++
//			stats.Keys += 1
//			keepOn = false
//		case *collisionLeaf:
//			stats.Nodes++
//			stats.Leafs++
//			stats.CollisionLeafs++
//			stats.Keys += uint(len(x.keys_))
//			keepOn = false
//		}
//		return keepOn
//	}
//
//	s.walk(statFn)
//	return stats
//}

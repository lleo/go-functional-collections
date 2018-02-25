// Package set implements a functional Set data structure. The internal data
// structure of set is a Hashed Array Mapped Trie
// (see https://en.wikipedia.org/wiki/Hash_array_mapped_trie).
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
	"sort"
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

// Set struct mainains an immutable collection of hash.Key entries.
type Set struct {
	root    tableI
	numEnts int
}

// New returns a properly initialized pointer to a Set struct.
func New() *Set {
	var s = new(Set)
	s.root = newRootTable()
	return s
}

func newRootTable() tableI {
	// fixedTable at root makes a noticable perf diff on small & large Maps.
	return newFixedTable(0, 0)
	//return newSparseTable(0, 0, 0)
}

// FIXME: generic version of newSparseTable & newFixedTable
func newTable(depth uint, hashVal hash.Val) tableI {
	//return newFixedTable(depth, hashVal)
	return newSparseTable(depth, hashVal, 0)
}

// FIXME: generic version of createSparseTable & createFixedTable
// FIXME: This should obviate createSparseTable & createFixedTable.
func createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if assertOn {
		assert(depth > 0, "createTable(): depth < 1")
		assertf(leaf1.hash().HashPath(depth) == leaf2.hash().HashPath(depth),
			"createTable(): hp1,%s != hp2,%s",
			leaf1.hash().HashPath(depth),
			leaf2.hash().HashPath(depth))
	}

	var retTable = newTable(depth, leaf1.hash())

	var idx1 = leaf1.hash().Index(depth)
	var idx2 = leaf2.hash().Index(depth)
	if idx1 != idx2 {
		retTable.insertInplace(idx1, leaf1)
		retTable.insertInplace(idx2, leaf2)
	} else { // idx1 == idx2
		var node nodeI
		if depth == hash.MaxDepth {
			node = newCollisionLeaf(append(leaf1.keys(), leaf2.keys()...))
		} else {
			node = createTable(depth+1, leaf1, leaf2)
		}
		retTable.insertInplace(idx1, node)
	}

	return retTable
}

// copy creates a shallow copy of the Set data structure and returns a pointer
// to that shallow copy.
func (s *Set) copy() *Set {
	var ns = new(Set)
	*ns = *s
	return ns
}

// IsSet searches the Set for a hash.Key value where the given key (k) matches
// a key in the Set (k0) such that k.Equals(k0) returns true. If the given
// key is found IsSet return true, otherwise it returns false.
func (s *Set) IsSet(key hash.Key) bool {
	if s.NumEntries() == 0 {
		return false
	}

	var hv = key.Hash()
	var curTable tableI = s.root

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
	var curTable tableI = s.root

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
	assert(s.root != nil, "s.root == nil")

	// downgrade() & upgrade() can return an unmodified table in NewFromList(),
	// BulkInsert(), BulkDelete(), and Merge(). Hence persist() is unnecessary.
	//if newTable == oldTable {
	//	return
	//}

	if s.root == oldTable {
		s.root = newTable
		return
	}

	var depth = uint(path.len()) //guaranteed depth > 0
	var parentDepth = depth - 1

	var parentIdx = oldTable.hash().Index(parentDepth)

	var oldParent = path.pop()
	var newParent tableI

	if newTable == nil {
		newParent = oldParent.remove(parentIdx)
	} else {
		newParent = oldParent.replace(parentIdx, newTable)
	}

	s.persist(oldParent, newParent, path)

	return
}

// Set inserts the give hash.Key into the Set and returns a new *Set. If an
// equivalent key exists in the receiver *Set nothing is done to the *Set and
// the original receiver *Set is returned.
//
// Equivalentcy of keys is determined by k.Equals(k0) where k is the given
// hash.Key and k0 is the hash.Key already stored in the *Set.
func (s *Set) Set(key hash.Key) *Set {
	var ns, _ = s.Add(key)
	return ns
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

	var newTable tableI

	if leaf == nil {
		newTable = curTable.insert(idx, newFlatLeaf(key))
		added = true
	} else {
		// This only happens when depth == MaxDepth
		var node nodeI
		if leaf.hash() != hv {
			// common case
			//node = createSparseTable(depth+1, leaf, newFlatLeaf(key))
			node = createTable(depth+1, leaf, newFlatLeaf(key))
			added = true
		} else {
			node, added = leaf.put(key)
		}

		newTable = curTable.replace(idx, node)
	}

	ns.persist(curTable, newTable, path)

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
	//if m.numEnts == 0 {
	//if m.root == nil {
	if s.NumEntries() == 0 {
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
	//var depth = uint(path.len())

	var ns = s.copy()

	ns.numEnts--
	if ns.numEnts < 0 {
		panic("WTF!?! new map.numEnts < 0")
	}

	var newTable tableI

	if newLeaf == nil {
		// leaf was a FlatLeaf
		newTable = curTable.remove(idx)
	} else {
		// leaf was a CollisionLeaf
		newTable = curTable.replace(idx, newLeaf)
	}

	ns.persist(curTable, newTable, path)

	return ns, deleted
}

func (s *Set) walkPreOrder(fn visitFunc) bool {
	return s.root.walkPreOrder(fn, 0)
}

// Iter returns an *Iter structure. You can call the Next() method on the *Iter
// structure sucessively until it returns a nil key value to walk the keys in
// the Set data structure. This is safe under any usage of the *Set because the
// Set is immutable.
func (s *Set) Iter() *Iter {
	if s.NumEntries() == 0 {
		return nil
	}

	//log.Printf("s.Iter: s=\n%s", s.treeString(""))
	var it = newIter(s.root)

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

// Keys returns a hash.Key slice that contains all the entries in Set.
func (s *Set) Keys() []hash.Key {
	var keys = make([]hash.Key, s.NumEntries())
	var i int
	s.Range(func(k hash.Key) bool {
		keys[i] = k
		i++
		return true
	})
	return keys
}

// NumEntries returns the number of hash.Keys in the *Set. This operation is
// O(1) because the count is maintained at the top level for the *Set and does
// not require a walk of the *Set data structure to return the count.
func (s *Set) NumEntries() int {
	return s.numEnts
}

// String prints a string representation of the Set. It is intended to be
// simmilar to fmt.Printf("%#v") of a golang set[].
func (s *Set) String() string {
	var ents = make([]string, s.NumEntries())
	var i int

	var it = s.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		ents[i] = fmt.Sprintf("%#v", k)
		i++
	}

	//s.Range(func(k hash.Key) bool {
	//	//log.Printf("i=%d, k=%#v\n", i, k)
	//	ents[i] = fmt.Sprintf("%#v", k)
	//	i++
	//	return true
	//})

	return "Set{" + strings.Join(ents, ",") + "}"
}

// TreeString returns a (potentially very large) string that represets the
// entire Set data structure.
func (s *Set) TreeString(indent string) string {
	var str string

	str = indent +
		fmt.Sprintf("Set{ numEnts: %d, root:\n", s.numEnts)
	str += indent + s.root.treeString(indent, 0)
	str += indent + "}"

	return str
}

// DeepCopy does a complete deep copy of a *Set returning an entirely new *Set.
func (s *Set) DeepCopy() *Set {
	var ns = s.copy()
	ns.root = s.root.deepCopy()
	return ns
}

// Equiv compares two *Set's by value.
func (s *Set) Equiv(s0 *Set) bool {
	//log.Printf("Set#Equiv: s.NumEntries(),%d != s0.NumEntries(),%d",
	//	s.NumEntries(), s0.NumEntries())

	if s.NumEntries() != s0.NumEntries() {
		return false
	}
	if !s.root.equiv(s0.root) {
		return false
	}
	return true
	//return s.NumEntries() == s0.NumEntries() && s.root.equiv(s0.root)
}

// Count recursively traverses the HAMT data structure to count every key.
func (s *Set) Count() int {
	return s.root.count()
}

// NewFromList constructs a new *Set structure containing all the keys
// of the given hash.Key slice.
//
// NewFromList is implemented more efficiently than repeated calls to Add.
func NewFromList(keys []hash.Key) *Set {
	var s = New()
	for _, k := range keys {
		var hv = k.Hash()
		var path, leaf, idx = s.find(hv)
		var curTable = path.pop()
		var depth = uint(path.len())
		var added bool
		if leaf == nil {
			curTable.insertInplace(idx, newFlatLeaf(k))
			//var _, isSparseTable = curTable.(*sparseTable)
			//if isSparseTable && curTable.slotsUsed() == upgradeThreshold {
			//if curTable.slotsUsed() == upgradeThreshold {
			if curTable.needsUpgrade() {
				var newTable = curTable.upgrade()
				s.persist(newTable, curTable, path)
			}
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k))
				node = createTable(depth+1, leaf, newFlatLeaf(k))
				added = true
			} else {
				node, added = leaf.put(k)
			}
			curTable.replaceInplace(idx, node)
		}
		if added {
			s.numEnts++
		}
	}
	return s
}

func insertPersist(
	s *Set,
	isOrigTable map[tableI]bool,
	k hash.Key,
) {
	var hv = k.Hash()
	var path, leaf, idx = s.find(hv)
	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool
	if isOrigTable[curTable] {
		var newTable tableI
		if leaf == nil {
			newTable = curTable.insert(idx, newFlatLeaf(k))
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k))
				node = createTable(depth+1, leaf, newFlatLeaf(k))
				added = true
			} else {
				node, added = leaf.put(k)
			}
			newTable = curTable.replace(idx, node)
		}
		s.persist(curTable, newTable, path)
	} else {
		if leaf == nil {
			curTable.insertInplace(idx, newFlatLeaf(k))
			//var _, isSparseTable = curTable.(*sparseTable)
			//if isSparseTable && curTable.slotsUsed() == upgradeThreshold {
			//if curTable.slotsUsed() == upgradeThreshold {
			if curTable.needsUpgrade() {
				var newTable = curTable.upgrade()
				s.persist(curTable, newTable, path)
			}
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k))
				node = createTable(depth+1, leaf, newFlatLeaf(k))
				added = true
			} else {
				node, added = leaf.put(k)
			}
			curTable.replaceInplace(idx, node)
		}
	}
	if added {
		s.numEnts++
	}
}

// BulkInsert stores all the given keys from the argument hash.Key slice  into
// the receiver Set.
//
// The returned Set maintains the structure sharing relationship with the
// receiver Set.
//
// BulkInsert is implemented more efficiently than repeated calls to Add.
func (s *Set) BulkInsert(keys []hash.Key) *Set {
	var isOrigTable = make(map[tableI]bool)
	s.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var ns = s.copy()
	for _, k := range keys {
		insertPersist(ns, isOrigTable, k)
	}
	return ns
}

// Merge returns a Set that contains all the entries from the receiver Set and
// the argument Set.
func (s *Set) Merge(other *Set) *Set {
	var big, sml = s, other
	if s.NumEntries() < other.NumEntries() {
		big, sml = other, s
	}
	// big is bigger then sml

	var isOrigTable = make(map[tableI]bool)
	big.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var ns = big.copy()
	var it = sml.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		insertPersist(ns, isOrigTable, k)
	}
	return ns
}

func removePersist(
	s *Set,
	isOrigTable map[tableI]bool,
	k hash.Key,
) bool {
	var hv = k.Hash()
	var path, leaf, idx = s.find(hv)
	if leaf == nil {
		return false // did not remove k
	}
	var newLeaf, found = leaf.del(k)
	if !found {
		return false // did not remove k
	}
	s.numEnts--
	var curTable = path.pop()
	if isOrigTable[curTable] {
		var newTable tableI
		if newLeaf == nil {
			newTable = curTable.remove(idx)
		} else {
			newTable = curTable.replace(idx, newLeaf)
		}
		s.persist(curTable, newTable, path)
	} else {
		if newLeaf == nil {
			if curTable.slotsUsed()-1 > 0 {
				curTable.removeInplace(idx)
				if curTable.needsDowngrade() {
					var newTable = curTable.downgrade()
					s.persist(curTable, newTable, path)
				}
			} else { // curTable.slotsUsed()-1 <= 0
				// we need to use persist cuz this will shrink empty tables
				var newTable = curTable.remove(idx)
				s.persist(curTable, newTable, path)
			}
		} else {
			curTable.replaceInplace(idx, newLeaf)
		}
	}
	return true // removed k
}

// BulkDelete removes all the keys in the given hash.Key slice. It returns a new
// persistent Set and a slice of the keys not found in the in the original Set.
//
// BulkDelete is implemented more efficiently than repeated calls to Remove.
func (s *Set) BulkDelete(keys []hash.Key) (*Set, []hash.Key) {
	var isOrigTable = make(map[tableI]bool)
	s.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var notFound []hash.Key
	var ns = s.copy()
	//KEYSLOOP:
	for _, k := range keys {
		if !removePersist(ns, isOrigTable, k) {
			notFound = append(notFound, k)
		}
	}
	return ns, notFound
}

// BulkDelete2 removes all the keys in the given hash.Key slice. It returns a
// new persistent Set.
//
// BulkDelete2 is implemented more efficiently than repeated calls to Remove.
func (s *Set) BulkDelete2(keys []hash.Key) *Set {
	var isOrigTable = make(map[tableI]bool)
	s.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var ns = s.copy()
	//KEYSLOOP:
	for _, k := range keys {
		removePersist(ns, isOrigTable, k)
	}
	return ns
}

// Union returns a Set that contains all entries of all given Sets.
//
// First it sorts all the sets from largest to smallest then progressively
// calculates the union of the biggest Set with each smaller Set.
//
//    var resultSet = sets[0] // biggest *Set
//    for _, s := range sets[1:] {
//      resultSet = resultSet.Union(s)
//    }
//    return resultSet
//
func Union(sets ...*Set) *Set {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	sort.Slice(sets, func(i, j int) bool {
		return sets[i].NumEntries() > sets[j].NumEntries()
	})
	// sets is now sorted from largest to smallest

	var resultSet = sets[0].copy()
	var rest = sets[1:]

	var isOrigTable = make(map[tableI]bool)
	resultSet.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	for _, s := range rest {
		var it = s.Iter()
		for k := it.Next(); k != nil; k = it.Next() {
			insertPersist(resultSet, isOrigTable, k)
		}
	}
	return resultSet
}

// Union returns a Set that contains all entries for the receiver Set and the
// argument Set.
func (s *Set) Union(other *Set) *Set {
	return s.Merge(other)
}

// Intersection returns a Set that is the Intersection  of all the given sets.
//
// First it sorts all the sets from largest to smallest then progressively
// calculates the intersection of the biggest Set with each smaller Set.
//
//    var resultSet = sets[0] // biggest *Set
//    for _, s := range sets[1:] {
//      resultSet = resultSet.Intersect(s)
//    }
//    return resultSet
//
func Intersection(sets ...*Set) *Set {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	sort.Slice(sets, func(i, j int) bool {
		return sets[i].NumEntries() > sets[j].NumEntries()
	})
	// sets is now sorted from largest to smallest

	var res = sets[0] // largest
	for _, s := range sets[1:] {
		res = res.Intersect(s)
	}
	return res
}

// Intersect returns a Set that contains only the entries that the receiver Set
// and the argument Set have in common.
//
// There is no structural sharing with either the receiver Set or argument Set.
func (s *Set) Intersect(other *Set) *Set {
	var intersectKeys []hash.Key
	var it = other.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		if s.IsSet(k) {
			intersectKeys = append(intersectKeys, k)
		}
	}
	return NewFromList(intersectKeys)
}

// Difference returns a new Set based on the receiver Set that contains none of
// the entries from the argument Set.
//
// Difference is implemented by repeated calls to Unset.
//
// NOTE: a.Difference(b) != b.Difference(a)
func (s *Set) Difference(other *Set) *Set {
	var diffSet = s
	var it = other.Iter()
	for k := it.Next(); k != nil; k = it.Next() {
		diffSet = diffSet.Unset(k)
	}
	return diffSet
}

// // Difference2 returns a new Set based on the receiver Set that contains none of
// // the entries from the argument Set.
// //
// // Difference2 is implemented by repeated calls to removePersist, the basic
// // function of BulkDelete.
// //
// // NOTE: a.Difference2(b) != b.Difference2(a)
// func (s *Set) Difference2(other *Set) *Set {
// 	var isOrigTable = make(map[tableI]bool)
// 	s.walkPreOrder(func(n nodeI, depth uint) bool {
// 		if t, isTable := n.(tableI); isTable {
// 			isOrigTable[t] = true
// 		}
// 		return true
// 	})
//
// 	var diffSet = s.copy()
// 	var it = other.Iter()
// 	for k := it.Next(); k != nil; k = it.Next() {
// 		removePersist(diffSet, isOrigTable, k)
// 	}
// 	return diffSet
// }

// // Difference1 returns a new Set based on the receiver Set that contains none of
// // the entries from the argument Set.
// //
// // Difference1 is calculated by creating a list of shared keys, then doing a
// // BulkDelete of those shared keys from the receiver Set.
// //
// // NOTE: a.Difference1(b) != b.Difference1(a)
// func (s *Set) Difference1(other *Set) *Set {
// 	var otherKeys = other.Keys()
// 	var diff, _ = s.BulkDelete(otherKeys)
// 	return diff
// }

// // Difference3 returns a new Set based on the receiver Set that contains none of
// // the entries from the argument Set.
// //
// // Difference3 is calculated by creating a list of shared keys, then doing a
// // BulkDelete of those shared keys from the receiver Set.
// //
// // NOTE: a.Difference3(b) != b.Difference3(a)
// func (s *Set) Difference3(other *Set) *Set {
// 	var otherKeys = other.Keys()
// 	var diff = s.BulkDelete2(otherKeys)
// 	return diff
// }

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
//		case *fixedTable:
//			stats.Nodes++
//			stats.Tables++
//			stats.FixedTables++
//			stats.TableCountsByNumEntries[x.slotsUsed()]++
//			stats.TableCountsByDepth[x.depth]++
//			if x.depth > stats.MaxDepth {
//				stats.MaxDepth = x.depth
//			}
//		case *sparseTable:
//			stats.Nodes++
//			stats.Tables++
//			stats.SparseTables++
//			stats.TableCountsByNumEntries[x.slotsUsed()]++
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

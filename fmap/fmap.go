package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/fmap/hash"
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

type Map struct {
	root    fixedTable
	numEnts uint
}

type MapKey interface {
	Hash() hash.HashVal
	Equals(MapKey) bool
	String() string
}

func New() *Map {
	return new(Map)
}

func (m *Map) copy() *Map {
	var nm = new(Map)
	*nm = *m
	return nm
}

func (m *Map) Get(key MapKey) interface{} {
	var v, _ = m.Load(key)
	return v
}

// Load retrieves the value related to the MapKey in the Map data structure.
// It also return a bool to indicate the value was found. This allows you to
// store nil values in the Map data structure.
func (m *Map) Load(key MapKey) (interface{}, bool) {
	if m.NumEntries() == 0 {
		return nil, false
	}

	var hv = key.Hash()
	var curTable tableI = &m.root

	var val interface{}
	var found bool

DepthIter:
	for depth := uint(0); depth <= hash.MaxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		switch n := curNode.(type) {
		case nil:
			val, found = nil, false
			break DepthIter
		case leafI:
			val, found = n.get(key)
			break DepthIter
		case tableI:
			curTable = n
		}
	}

	return val, found
}

// find() traverses the path defined by the given HashVal till it encounters
// a leafI, then it returns the table path leading to the current table (also
// returned) and the Index in the current table the leaf is at.
func (m *Map) find(hv hash.HashVal) (*tableStack, leafI, uint) {
	var curTable tableI = &m.root

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
func (m *Map) persist(oldTable, newTable tableI, path *tableStack) {
	// Removed the case where path.len() == 0 on the first call to nh.perist(),
	// because that case is handled in Put & Del now. It is handled in Put & Del
	// because otherwise we were allocating an extraneous fixedTable for the
	// old m.root.
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
		m.root = *oldParent.(*fixedTable)
		newParent = &m.root
	} else {
		newParent = oldParent.copy()
	}

	if newTable == nil {
		newParent.remove(parentIdx)
	} else {
		newParent.replace(parentIdx, newTable)
	}

	if path.len() > 0 {
		m.persist(oldParent, newParent, path)
	}

	return
}

func (m *Map) LoadOrStore(key MapKey, val interface{}) (
	*Map, interface{}, bool,
) {
	var hv = key.Hash()

	var path, leaf, idx = m.find(hv)
	var curTable = path.pop()

	var depth = uint(path.len())

	var foundVal interface{}
	var found bool
	var added bool //probably not necessary added == !found

	var nm *Map

	if curTable == &m.root {
		if leaf == nil {
			nm = m.copy()

			nm.root.insert(idx, newFlatLeaf(key, val))
		} else {
			foundVal, found = leaf.get(key)
			if found {
				return m, foundVal, true // result of Loaded value
			}
			//else

			nm = m.copy()

			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key, val)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
				added = true
			}
			nm.root.replace(idx, node)
		}
	} else {
		var newTable tableI

		if leaf == nil {
			if (curTable.numEntries() + 1) == upgradeThreshold {
				newTable = upgradeToFixedTable(
					curTable.hash(), depth, curTable.entries())
			} else {
				newTable = curTable.copy()
			}

			newTable.insert(idx, newFlatLeaf(key, val))
			added = true
		} else {
			foundVal, found = leaf.get(key)
			if found {
				return m, foundVal, true // result of Loaded value
			}
			//else
			newTable = curTable.copy()

			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key, val)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
				added = true
			}

			newTable.replace(idx, node)
		}

		nm = m.copy()

		nm.persist(curTable, newTable, path)
	}

	if added {
		nm.numEnts++
	}

	return nm, nil, false // result for a Stored value
}

func (m *Map) Put(key MapKey, val interface{}) *Map {
	m, _ = m.Store(key, val)
	return m
}

// Store stores a new (key,value) pair in the Map data structure. It returns the
// new *Map data structure and a bool indicating if a new pair was added (true)
// or if the value merely replaced a prior value (false).
func (m *Map) Store(key MapKey, val interface{}) (*Map, bool) {
	var nm = m.copy()

	var hv = key.Hash()

	var path, leaf, idx = m.find(hv)
	var curTable = path.pop()

	var depth = uint(path.len())

	var added bool

	if curTable == &m.root {
		// Special handling of root table.
		// 1) It is never upgraded.
		// 2) It doesn't need to be copied, cuz the new Map has a fixed root.
		if leaf == nil {
			// if curTable.get(idx) is slot is empty
			nm.root.insert(idx, newFlatLeaf(key, val))
			added = true
		} else {
			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key, val)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
				added = true
			}
			nm.root.replace(idx, node)
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

			newTable.insert(idx, newFlatLeaf(key, val))
			added = true
		} else {
			newTable = curTable.copy()

			var node nodeI
			if leaf.hash() == hv {
				node, added = leaf.put(key, val)
			} else {
				node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
				added = true
			}

			newTable.replace(idx, node)
		}

		nm.persist(curTable, newTable, path)
	}

	if added {
		nm.numEnts++
	}

	return nm, added
}

func (m *Map) Del(key MapKey) *Map {
	m, _, _ = m.Remove(key)
	return m
}

func (m *Map) Delete(key MapKey) *Map {
	return m.Del(key)
}

func (m *Map) Remove(key MapKey) (*Map, interface{}, bool) {
	if m.numEnts == 0 {
		return m, nil, false
	}

	var hv = key.Hash()
	var path, leaf, idx = m.find(hv)

	if leaf == nil {
		return m, nil, false
	}

	var newLeaf, val, deleted = leaf.del(key)

	if !deleted {
		return m, nil, false
	}

	var curTable = path.pop()
	var depth = uint(path.len())

	var nm = m.copy()

	nm.numEnts--

	if curTable == &m.root {
		//copying all m.root into nm.root already done in *nm = *m
		if newLeaf == nil { //leaf was a FlatLeaf
			nm.root.remove(idx)
		} else { //leaf was a CollisionLeaf
			nm.root.replace(idx, newLeaf)
		}
	} else {
		var newTable = curTable.copy()

		if newLeaf == nil { //leaf was a FlatLeaf
			newTable.remove(idx)

			// Side-Effects of removing a KeyVal from the table
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

		nm.persist(curTable, newTable, path)
	}

	return nm, val, deleted
}

func (m *Map) walk(fn visitFn) bool {
	var err, keepOn = m.root.visit(fn, 0)
	if err != nil {
		panic(err)
	}
	return keepOn
}

//func (m *Map) walk(fn visitFn) error {
//	var curTable tableI = &m.root
//
//	for idx := uint(0); idx < hash.IndexLimit; idx++ {
//		var n = curTable.get(idx)
//
//		switch x := n.(type) {
//		case nil:
//
//		case leafI:
//
//		case tableI:
//
//		}
//	}
//}

func (m *Map) Iter() *Iter {
	if m.NumEntries() == 0 {
		return nil
	}

	//log.Printf("m.Iter: m=\n%s", m.treeString(""))
	var it = newIter(&m.root)

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

func (m *Map) Range(fn func(MapKey, interface{}) bool) {
	//var visitLeafs = func(n nodeI, depth uint) bool {
	//	if leaf, ok := n.(leafI); ok {
	//		for _, kv := range leaf.keyVals() {
	//			if !fn(kv.Key, kv.Val) {
	//				return false
	//			}
	//		}
	//	}
	//	return true
	//} //end: visitLeafsFn = func(nodeI)
	//m.walk(visitLeafs)
	var it = m.Iter()
	for k, v := it.Next(); k != nil; k, v = it.Next() {
		if !fn(k, v) {
			break
		}
	}
}

func (m *Map) NumEntries() uint {
	return m.numEnts
}

// String prints a string representation of the Map. It is intended to be
// simmilar to fmt.Printf("%#v") of a golang map[].
func (m *Map) String() string {
	var ents = make([]string, m.NumEntries())
	var i int = 0
	m.Range(func(k MapKey, v interface{}) bool {
		//log.Printf("i=%d, k=%#v, v=%#v\n", i, k, v)
		ents[i] = fmt.Sprintf("%#v:%#v", k, v)
		i++
		return true
	})
	return "Map{" + strings.Join(ents, ",") + "}"
}

// treeString returns a (potentially very large) string that represets the
// entire Map data structure. It is for print debugging.
func (m *Map) treeString(indent string) string {
	var str string

	str = indent +
		fmt.Sprintf("Map{ numEnts: %d, root:\n", m.numEnts)
	str += indent + m.root.treeString(indent, 0)
	str += indent + "}"

	return str
}

type Stats struct {
	DeepestKeys struct {
		Keys  []MapKey
		Depth uint
	}

	// Depth of deepest table
	MaxDepth uint

	// TableCountsByNumEntries is a Hash table of the number of tables with each
	// given number of entries in the tatble. There are slots for
	// [0..IndexLimit] inclusive (so there are IndexLimit+1 slots). Technically,
	// there should never be a table with zero entries, but I allow counting
	// tables with zero entries just to catch those errors.
	// [0..IndexLimit] inclusive
	TableCountsByNumEntries [hash.IndexLimit + 1]uint

	// TableCountsByDepth is a Hash table of the number of tables at a given
	// depth. There are slots for [0..DepthLimit).
	// [0..DepthLimit)
	TableCountsByDepth [hash.DepthLimit]uint

	// Nils is the total count of allocated slots that are unused in the Map.
	Nils uint

	// Nodes is the total count of nodeI capable structs in the Map.
	Nodes uint

	// Tables is the total count of tableI capable structs in the Map.
	Tables uint

	// Leafs is the total count of leafI capable structs in the Map.
	Leafs uint

	// FixedTables is the total count of fixedTable structs in the Map.
	FixedTables uint

	// SparseTables is the total count of sparseTable structs in the Map.
	SparseTables uint

	// FlatLeafs is the total count of flatLeaf structs in the Map.
	FlatLeafs uint

	// CollisionLeafs is the total count of collisionLeaf structs in the Map.
	CollisionLeafs uint

	// KeyVals is the total number of KeyVal pairs int the Map.
	KeyVals uint
}

// Stats walks the Hamt in a pre-order traversal and populates a Stats data
// struture which it returns.
func (m *Map) Stats() *Stats {
	var stats = new(Stats)

	// statFn closes over the stats variable
	var statFn = func(n nodeI, depth uint) bool {
		var keepOn = true
		switch x := n.(type) {
		case nil:
			stats.Nils++
			keepOn = false
		case *fixedTable:
			stats.Nodes++
			stats.Tables++
			stats.FixedTables++
			stats.TableCountsByNumEntries[x.numEntries()]++
			stats.TableCountsByDepth[x.depth]++
			if x.depth > stats.MaxDepth {
				stats.MaxDepth = x.depth
			}
		case *sparseTable:
			stats.Nodes++
			stats.Tables++
			stats.SparseTables++
			stats.TableCountsByNumEntries[x.numEntries()]++
			stats.TableCountsByDepth[x.depth]++
			if x.depth > stats.MaxDepth {
				stats.MaxDepth = x.depth
			}
		case *flatLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.FlatLeafs++
			stats.KeyVals += 1
			keepOn = false
		case *collisionLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.CollisionLeafs++
			stats.KeyVals += uint(len(x.kvs))
			keepOn = false
		}
		return keepOn
	}

	m.walk(statFn)
	return stats
}

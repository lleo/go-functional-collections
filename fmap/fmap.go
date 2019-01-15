// Package fmap implements a functional Map data structure; mapping a key to a
// value. The package name cannot be map because that is a reserved keyword in
// Golang, so "fmap" was used instead. The internal data structure of fmap is a
// Hashed Array Mapped Trie
// (see https://en.wikipedia.org/wiki/Hash_array_mapped_trie).
//
// Functional means that each data structure is immutable and persistent.
// The Map is immutable because you never modify a Map in place, but rather
// every modification (like a Store or Remove) creates a new Map with that
// modification. This is not as inefficient as it sounds like it would be. Each
// modification only changes the smallest  branch of the data structure it needs
// to in order to effect the new mapping. Otherwise, the new data structure
// shares the majority of the previous data structure. That is the persistent
// property.
//
// Each method call that potentially modifies the Map, returns a new Map data
// structure in addition to the other pertinent return values.
//
// Every key in the key/value mapping must implement the key.Hash interface.
//
// Any value can be stored in the key/value mapping, because values are treated
// and returned as interface{} values. That means the values returned from
// methods Get, Load, and LoadOrStore, must be type asserted back to their
// original type by the user of this library.
package fmap

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/key"
	"github.com/lleo/go-functional-collections/key/hash"
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

// The Map struct maintains a immutable collection of key/value mappings.
type Map struct {
	root    tableI
	numEnts int
}

// New returns a properly initialize pointer to a fmap.Map struct.
func New() *Map {
	var m = new(Map)
	m.root = newRootTable()
	return m
}

func newRootTable() tableI {
	// fixedTable at root makes a noticable perf diff on small & large Maps.
	return newFixedTable(0, 0)
	//return newSparseTable(0, 0, 0)
}

// newTable is a generic version of newSparseTable & newFixedTable
func newTable(depth uint, hashVal hash.Val) tableI {
	//return newFixedTable(depth, hashVal)
	return newSparseTable(depth, hashVal, 0)
}

// createTable is a  generic version of createSparseTable & createFixedTable
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
			node = newCollisionLeaf(append(leaf1.keyVals(), leaf2.keyVals()...))
		} else {
			node = createTable(depth+1, leaf1, leaf2)
		}
		retTable.insertInplace(idx1, node)
	}

	return retTable
}

// copy creates a shallow copy of the Map data structure and returns a pointer
// to that shallow copy.
func (m *Map) copy() *Map {
	var nm = new(Map)
	*nm = *m
	return nm
}

// Get loads the value stored for the given key. If the key doesn't exist in the
// Map a nil is returned. If you need to store nil values and want to
// distinguish between a found existing mapping of the key to nil and a
// non-existent mapping for the key, you must use the Load method.
func (m *Map) Get(key key.Hash) interface{} {
	var v, _ = m.Load(key)
	return v
}

// Load retrieves the value related to the key.Hash in the Map data structure.
// It also return a bool to indicate the value was found. This allows you to
// store nil values in the Map data structure and distinguish between a found
// nil key/value mapping and a non-existant key/value mapping.
func (m *Map) Load(key key.Hash) (interface{}, bool) {
	if m.NumEntries() == 0 {
		return nil, false
	}

	var hv = key.Hash()
	var curTable = m.root

	var val interface{}
	var found bool

DepthIter:
	for depth := uint(0); depth <= hash.MaxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) // nodeI

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

// find() traverses the path defined by the given Val till it encounters
// a leafI, then it returns the table path leading to the current table (also
// containted in the *tableStack path) and the Index in the current table the
// leaf is at.
func (m *Map) find(hv hash.Val) (*tableStack, leafI, uint) {
	var curTable = m.root

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
	_ = assertOn && assert(m.root != nil, "m.root == nil")

	// downgrade() & upgrade() can return an unmodified table in NewFromList(),
	// BulkInsert(), BulkDelete(), and Merge(). Hence persist() is unnecessary.
	if newTable == oldTable {
		return
	}

	if m.root == oldTable {
		m.root = newTable
		return
	}

	var depth = uint(path.len()) // guaranteed depth > 0
	var parentDepth = depth - 1

	var parentIdx = oldTable.hash().Index(parentDepth)

	var oldParent = path.pop()
	var newParent tableI

	if newTable == nil {
		newParent = oldParent.remove(parentIdx)
	} else {
		newParent = oldParent.replace(parentIdx, newTable)
	}

	m.persist(oldParent, newParent, path)
}

// LoadOrStore returns the existing value for the key if present. Otherwise,
// it stores a new key/value mapping and returns the given value. The loaded
// result is true if the key/value was loaded, false if a new key/value mapping
// was created. Lastly, if an existing key/value mapping was loaded then the
// returned map is the original *Map, if the a new key/value mapping was
// created returned *Map is a new persistent *Map.
func (m *Map) LoadOrStore(key key.Hash, val interface{}) (
	*Map, interface{}, bool,
) {
	var hv = key.Hash()

	var path, leaf, idx = m.find(hv)
	var curTable = path.pop()

	var depth = uint(path.len())

	var foundVal interface{}
	var found bool
	var added bool // probably not necessary added == !found

	var nm *Map

	var newTable tableI

	if leaf == nil {
		newTable = curTable.insert(idx, newFlatLeaf(key, val))
		added = true
	} else {
		foundVal, found = leaf.get(key)
		if found {
			return m, foundVal, true // result of Loaded value
		}
		// else

		var node nodeI
		if leaf.hash() != hv {
			// common case
			//node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
			node = createTable(depth+1, leaf, newFlatLeaf(key, val))
			added = true
		} else {
			// hash collision; very rare case; leaf.hash() == key.Hash()
			node, added = leaf.put(key, val)
		}

		newTable = curTable.replace(idx, node)
	}

	nm = m.copy()

	nm.persist(curTable, newTable, path)

	if added {
		nm.numEnts++
	}

	return nm, nil, false // result for a Stored value
}

// Put stores a new key/value mapping. It returns a new persistent *Map data
// structure.
func (m *Map) Put(key key.Hash, val interface{}) *Map {
	m, _ = m.Store(key, val)
	return m
}

// Store stores a new key/value mapping. It returns a new persistent
// *Map data structure and a bool indicating if a new pair was added (true)
// or if the value merely replaced a prior value (false). Regardless of
// whether a new key/value mapping was created or mearly replaced, a new
// *Map is created.
func (m *Map) Store(key key.Hash, val interface{}) (*Map, bool) {
	var nm = m.copy()

	var hv = key.Hash()

	var path, leaf, idx = nm.find(hv)
	var curTable = path.pop()

	var depth = uint(path.len())

	var added bool

	var newTable tableI

	if leaf == nil {
		newTable = curTable.insert(idx, newFlatLeaf(key, val))
		added = true
	} else {
		// This only happens when depth == MaxDepth
		var node nodeI
		if leaf.hash() != hv {
			// common case
			//node = createSparseTable(depth+1, leaf, newFlatLeaf(key, val))
			node = createTable(depth+1, leaf, newFlatLeaf(key, val))
			added = true
		} else {
			// hash collision; very rare case; leaf.hash() == key.Hash()
			node, added = leaf.put(key, val)
		}

		newTable = curTable.replace(idx, node)
	}

	nm.persist(curTable, newTable, path)

	if added {
		nm.numEnts++
	}

	return nm, added
}

// Del deletes any entry with the given key, but does not indicate if the key
// existed or not. However, if the key did not exist the returned *Map will be
// the original *Map.
func (m *Map) Del(key key.Hash) *Map {
	m, _, _ = m.Remove(key)
	return m
}

// Remove deletes any key/value mapping for the given key. It returns a
// *Map data structure, the possible value that was stored for that key,
// and a boolean idicating if the key was found and deleted. If the key didn't
// exist, then the value is set nil, and the original *Map is returned.
func (m *Map) Remove(key key.Hash) (*Map, interface{}, bool) {
	//if m.numEnts == 0 {
	//if m.root == nil {
	if m.NumEntries() == 0 {
		return m, nil, false
	}

	var hv = key.Hash()
	var path, leaf, idx = m.find(hv)

	if leaf == nil {
		//log.Println("leaf == nil")
		return m, nil, false
	}

	var newLeaf, val, deleted = leaf.del(key)

	if !deleted {
		return m, nil, false
	}

	var curTable = path.pop()
	//var depth = uint(path.len())

	var nm = m.copy()

	nm.numEnts--
	if nm.numEnts < 0 {
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

	nm.persist(curTable, newTable, path)

	return nm, val, deleted
}

func (m *Map) walkPreOrder(fn visitFunc) bool {
	return m.root.walkPreOrder(fn, 0)
}

// Iter returns an *Iter structure. You can call the Next() method on the *Iter
// structure sucessively until it return a nil key value, to walk the key/value
// mappings in the Map data structure. This is safe under any usage of the *Map
// because the Map is immutable.
func (m *Map) Iter() *Iter {
	//if m.NumEntries() == 0 {
	//	return nil
	//}

	var it = newIter(m.root)

	// find left-most leaf
LOOP:
	for {
		var curNode = it.tblNextNode()
		switch x := curNode.(type) {
		case nil:
			// EMPTY FMAP
			// current 'it' is in end state, so just return.
			break LOOP
		case tableI:
			it.stack.push(it.tblNextNode)
			it.tblNextNode = x.iter()
			_ = assertOn && assert(it.tblNextNode != nil, "it.tblNextNode==nil")
			break // switch
		case leafI:
			it.curLeaf = x
			break LOOP
		default:
			panic("finding first leaf; unknown type")
		}
	}

	return it
}

// Range applies the given function for every key/value mapping in the *Map
// data structure. Given that the *Map is immutable there is no danger with
// concurrent use of the *Map while the Range method is executing.
func (m *Map) Range(fn func(KeyVal) bool) {
	var it = m.Iter()
	for kv := it.Next(); kv.Key != nil; kv = it.Next() {
		if !fn(kv) {
			break
		}
	}
}

// NumEntries returns the number of key/value entries in the *Map. This
// operation is O(1), because a current count of the number of entries is
// maintained at the top level of the *Map data structure, so walking the data
// structure is not required to get the current count of key/value entries.
func (m *Map) NumEntries() int {
	return m.numEnts
}

// String prints a string list all the key/value mappings in the *Map. It is
// intended to be simmilar to fmt.Printf("%#v") of a golang builtin map.
func (m *Map) String() string {
	var ents = make([]string, m.NumEntries())
	var i int

	var it = m.Iter()
	for kv := it.Next(); kv.Key != nil; kv = it.Next() {
		ents[i] = fmt.Sprintf("%#v:%#v", kv.Key, kv.Val)
		i++
	}

	//m.Range(func(kv KeyVal) bool {
	//	ents[i] = fmt.Sprintf("%#v:%#v", kv.Key, kv.Val)
	//	i++
	//	return true
	//})

	return "Map{" + strings.Join(ents, ",") + "}"
}

// TreeString returns a (potentially very large) string that represets the
// entire Map data structure. It is for print debugging.
func (m *Map) TreeString(indent string) string {
	var str string

	str = indent +
		fmt.Sprintf("Map{ numEnts: %d, root:\n", m.numEnts)
	str += indent + m.root.treeString(indent, 0)
	str += indent + "}"

	return str
}

// DeepCopy does a complete deep copy of a *Map returning an entirely new *Map.
func (m *Map) DeepCopy() *Map {
	var nm = m.copy()
	nm.root = m.root.deepCopy()
	return nm
}

// Equiv compares two *Map's by value.
func (m *Map) Equiv(m0 *Map) bool {
	if m.NumEntries() != m0.NumEntries() {
		return false
	}
	if !m.root.equiv(m0.root) {
		return false
	}
	return true
	//return m.NumEntries() == m0.NumEntries() && m.root.equiv(m0.root)
}

// Count recursively traverses the HAMT data structure to count every key,value
// pair.
func (m *Map) Count() int {
	return m.root.count()
}

// NewFromList constructs a new *Map structure containing all the key,value
// pairs of the given KeyVal slice.
func NewFromList(kvs []KeyVal) *Map {
	var m = New()
	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		var hv = k.Hash()
		var path, leaf, idx = m.find(hv)
		var curTable = path.pop()
		var depth = uint(path.len())
		var added bool
		if leaf == nil {
			curTable.insertInplace(idx, newFlatLeaf(k, v))
			//var _, isSparseTable = curTable.(*sparseTable)
			//if isSparseTable && curTable.slotsUsed() == upgradeThreshold {
			//if curTable.slotsUsed() == upgradeThreshold {
			if curTable.needsUpgrade() {
				var newTable = curTable.upgrade()
				m.persist(curTable, newTable, path)
			}
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k, v))
				node = createTable(depth+1, leaf, newFlatLeaf(k, v))
				added = true
			} else {
				node, added = leaf.put(k, v)
			}
			curTable.replaceInplace(idx, node)
		}
		if added {
			m.numEnts++
		}
	}
	return m
}

func insertPersist(
	m *Map,
	isOrigTable map[tableI]bool,
	resolve ResolveConflictFunc,
	k key.Hash,
	v interface{},
) {
	var hv = k.Hash()
	var path, leaf, idx = m.find(hv)
	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool
	if isOrigTable[curTable] {
		var newTable tableI
		if leaf == nil {
			newTable = curTable.insert(idx, newFlatLeaf(k, v))
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k, v))
				node = createTable(depth+1, leaf, newFlatLeaf(k, v))
				added = true
			} else {
				node, added = leaf.putResolve(k, v, resolve)
			}
			newTable = curTable.replace(idx, node)
		}
		m.persist(curTable, newTable, path)
	} else {
		if leaf == nil {
			curTable.insertInplace(idx, newFlatLeaf(k, v))
			//var _, isSparseTable = curTable.(*sparseTable)
			//if isSparseTable && curTable.slotsUsed() == upgradeThreshold {
			//if curTable.slotsUsed() == upgradeThreshold {
			if curTable.needsUpgrade() {
				var newTable = curTable.upgrade()
				m.persist(curTable, newTable, path)
			}
			added = true
		} else {
			var node nodeI
			if leaf.hash() != hv {
				//node = createSparseTable(depth+1, leaf, newFlatLeaf(k, v))
				node = createTable(depth+1, leaf, newFlatLeaf(k, v))
				added = true
			} else {
				node, added = leaf.putResolve(k, v, resolve)
			}
			curTable.replaceInplace(idx, node)
		}
	}
	if added {
		m.numEnts++
	}
}

// BulkInsert stores all the given key,value pairs into the Map while
// resolving any conflict with the given resolve function.
//
// The returned Map maintains the structure sharing relationship with the
// original Map.
//
// BulkInsert is implemented more efficiently than repeated calls to Store.
func (m *Map) BulkInsert(kvs []KeyVal, resolve ResolveConflictFunc) *Map {
	var isOrigTable = make(map[tableI]bool)
	m.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var nm = m.copy()
	for _, kv := range kvs {
		var k, v = kv.Key, kv.Val
		insertPersist(nm, isOrigTable, resolve, k, v)
	}
	return nm
}

// Merge inserts all the key,value pairs from the Map provided as an argument.
//
// If the argument Map has a key that is Equal to a key in the receiver Map,
// then the receiver key, its corrosponding value, and the value for the Equal
// key in the argument Map, will be passed into the ResolveConflictFunc; the
// result of which will be stored as the new key,value pair in the resulting
// Map.
func (m *Map) Merge(om *Map, resolve ResolveConflictFunc) *Map {
	var isOrigTable = make(map[tableI]bool)
	m.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var nm = m.copy()
	var it = om.Iter()
	for kv := it.Next(); kv.Key != nil; kv = it.Next() {
		insertPersist(nm, isOrigTable, resolve, kv.Key, kv.Val)
	}
	return nm
}

// BulkDelete removes all the keys in the given key.Hash slice. It then returns
// a new persistent Map and a slice of the keys not found in the in the
// original Map. BulkDelete is implemented more efficiently than repeated calls
// to Remove.
func (m *Map) BulkDelete(keys []key.Hash) (*Map, []key.Hash) {
	var isOrigTable = make(map[tableI]bool)
	m.walkPreOrder(func(n nodeI, depth uint) bool {
		if t, isTable := n.(tableI); isTable {
			isOrigTable[t] = true
		}
		return true
	})

	var notFound []key.Hash
	var nm = m.copy()
KEYSLOOP:
	for _, k := range keys {
		var hv = k.Hash()
		var path, leaf, idx = nm.find(hv)
		if leaf == nil {
			notFound = append(notFound, k)
			continue KEYSLOOP
		}
		var newLeaf, _, found = leaf.del(k)
		if !found {
			notFound = append(notFound, k)
			continue KEYSLOOP
		}
		nm.numEnts--
		var curTable = path.pop()
		if isOrigTable[curTable] {
			var newTable tableI
			if newLeaf == nil {
				newTable = curTable.remove(idx)
			} else {
				newTable = curTable.replace(idx, newLeaf)
			}
			nm.persist(curTable, newTable, path)
		} else {
			if newLeaf == nil {
				if curTable.slotsUsed()-1 > 0 {
					curTable.removeInplace(idx)
					//if curTable.slotsUsed() == downgradeThreshold {
					if curTable.needsDowngrade() {
						var newTable = curTable.downgrade()
						nm.persist(curTable, newTable, path)
					}
				} else { // curTable.slotsUsed()-1 <= 0
					// we need to use persist cuz this will shrink empty tables
					var newTable = curTable.remove(idx)
					nm.persist(curTable, newTable, path)
				}
			} else {
				curTable.replaceInplace(idx, newLeaf)
			}
		}
	}
	return nm, notFound
}

//type Stats struct {
//	DeepestKeys struct {
//		Keys  []key.Hash
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
//	// Nils is the total count of allocated slots that are unused in the Map.
//	Nils uint
//
//	// Nodes is the total count of nodeI capable structs in the Map.
//	Nodes uint
//
//	// Tables is the total count of tableI capable structs in the Map.
//	Tables uint
//
//	// Leafs is the total count of leafI capable structs in the Map.
//	Leafs uint
//
//	// FixedTables is the total count of fixedTable structs in the Map.
//	FixedTables uint
//
//	// SparseTables is the total count of sparseTable structs in the Map.
//	SparseTables uint
//
//	// FlatLeafs is the total count of flatLeaf structs in the Map.
//	FlatLeafs uint
//
//	// CollisionLeafs is the total count of collisionLeaf structs in the Map.
//	CollisionLeafs uint
//
//	// KeyVals is the total number of KeyVal pairs int the Map.
//	KeyVals uint
//}
//
//// Stats walks the Hamt in a pre-order traversal and populates a Stats data
//// struture which it returns.
//func (m *Map) Stats() *Stats {
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
//			stats.KeyVals += 1
//			keepOn = false
//		case *collisionLeaf:
//			stats.Nodes++
//			stats.Leafs++
//			stats.CollisionLeafs++
//			stats.KeyVals += uint(len(x.kvs))
//			keepOn = false
//		}
//		return keepOn
//	}
//
//	m.walkPreOrder(statFn)
//	return stats
//}

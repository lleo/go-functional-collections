package set

import (
	"fmt"
	"log"
	"strings"

	"github.com/lleo/go-functional-collections/key/hash"
)

// sparseTableInitCap constant sets the default capacity of a new
// sparseTable.
const sparseTableInitCap int = 2

type sparseTable struct {
	nodes    []nodeI
	depth    uint
	hashPath hash.Val
	nodeMap  bitmap
}

func newSparseTable(depth uint, hashVal hash.Val, size uint) *sparseTable {
	_ = assertOn && assertf(size <= hash.IndexLimit, "size,%d <= hash.IndexLimit", size)
	var t = new(sparseTable)
	var capacity int
	if size < hash.IndexLimit {
		capacity = pow2GreaterThan(size)
	} else {
		capacity = hash.IndexLimit
	}
	t.nodes = make([]nodeI, size, capacity)
	t.depth = depth
	t.hashPath = hashVal.HashPath(depth)
	return t
}

func (t *sparseTable) copy() tableI {
	var nt = new(sparseTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nodeMap = t.nodeMap

	nt.nodes = make([]nodeI, len(t.nodes), cap(t.nodes))
	copy(nt.nodes, t.nodes)

	return nt
}

func (t *sparseTable) deepCopy() tableI {
	var nt = new(sparseTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nodeMap = t.nodeMap

	nt.nodes = make([]nodeI, len(t.nodes), cap(t.nodes))
	//for i := 0; i < len(t.nodes); i++ {
	//	if table, isTable := t.nodes[i].(tableI); isTable {
	//		nt.nodes[i] = table.deepCopy()
	//	} else {
	//		// leafI's are functional, so no need to copy them.
	//		// nils can be copied just fine; duh!
	//		nt.nodes[i] = t.nodes[i]
	//	}
	//}
	for i, n := range t.nodes {
		switch x := n.(type) {
		case tableI:
			nt.nodes[i] = x.deepCopy()
		case leafI:
			// leafI's are functional, so no need to copy them.
			nt.nodes[i] = x
		case nil:
			panic("found a nil entry in sparseTable")
			// nils can be copied just fine; duh!
			nt.nodes[i] = x
		default:
			panic("unknown entry in table")
		}
	}

	return nt
}

// equiv compares the *sparseTable to another node by value. This ultimately
// becomes a deep comparison of tables.
func (t *sparseTable) equiv(other nodeI) bool {
	var ot, ok = other.(*sparseTable)
	if !ok {
		log.Println("other is not a *sparseTable")
		return false
	}
	ok = ok && t.depth == ot.depth
	ok = ok && t.hashPath == ot.hashPath
	ok = ok && t.nodeMap == ot.nodeMap
	ok = ok && len(t.nodes) == len(ot.nodes)
	if !ok {
		log.Printf("t,%s != ot,%s", t, ot)
		return false
	}
	for i, n := range t.nodes {
		if n == nil && n != ot.nodes[i] {
			log.Printf("n == nil && n != ot.nodes[%d],%s", i, ot.nodes[i])
			return false
		}
		if n != nil && !n.equiv(ot.nodes[i]) {
			log.Printf("!n.equiv(ot.nodes[%d])", i)
			return false
		}
	}
	return true
}

func createSparseTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if assertOn {
		assert(depth > 0, "createSparseTable(): depth < 1")
		assertf(leaf1.hash().HashPath(depth) == leaf2.hash().HashPath(depth),
			"createSparseTable(): hp1,%s != hp2,%s",
			leaf1.hash().HashPath(depth),
			leaf2.hash().HashPath(depth))
	}

	var retTable = newSparseTable(depth, leaf1.hash(), 0)

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
			node = createSparseTable(depth+1, leaf1, leaf2)
		}
		retTable.insertInplace(idx1, node)
	}

	return retTable
}

// downgradeToSparseTable() converts fixedTable structs that have less than
// or equal to downgradeThreshold tableEntry's. One important thing we know is
// that none of the entries will collide with another.
//
// The ents []tableEntry slice is guaranteed to be in order from lowest idx to
// highest. tableI.entries() also adhears to this contract.
func downgradeToSparseTable(
	hashPath hash.Val,
	depth uint,
	ents []tableEntry,
) *sparseTable {
	var nt = new(sparseTable)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]nodeI, len(ents), len(ents)+1)

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		nt.nodeMap.set(ent.idx)
		nt.nodes[i] = ent.node
	}

	return nt
}

// hash returns an incomplete hash of this table. Any levels past it's current
// depth should be zero.
func (t *sparseTable) hash() hash.Val {
	return t.hashPath
}

// String return a string representation of this table including the hashPath,
// depth, and number of entries.
func (t *sparseTable) String() string {
	return fmt.Sprintf("sparseTable{hashPath:%s, depth=%d, slotsUsed()=%d}",
		t.hashPath.HashPathString(t.depth), t.depth, t.slotsUsed())
}

// treeString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *sparseTable) treeString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent +
		fmt.Sprintf("sparseTable{hashPath=%s, depth=%d, slotsUsed()=%d,",
			t.hashPath.HashPathString(depth), t.depth, t.slotsUsed())

	strs[1] = indent + "\tnodeMap=" + t.nodeMap.String() + ","

	for i, n := range t.nodes {
		var idx = n.hash().Index(depth)
		if t, isTable := n.(tableI); isTable {
			strs[2+i] = indent +
				fmt.Sprintf("\tt.nodes[%d]:\n%s",
					idx, t.treeString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", idx, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *sparseTable) slotsUsed() uint {
	if t == nil {
		log.Printf("t,%#p.slotsUsed()=0", t)
		return 0
	}
	return uint(len(t.nodes))
	//return t.nodeMap.count(hash.IndexLimit)
}

func (t *sparseTable) entries() []tableEntry {
	var n = t.slotsUsed()
	var ents = make([]tableEntry, n)

	for j := uint(0); j < n; j++ {
		idx := t.nodes[j].hash().Index(t.depth)
		ents[j] = tableEntry{idx, t.nodes[j]}
	}

	return ents
}

func (t *sparseTable) get(idx uint) nodeI {
	if !t.nodeMap.isSet(idx) {
		return nil
	}

	var j = t.nodeMap.count(idx)

	return t.nodes[j]
}

func (t *sparseTable) needsUpgrade() bool {
	return t.slotsUsed() == upgradeThreshold
}

func (t *sparseTable) needsDowngrade() bool {
	return false
}

func (t *sparseTable) upgrade() tableI {
	var nt = newFixedTable(t.depth, t.hashPath)
	var slots = t.slotsUsed()
	for j := uint(0); j < slots; j++ {
		var idx0 = t.nodes[j].hash().Index(t.depth)
		nt.insertInplace(idx0, t.nodes[j])
	}
	return nt
}

func (t *sparseTable) downgrade() tableI {
	//panic("downgrade() invalid op")
	return t
}

func (t *sparseTable) insertInplace(idx uint, n nodeI) {
	var j = int(t.nodeMap.count(idx))
	if j == len(t.nodes) {
		t.nodes = append(t.nodes, n)
	} else {
		// slower and more obscure method
		//t.nodes = append(t.nodes[:j], append([]nodeI{n}, t.nodes[j:]...)...)

		// faster and more understandable method
		t.nodes = append(t.nodes, nodeI(nil))
		copy(t.nodes[j+1:], t.nodes[j:])
		t.nodes[j] = n
	}
	t.nodeMap.set(idx)
}

func (t *sparseTable) insert(idx uint, n nodeI) tableI {
	_ = assertOn && assert(!t.nodeMap.isSet(idx),
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	//if t.slotsUsed()+1 == upgradeThreshold {
	//	var nt = t.upgrade()
	//	nt.insertInplace(idx, n)
	//	return nt
	//}

	var nt = t.copy()
	nt.insertInplace(idx, n)
	return nt
}

func (t *sparseTable) replaceInplace(idx uint, n nodeI) {
	var j = t.nodeMap.count(idx)
	t.nodes[j] = n
}

func (t *sparseTable) replace(idx uint, n nodeI) tableI {
	_ = assertOn && assert(t.nodeMap.isSet(idx),
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	var nt = t.copy()
	nt.replaceInplace(idx, n)
	//var nt = t.copy().(*sparseTable)
	//var j = nt.nodeMap.count(idx)
	//nt.nodes[j] = n
	return nt
}

func (t *sparseTable) removeInplace(idx uint) {
	var j = int(t.nodeMap.count(idx))
	if j == len(t.nodes)-1 {
		t.nodes = t.nodes[:j]
	} else {
		// No obvious performance difference, but append code is more obvious
		t.nodes = append(t.nodes[:j], t.nodes[j+1:]...)
		//t.nodes = t.nodes[:j+copy(t.nodes[j:], t.nodes[j+1:])]
	}
	t.nodeMap.unset(idx)
}

func (t *sparseTable) remove(idx uint) tableI {
	_ = assertOn && assert(t.nodeMap.isSet(idx),
		"t.remove(idx) where idx slot is already empty")

	// If the table is the root table (ie t.depth == 0), do NOT return nil.
	// If the table only has one entry, cut to the chase and return nil.
	if t.depth > 0 && t.slotsUsed() == 1 {
		return nil
	}

	var nt = t.copy()
	nt.removeInplace(idx)
	//var nt = t.copy().(*sparseTable)
	//var j = int(nt.nodeMap.count(idx))
	//if j == len(nt.nodes)-1 {
	//	nt.nodes = nt.nodes[:j]
	//} else {
	//	// No obvious performance difference, but append code is more obvious
	//	nt.nodes = append(nt.nodes[:j], nt.nodes[j+1:]...)
	//	//nt.nodes = nt.nodes[:j+copy(nt.nodes[j:], nt.nodes[j+1:])]
	//}
	//nt.nodeMap.unset(idx)
	return nt
}

// walkPreOrder executes the visitFunc in pre-order traversal. If there is no
// node for a given idx, walkPreOrder skips that idx.
//
// The traversal stops if the visitFunc function returns false.
func (t *sparseTable) walkPreOrder(fn visitFunc, depth uint) bool {
	_ = assertOn && assertf(depth == t.depth, "depth,%d != t.depth=%d; t=%s", depth, t.depth, t)

	depth++

	if !fn(t, depth) {
		return false
	}

	for j, n := range t.nodes {
		_ = assertOn && assertf(n != nil, "n == nil; j=%d", j)
		if !n.walkPreOrder(fn, depth) {
			return false
		}
	}

	return true
}

func (t *sparseTable) iter() tableIterFunc {
	var j = -1

	return func() nodeI {
		if j < len(t.nodes)-1 {
			j++
			return t.nodes[j]
		}
		return nil
	}
}

func (t *sparseTable) count() int {
	var i int
	for _, n := range t.nodes {
		i += n.count()
	}
	return i
}

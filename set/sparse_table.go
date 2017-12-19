package set

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/set/hash"
)

// sparseTableInitCap constant sets the default capacity of a new
// sparseTable.
const sparseTableInitCap int = 2

type sparseTable struct {
	nodes    []nodeI
	depth    uint
	hashPath hash.HashVal
	nodeMap  bitmap
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
	for i := 0; i < len(t.nodes); i++ {
		if table, isTable := t.nodes[i].(tableI); isTable {
			nt.nodes[i] = table.deepCopy()
		} else {
			//leafI's are functional, so no need to copy them.
			//nils can be copied just fine; duh!
			nt.nodes[i] = t.nodes[i]
		}
	}
	//for i, n := range t.nodes {
	//	switch x := n.(type) {
	//	case tableI:
	//		nt.nodes[i] = x.deepCopy()
	//	default:
	//		nt.nodes[i] = x
	//	}
	//}

	return nt
}

func createSparseTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if assertOn {
		assert(depth > 0, "createSparseTable(): depth < 1")
		assertf(leaf1.hash().HashPath(depth) == leaf2.hash().HashPath(depth),
			"createSparseTable(): hp1,%s != hp2,%s",
			leaf1.hash().HashPath(depth),
			leaf2.hash().HashPath(depth))
	}

	var retTable = new(sparseTable)
	retTable.hashPath = leaf1.hash().HashPath(depth)
	retTable.depth = depth
	//retTable.nodeMap = 0
	retTable.nodes = make([]nodeI, 0, sparseTableInitCap)

	var idx1 = leaf1.hash().Index(depth)
	var idx2 = leaf2.hash().Index(depth)
	if idx1 != idx2 {
		retTable.insert(idx1, leaf1)
		retTable.insert(idx2, leaf2)
	} else { //idx1 == idx2
		var node nodeI
		if depth == hash.MaxDepth {
			node = newCollisionLeaf(append(leaf1.keys(), leaf2.keys()...))
		} else {
			node = createSparseTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
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
	hashPath hash.HashVal,
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
func (t *sparseTable) hash() hash.HashVal {
	return t.hashPath
}

// String return a string representation of this table including the hashPath,
// depth, and number of entries.
func (t *sparseTable) String() string {
	return fmt.Sprintf("sparseTable{hashPath:%s, depth=%d, numEntries()=%d}",
		t.hashPath.HashPathString(t.depth), t.depth, t.numEntries())
}

// LongString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *sparseTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent +
		fmt.Sprintf("sparseTable{hashPath=%s, depth=%d, numEntries()=%d,",
			t.hashPath.HashPathString(depth), t.depth, t.numEntries())

	strs[1] = indent + "\tnodeMap=" + t.nodeMap.String() + ","

	for i, n := range t.nodes {
		var idx = n.hash().Index(depth)
		if t, isTable := n.(tableI); isTable {
			strs[2+i] = indent +
				fmt.Sprintf("\tt.nodes[%d]:\n%s",
					idx, t.LongString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", idx, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *sparseTable) numEntries() uint {
	return uint(len(t.nodes))
	//return t.nodeMap.count(hash.IndexLimit)
}

func (t *sparseTable) entries() []tableEntry {
	var n = t.numEntries()
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

func (t *sparseTable) insert(idx uint, n nodeI) {
	_ = assertOn && assert(!t.nodeMap.isSet(idx),
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	var j = int(t.nodeMap.count(idx))
	if j == len(t.nodes) {
		t.nodes = append(t.nodes, n)
	} else {
		// Second code is significantly faster
		// Also I believe the second code is more understandable.

		//t.nodes = append(t.nodes[:j], append([]nodeI{n}, t.nodes[j:]...)...)

		t.nodes = append(t.nodes, nodeI(nil))
		copy(t.nodes[j+1:], t.nodes[j:])
		t.nodes[j] = n
	}

	t.nodeMap.set(idx)
}

func (t *sparseTable) replace(idx uint, n nodeI) {
	_ = assertOn && assert(t.nodeMap.isSet(idx),
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	var j = t.nodeMap.count(idx)
	t.nodes[j] = n
}

func (t *sparseTable) remove(idx uint) {
	_ = assertOn && assert(t.nodeMap.isSet(idx),
		"t.remove(idx) where idx slot is already empty")

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

// visit executes the visitFn in pre-order traversal. If there is no node for
// a given node, slot visit calls the visitFn on nil.
//
// The traversal stops if the visitFn function returns false.
func (t *sparseTable) visit(fn visitFn, depth uint) (error, bool) {
	if depth != t.depth {
		var err = fmt.Errorf("depth,%d != t.depth=%d; t=%s", depth, t.depth, t)
		return err, false
	}

	depth++

	if !fn(t, depth) {
		return nil, false
	}

	for idx := uint(0); idx < hash.IndexLimit; idx++ {
		var n = t.get(idx)
		if n == nil {
			if !fn(n, depth) {
				return nil, false
			}
		} else if err, keepOn := n.visit(fn, depth); !keepOn || err != nil {
			return err, keepOn
		}
	}

	return nil, true
}

func (t *sparseTable) iter() tableIterFunc {
	var j int = -1

	return func() nodeI {
		if j < len(t.nodes)-1 {
			j++
			return t.nodes[j]
		}
		return nil
	}
}
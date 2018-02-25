package set

import (
	"fmt"
	"log"
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

type fixedTable struct {
	nodes     [hash.IndexLimit]nodeI
	depth     uint
	usedSlots uint //numEnts  uint
	hashPath  hash.Val
}

func newFixedTable(depth uint, hashVal hash.Val) *fixedTable {
	var t = new(fixedTable)
	t.depth = depth
	t.hashPath = hashVal.HashPath(depth)
	return t
}

func (t *fixedTable) copy() tableI {
	var nt = new(fixedTable)
	*nt = *t
	return nt
}

func (t *fixedTable) deepCopy() tableI {
	var nt = new(fixedTable)

	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.usedSlots = t.usedSlots

	//for i := 0; i < len(t.nodes); i++ {
	//	if table, isTable := t.nodes[i].(tableI); isTable {
	//		nt.nodes[i] = table.deepCopy()
	//	} else {
	//		// leafs are functional, so no need to copy
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
			// nils can be copied just fine; duh!
			nt.nodes[i] = x
		default:
			panic("unknown entry in table")
		}
	}

	return nt
}

// equiv compares the *fixedTable to another node by value. This ultimately
// becomes a deep comparison of tables.
func (t *fixedTable) equiv(other nodeI) bool {
	var ot, ok = other.(*fixedTable)
	if !ok {
		log.Println("other is not a *fixedTable")
		return false
	}
	if t.depth != ot.depth {
		log.Printf("t.depth,%d != ot.depth,%d", t.depth, ot.depth)
		return false
	}
	if t.usedSlots != ot.usedSlots {
		log.Printf("t.usedSlots,%d != ot.usedSlots,%d", t.usedSlots, ot.usedSlots)
		return false
	}
	if t.hashPath != ot.hashPath {
		log.Printf("t.hashPath,%s != ot.hashPath,%s", t.hashPath, ot.hashPath)
		return false
	}
	//ok = ok && t.depth == ot.depth
	//ok = ok && t.usedSlots == ot.usedSlots
	//ok = ok && t.hashPath == ot.hashPath
	//if !ok {
	//	log.Printf("t,%s != ot,%s", t, ot)
	//	return false
	//}
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

func createFixedTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if assertOn {
		assertf(depth > 0, "createFixedTable(): depth,%d < 1", depth)
		assertf(leaf1.hash().HashPath(depth) == leaf2.hash().HashPath(depth),
			"createFixedTable(): hp1,%s != hp2,%s",
			leaf1.hash().HashPath(depth),
			leaf2.hash().HashPath(depth))
	}

	var retTable = newFixedTable(depth, leaf1.hash())

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
			node = createFixedTable(depth+1, leaf1, leaf2)
		}
		retTable.insertInplace(idx1, node)
	}

	return retTable
}

func upgradeToFixedTable(
	hashPath hash.Val,
	depth uint,
	ents []tableEntry,
) *fixedTable {
	var ft = new(fixedTable)

	ft.hashPath = hashPath
	ft.depth = depth
	ft.usedSlots = uint(len(ents))

	for _, ent := range ents {
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

// hash returns an incomplete hash of this table. Any levels past it's current
// depth should be zero.
func (t *fixedTable) hash() hash.Val {
	return t.hashPath
}

// String return a string representation of this table including the hashPath,
// depth, and number of entries.
func (t *fixedTable) String() string {
	return fmt.Sprintf("fixedTable{hashPath=%s, depth=%d, slotsUsed()=%d}",
		t.hashPath.HashPathString(t.depth), t.depth, t.slotsUsed())
}

// treeString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *fixedTable) treeString(indent string, depth uint) string {
	var strs = make([]string, 3+t.slotsUsed())

	strs[0] = indent + "fixedTable{"
	strs[1] = indent + fmt.Sprintf("\thashPath=%s, depth=%d, slotsUsed()=%d,",
		t.hashPath.HashPathString(depth), t.depth, t.slotsUsed())

	var j = 0
	for i, n := range t.nodes {
		if t.nodes[i] != nil {
			if t, isTable := t.nodes[i].(tableI); isTable {
				strs[2+j] = indent + fmt.Sprintf("\tnodes[%d]:\n", i) +
					t.treeString(indent+"\t", depth+1)
			} else {
				strs[2+j] = indent + fmt.Sprintf("\tnodes[%d]: %s", i, n)
			}
			j++
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *fixedTable) slotsUsed() uint {
	if t == nil {
		log.Printf("t,%#p.slotsUsed()=0", t)
		return 0
	}
	return t.usedSlots
}

func (t *fixedTable) entries() []tableEntry {
	var n = t.slotsUsed()
	var ents = make([]tableEntry, n)
	var i, j uint
	for i, j = 0, 0; j < n && i < hash.IndexLimit; i++ {
		if t.nodes[i] != nil {
			ents[j] = tableEntry{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fixedTable) get(idx uint) nodeI {
	return t.nodes[idx]
}

func (t *fixedTable) needsUpgrade() bool {
	return false
}

func (t *fixedTable) needsDowngrade() bool {
	return t.slotsUsed() == downgradeThreshold
}

func (t *fixedTable) upgrade() tableI {
	//panic("upgrade() invalid op")
	return t
}

func (t *fixedTable) downgrade() tableI {
	var nt = newSparseTable(t.depth, t.hashPath, t.slotsUsed())
	for idx := uint(0); idx < hash.IndexLimit; idx++ {
		if t.nodes[idx] != nil {
			nt.insertInplace(idx, t.nodes[idx])
		}
	}
	return nt
}

func (t *fixedTable) insertInplace(idx uint, n nodeI) {
	t.nodes[idx] = n
	t.usedSlots++
}

func (t *fixedTable) insert(idx uint, n nodeI) tableI {
	_ = assertOn && assert(t.nodes[idx] == nil,
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	var nt = t.copy().(*fixedTable)
	nt.nodes[idx] = n
	nt.usedSlots++
	return nt
}

func (t *fixedTable) replaceInplace(idx uint, n nodeI) {
	t.nodes[idx] = n
}

func (t *fixedTable) replace(idx uint, n nodeI) tableI {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	var nt = t.copy().(*fixedTable)
	nt.replaceInplace(idx, n)
	//nt.nodes[idx] = n
	return nt
}

func (t *fixedTable) removeInplace(idx uint) {
	t.nodes[idx] = nil
	t.usedSlots--
}

func (t *fixedTable) remove(idx uint) tableI {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.remove(idx) where idx slot is already empty")

	if t.depth > 0 { //non-root table
		if t.slotsUsed() == 1 {
			return nil
		}

		if t.slotsUsed()-1 == downgradeThreshold {
			var nt = t.downgrade()
			nt.removeInplace(idx)
			return nt
		}
	}

	var nt = t.copy().(*fixedTable)
	nt.removeInplace(idx)
	//nt.nodes[idx] = nil
	//nt.usedSlots--
	return nt
}

// walkPreOrder executes the visitFunc in pre-order traversal. If there is no
// node for a given idx, walkPreOrder skips that idx.
//
// The traversal stops if the visitFunc function returns false.
func (t *fixedTable) walkPreOrder(fn visitFunc, depth uint) bool {
	_ = assertOn && assertf(depth == t.depth, "depth,%d != t.depth=%d; t=%s", depth, t.depth, t)

	depth++

	if !fn(t, depth) {
		return false
	}

	for _, n := range t.nodes {
		if n != nil {
			if !n.walkPreOrder(fn, depth) {
				return false
			}
		}
	}

	return true
}

func (t *fixedTable) iter() tableIterFunc {
	var i = -1

	return func() nodeI {
		for i < int(hash.IndexLimit-1) {
			i++
			if t.nodes[i] != nil {
				return t.nodes[i]
			}
		}

		return nil
	}
}

func (t *fixedTable) count() int {
	var i int
	for _, n := range t.nodes {
		if n != nil {
			i += n.count()
		}
	}
	return i
}

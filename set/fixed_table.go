package set

import (
	"fmt"
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

type fixedTable struct {
	nodes    [hash.IndexLimit]nodeI
	depth    uint
	numEnts  uint
	hashPath hash.Val
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
	nt.numEnts = t.numEnts

	for i := 0; i < len(t.nodes); i++ {
		if table, isTable := t.nodes[i].(tableI); isTable {
			nt.nodes[i] = table.deepCopy()
		} else {
			//leafs are functional, so no need to copy
			//nils can be copied just fine; duh!
			nt.nodes[i] = t.nodes[i]
		}
	}

	return nt
}

//func createRootFixedTable(lf leafI) tableI {
//	var idx = lf.hash().Index(0)
//
//	var ft = new(fixedTable)
//	//ft.hashPath = 0
//	//ft.depth = 0
//	//ft.numEnts = 0
//	ft.insert(idx, lf)
//
//	return ft
//}

func createFixedTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if assertOn {
		assertf(depth > 0, "createFixedTable(): depth,%d < 1", depth)
		assertf(leaf1.hash().HashPath(depth) == leaf2.hash().HashPath(depth),
			"createFixedTable(): hp1,%s != hp2,%s",
			leaf1.hash().HashPath(depth),
			leaf2.hash().HashPath(depth))
	}

	var retTable = new(fixedTable)

	retTable.hashPath = leaf1.hash().HashPath(depth)
	retTable.depth = depth

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
			node = createFixedTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
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
	ft.numEnts = uint(len(ents))

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
	return fmt.Sprintf("fixedTable{hashPath=%s, depth=%d, numEntries()=%d}",
		t.hashPath.HashPathString(t.depth), t.depth, t.numEntries())
}

// treeString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *fixedTable) treeString(indent string, depth uint) string {
	var strs = make([]string, 3+t.numEntries())

	strs[0] = indent + "fixedTable{"
	strs[1] = indent + fmt.Sprintf("\thashPath=%s, depth=%d, numEntries()=%d,",
		t.hashPath.HashPathString(depth), t.depth, t.numEntries())

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

func (t *fixedTable) numEntries() uint {
	return t.numEnts
}

func (t *fixedTable) entries() []tableEntry {
	var n = t.numEntries()
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

func (t *fixedTable) insert(idx uint, n nodeI) {
	_ = assertOn && assert(t.nodes[idx] == nil,
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	t.nodes[idx] = n
	t.numEnts++
}

func (t *fixedTable) replace(idx uint, n nodeI) {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	t.nodes[idx] = n
}

func (t *fixedTable) remove(idx uint) {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.remove(idx) where idx slot is already empty")

	t.nodes[idx] = nil
	t.numEnts--
}

// visit executes the visitFn in pre-order traversal. If there is no node for
// a given node slot, visit calls the visitFn on nil.
//
// The traversal stops if the visitFn function returns false.
func (t *fixedTable) visit(fn visitFn, depth uint) (bool, error) {
	if depth != t.depth {
		var err = fmt.Errorf("depth,%d != t.depth=%d; t=%s", depth, t.depth, t)
		return false, err
	}

	depth++

	if !fn(t, depth) {
		return false, nil
	}

	for _, n := range t.nodes {
		if n == nil {
			if !fn(n, depth) {
				return false, nil
			}
		} else if keepOn, err := n.visit(fn, depth); !keepOn || err != nil {
			return keepOn, err
		}
	}

	return true, nil
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

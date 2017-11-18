package sorted_map

import (
	"fmt"
	"log"
	"strings"
)

//type Map interface {
//	Get(MapKey) interface{}
//	Load(MapKey) (interface{}, bool)
//	LoadOrStore(MapKey, interface{}) (Map, interface{}, bool)
//	Put(MapKey, interface{}) Map
//	Store(MapKey, interface{}) (Map, bool)
//	Del(MapKey) Map
//	Delete(MapKey) Map
//	Remove(MapKey) (Map, interface{}, bool)
//	Range(func(MapKey, interface{}) bool)
//	NumEntries() uint
//	String() string
//	LongString(string) string
//	Stats() *Stats
//}

type Map struct {
	root    *node
	numEnts uint
}

func New() *Map {
	var m = new(Map)
	return m
}

func (m *Map) copy() *Map {
	var nm = new(Map)
	*nm = *m
	return nm
}

func (m *Map) iterAll() *nodeIter {
	return m.iterRange(ninf, pinf)
}

func (m *Map) iterRange(startKey, endKey MapKey) *nodeIter {
	var cur, path = m.root.findNodeWithPath(startKey)
	var dir = less(startKey, endKey)
	if cur == nil {
		cur = path.pop()
		if dir {
			if less(cur.key, startKey) {
				cur = path.pop()
			}
		} else {
			if less(endKey, cur.key) {
				cur = path.pop()
			}
		}
	}
	return newNodeIter(cur, endKey, path)
}

func (m *Map) Get(k MapKey) interface{} {
	var v, _ = m.Load(k)
	return v
}

func (m *Map) Load(k MapKey) (interface{}, bool) {
	var n = m.root.findNode(k)

	if n == nil {
		return nil, false
	}

	return n.val, true
}

// LoadOrStore() finds the value for a given key. If the key is found then
// it simply return the current map, the value found, and a true value
// indicating is was found. If the key is NOT found then it stores the
// new key:value pair and returns the new map value, a nil for the previous
// value, and a false value indicating the key was not found and the store
// occured.
func (m *Map) LoadOrStore(k MapKey, v interface{}) (*Map, interface{}, bool) {
	//MAYBE:var n, path = m.root.findNodeWithPath(k)
	//MAYBE:
	//MAYBE:var nm *Map
	//MAYBE:if n == nil {
	//MAYBE:	nm = m.store(k, v, n, path)
	//MAYBE:	return nm, nil, false
	//MAYBE:}
	//MAYBE:
	//MAYBE:return m, n.val, true
	panic("not implemented")
	return nil, nil, false
}

func (m *Map) Put(k MapKey, v interface{}) *Map {
	var nm, _ = m.Store(k, v)
	return nm
}

// Store() inserts a new key:val pair and returns a new Map and a boolean
// indicatiing if the key:val was added(true) or merely replaced(false).
func (m *Map) Store(k MapKey, v interface{}) (*Map, bool) {
	var n, path = m.root.findNodeWithPath(k)

	var nm = m.copy()

	if n != nil {
		nm.replace(k, v, n, path)
		return nm, false
	}

	nm.insert(k, v, path)
	nm.numEnts++

	return nm, true
}

// replace() simply replaces the value on a copy of the node that contains the
// old value, then calls persist() on the *Map.
//
// replace() MUST be called on a new *Map.
func (m *Map) replace(k MapKey, v interface{}, on *node, path *nodeStack) {
	assert(cmp(on.key, k) == 0, "on.key != nn.key")

	var nn = on.copy()
	nn.val = v

	m.persistAll(on, nn, path) //ignoring return vals

	return
}

// insert() inserts a new node, create from a key-value pair,  into the tip of
// the path, then balances and persists the *Map.
//
// path MUST be non-zero length.
//
// insert() MUST be called on a new *Map.
func (m *Map) insert(k MapKey, v interface{}, path *nodeStack) {
	var on *node // = nil
	var nn = newNode(k, v)

	//m.insertBalance(on, nn, path)

	m.insertRepair(on, nn, path)

	return
}

// persistAll() calls persit all the way to root
func (m *Map) persistAll(on, nn *node, path *nodeStack) {
	log.Printf("persistAll:\non=%s\nnn=%s\npath=%s", on, nn, path)
	m.persistTill(on, nn, nil, path) //ignoring return vals
	return
}

// persistTill() updates node pointers down the path to the root creating a
// persistent data structure.
//
// persistTill() MUST be called on a new *Map.
func (m *Map) persistTill(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("persistTill:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)
	// on is the old node
	// nn is the new node
	if path.len() == 0 {
		log.Printf("persistTill: SETTING ROOT:\n"+
			"OLD m.root=%p; m.root=%v\nNEW ROOT nn=%s\n", m.root, m.root, nn)
		//if nn.isRed() {
		//	log.Println("persistAll: ROOT: ************ nn.isRed ************")
		//	nn.setBlack()
		//}
		m.root = nn
		return m.root, nn, path
	}
	//path.peek() != nil

	var oparent = path.pop()
	var nparent = oparent.copy()

	if on != nil {
		if on.isLeftChildOf(oparent) {
			nparent.ln = nn
		} else {
			nparent.rn = nn
		}
	} else {
		//NOTE: oparent.key == nn.key CAN NOT happen, cuz reasons...
		if less(nn.key, oparent.key) {
			nparent.ln = nn
		} else {
			nparent.rn = nn
		}
	}

	if oparent == term {
		return oparent, nparent, path
	}

	log.Printf("persistTill: recursing:\noparent=%s\nnparent=%s\npath=%s\n",
		oparent, nparent, path)
	return m.persistTill(oparent, nparent, term, path)
}

// rotateLeft() takes the target node(n) and its parent(p). We are rotating on
// target node(n) left.
//
//            p                      p
//            |                      |
//            n                      r
//          /   \      --->        /   \
//        l      r               n      y
//              / \             / \
//             x   y           l   x
//
// We are returning the original target node(n) and its new parent(r), because
// they changed position and swapped their parent-child relationship.
// Only p, n, and r changed values.
func (m *Map) rotateLeft(n, p *node) (*node, *node) {
	//assert(n == p.rn, "new node is not right child of new parent")
	//assert(p == nil || n == p.ln,
	//	"node is not the left child of parent")
	var r = n.rn
	//log.Printf("rotateLeft:\nn = %s\np = %s\nr = %s\n", n, p, r)

	if p != nil {
		if n.isLeftChildOf(p) {
			p.ln = r
		} else {
			p.rn = r
		}
	} else {
		m.root = r
	}

	n.rn = r.ln //handle anticipated orphaned node
	r.ln = n    //now orphan it

	return n, r
}

// rotateRight() takes the target node(n) and its parent(p). We are rotating on
// target node(n) right.
//
//            p                      p
//            |                      |
//            n                      l
//          /   \      --->        /   \
//        l      r               x      n
//       / \                           / \
//      x   y                         y   r
//
// We are returning the original target node(n) and its new parent(l), because
// they changed position and swapped their parent-child relationship.
// Only p, n, and l changed values.
func (m *Map) rotateRight(n, p *node) (*node, *node) {
	//assert(l == n.ln, "new node is not left child of new parent")
	//assert(p == nil || n == p.rn,
	//	"node is not the right child of parent")
	var l = n.ln
	//log.Printf("rotateRight:\nn = %s\np = %s\nl = %s\n", n, p, l)

	if p != nil {
		if n.isLeftChildOf(p) {
			p.ln = l
		} else {
			p.rn = l
		}
	} else {
		m.root = l
	}

	n.ln = l.rn //handle anticipated orphaned node
	l.rn = n    //now orphan it

	return n, l
}

// insertRepair() MUST be called on a new *Map.
func (m *Map) insertRepair(on, nn *node, path *nodeStack) {
	var parent, gp, uncle *node

	parent = path.peek()

	gp = path.peekN(1) // peek() == peekN(0); peekN is index from top

	if gp != nil {
		if parent.isLeftChildOf(gp) {
			uncle = gp.rn
		} else {
			uncle = gp.ln
		}
	}

	if parent == nil {
		m.insertCase1(on, nn, path)
	} else if parent.isBlack() {
		// we know:
		// parent exists and is black
		m.insertCase2(on, nn, path)
	} else if uncle.isRed() {
		// we know:
		// parent.isRed becuase of the previous condition
		// grandparent exists because root is never red
		// grandparent is black because parent is red
		m.insertCase3(on, nn, path)
	} else {
		// we know:
		// the nn side is longer than the uncle side because
		// parent.isRed and nn.isRed and uncle.isBlack.
		m.insertCase4(on, nn, path)
	}
}

// insertCase1() MUST be called on a new *Map.
func (m *Map) insertCase1(on, nn *node, path *nodeStack) {
	assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")
	assert(m.root == on, "path.peek()==nil BUT m.root != on")

	nn.setBlack()
	m.persistAll(on, nn, path)
	return
}

// insertCase2() MUST be called on a new *Map.
func (m *Map) insertCase2(on, nn *node, path *nodeStack) {
	m.persistAll(on, nn, path)
	return
}

// insertCase3() MUST be called on a new *Map.
func (m *Map) insertCase3(on, nn *node, path *nodeStack) {
	var parent = path.pop()
	var gp = path.pop() //gp means grandparent
	var uncle *node
	if parent.isLeftChildOf(gp) {
		uncle = gp.rn
	} else {
		uncle = gp.ln
	}

	var nparent = parent.copy() //new parent, cuz I am mutating it.
	nparent.setBlack()

	if less(nn.key, parent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	var nuncle = uncle.copy() //new uncle, cuz I am mutating it.
	nuncle.setBlack()

	var ngp = gp.copy() //new grandparent, cuz I am mutating it.
	ngp.setRed()

	if parent.isLeftChildOf(gp) {
		ngp.ln = nparent
		ngp.rn = nuncle
	} else {
		ngp.ln = nuncle
		ngp.rn = nparent
	}

	if parent.isLeftChildOf(gp) {
		ngp.ln = nparent
		ngp.rn = nuncle
	} else {
		ngp.ln = nuncle
		ngp.rn = nparent
	}

	m.insertRepair(gp, ngp, path)
	return
}

// insertCase4() MUST be called on a new *Map.
func (m *Map) insertCase4(on, nn *node, path *nodeStack) {
	var parent = path.pop()
	var gp = path.pop() //gp means grandparent

	var nparent = parent.copy()
	if less(nn.key, parent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	var ngp = gp.copy()
	if parent.isLeftChildOf(gp) {
		ngp.ln = nparent
	} else {
		ngp.rn = nparent
	}

	// insert_case4.1: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	if nn.isRightChildOf(nparent) && nparent.isLeftChildOf(ngp) {
		nn, nparent = m.rotateLeft(nparent, ngp)
	} else if nn.isLeftChildOf(nparent) && nparent.isRightChildOf(ngp) {
		nn, nparent = m.rotateRight(nparent, ngp)
	}

	// insert_case4.2: balancing rotate
	nparent.setBlack()
	ngp.setRed()

	if nn.isLeftChildOf(nparent) {
		nparent, ngp = m.rotateRight(ngp, path.peek())
	} else {
		nparent, ngp = m.rotateLeft(ngp, path.peek())
	}

	m.persistAll(gp, ngp, path)
	return
}

// insertBalance() was the what insertRepair() does but in one large function.
// which is annotated with commentary for me to understand the insert-n-balance
// a reb-black tree algorithm.
//
// insertBalance() rebalances the red-black tree from a given node to
// the root.
//
// insertBalance() MUST be called on a new *Map.
func (m *Map) insertBalance(on, nn *node, path *nodeStack) {

	if path.peek() == nil {
		// INSERT CASE #1
		//log.Println("insertBalance: path.peek() == nil")
		// path.peek() is the parent and the only way it can be nil is if
		// path.len()==0 and m.root == nil.
		assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")
		assert(m.root == on, "path.peek()==nil BUT m.root != on")

		nn.setBlack() //to enforce RBT#2

		m.persistAll(on, nn, path)
		return
	}
	//Fact#2: parent (aka path.peek()) != nil
	//log.Println("insertBalance: path.peek() != nil")

	if path.peek().isBlack() {
		// INSERT CASE #2
		//log.Printf("insertBalance: insert_case #2: path.peek().isBlack()")
		m.persistAll(on, nn, path) //persist will stitch nn into parent->child
		return
	}

	// INSERT CASE #3

	var parent = path.pop()

	//Fact#3: parent.isRed() == true
	// Fact#1 && Fact#3 violate RBT#4,
	// so we have to fix this with rotations.

	//Fact#4: parent has a parent; aka the grandparent exists
	// This is because of Fact#3(parent is RED) & RBT#2(root MUST be black).
	// Reasoning: If there is no grandparent then parent would be root and
	// hence black, but the parent is RED, so parent must have a parent.
	var gp = path.pop() //gp means grandparent
	//NOTE: path is now relative to grandparent

	//find uncle
	var uncle = parent.sibling(gp)

	// insert_case3: parent.isRed() && uncle.isRed()
	if uncle.isRed() {
		//log.Println("insertBalance: insert case #3: uncle.isRed()")
		//NOTE: isRed() method works when object is nil (it returns false).
		//Local Fact: if uncle is RED, then uncle != nil

		var nparent = parent.copy()
		nparent.setBlack()

		if less(nn.key, parent.key) {
			nparent.ln = nn
		} else {
			nparent.rn = nn
		}

		var nuncle = uncle.copy() //new uncle, cuz I am mutating it.
		nuncle.setBlack()

		var ngp = gp.copy() //new grandparent, cuz I am mutating it.
		ngp.setRed()

		if parent.isLeftChildOf(gp) {
			ngp.ln = nparent
			ngp.rn = nuncle
		} else {
			ngp.ln = nuncle
			ngp.rn = nparent
		}

		//nn = ngp //This is ok, cuz path is relative to gp

		m.insertBalance(gp, ngp, path)
		return
	}

	// From here on, we are dealing with nn, parent, and gp where path
	// is relative to gp.

	// INSERT CASE #4.1

	// create new parent and stitch nn into nparent
	var nparent = parent.copy()
	if less(nn.key, parent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}
	//log.Println("insertBalance: created new parent and stitched new node into it.")

	// create new grandparent and stitch nn into ngp
	var ngp = gp.copy()
	if parent.isLeftChildOf(gp) {
		ngp.ln = nparent
	} else {
		ngp.rn = nparent
	}
	//log.Println("insertBalance: create new grand parent and stitched parent into it.")

	// insert_case4: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	if nn.isRightChildOf(nparent) && nparent.isLeftChildOf(ngp) {
		//log.Println("insertBalance: New node is inner child of grandparent.")
		//log.Println("insertBalance: Doing prep-rotateLeft on parent.")
		nn, nparent = m.rotateLeft(nparent, ngp)
	} else if nn.isLeftChildOf(nparent) && nparent.isRightChildOf(ngp) {
		//log.Println("insertBalance: New node is inner child of grandparent.")
		//log.Println("insertBalance: Doing prep-rotateLeft on parent.")
		nn, nparent = m.rotateRight(nparent, ngp)
	}

	//nn.isRed; FYI nn was nparent
	nparent.setBlack() //FYI nparent was nn
	ngp.setRed()       //gp was black cuz parent was red

	//log.Printf("insertBalance: ngp     = %s\n", ngp)
	//log.Printf("insertBalance: nparent = %s\n", nparent)
	//log.Printf("insertBalance: nn      = %s\n", nn)
	//log.Printf("insertBalance: path = %s\n", path)

	// insert_case4: final rotate
	// This rotate makes nparent the root of this sub-tree. Hence, nparent
	// is the node we want persisted.
	if nn.isLeftChildOf(nparent) {
		//log.Println("insertBalance: insert_case #5: Doing rotateRight.")
		nparent, ngp = m.rotateRight(ngp, path.peek())
	} else {
		//log.Println("insertBalance: insert_case #5: Doing rotateLeft.")
		nparent, ngp = m.rotateLeft(ngp, path.peek())
	}
	//IMPORTANT: nn is no longer VALID from here on
	// Only, ngp and nparent are correct.

	//log.Println("insertBalance: After insert case #5:")
	//log.Printf("insertBalance: ngp     = %s\n", ngp)
	//log.Printf("insertBalance: nparent = %s\n", nparent)
	//log.Printf("insertBalance: path = %s\n", path)

	m.persistAll(gp, ngp, path)
	return
}

// Del() calls Remove() but only returns the modified *Map.
//
// I wonder if this is inlined as Delete() may have.
func (m *Map) Del(k MapKey) *Map {
	var nm /*deletedVal*/, _ /*wasDeleted*/, _ = m.Remove(k)
	return nm
}

// Delete() calls m.Del() and any call to id SHOULD be eliminated by compiler
// replaced by inlined call to m.Del() for go version >= "1.8 (1.7 on amd64)".
func (m *Map) Delete(k MapKey) *Map {
	return m.Del(k)
}

// Remove() eliminates the node pointed to by the MapKey argument (and
// rebalances) a persistent version of the given *Map.
func (m *Map) Remove(k MapKey) (*Map, interface{}, bool) {
	var on, path = m.root.findNodeWithPath(k)

	if on == nil {
		return m, nil, false
	}
	//found node associated with k
	var retVal = on.val

	var nm = m.copy()

	var oterm *node
	var nterm *node

	// if node has two children swap values with previous in-order node, then
	// delete that child, which will have at most one child of its' own.
	if on.ln != nil && on.rn != nil {
		log.Printf("Remove: on has two childred; on=%s", on)
		oterm = on
		nterm = oterm.copy() //modifiable copy of oterm

		//find victim building path
		path.push(on)
		var on = on.ln
		for on.rn != nil {
			path.push(on)
			on = on.rn
		}
		//on now points to previous node

		//trade oterm.val for previous node's val
		nterm.val = on.val
		log.Printf("Remove: setting nterm=%s", nterm)
	} else {
		log.Printf("Remove: setting oterm to m.root; oterm=%s", oterm)
		oterm = m.root
	}

	var ochild, nchild *node
	log.Printf("Remove: calling remove_one_child() on om=%s\noterm=%s\npath=%s",
		on, oterm, path)
	//remove_one_child should be remove_node_with_zero_or_one_child but that
	//is to long.
	ochild, nchild, path = nm.remove_one_child(on, oterm, path)
	//return values:
	//  ochild points to the non-persisted parent of oterm.
	//  nchild points to a modified (ie persisted) oterm.
	//  path points to the path from m.root to ochild (not including ochild).
	log.Printf("Remove: remove_one_child returned: "+
		"ochild=%s\nnchild=%s\npath=%s", ochild, nchild, path)

	if nterm != nil {
		//nterm has been allocated. This means oterm has been set to "on" and
		//"on" has be set to the node previous the origninal on. Hence we have
		//to set one of the children of the new on to ...
		// FUCK IT this explination isn't working...
		log.Println("Remove: nterm != nil")
		if oterm.ln == ochild {
			nterm.ln = nchild
		} else if oterm.rn == ochild {
			nterm.rn = nchild
		} else {
			panic("ochild != oterm child")
		}
	} else {
		//nterm has been allocated. So this means oterm was set to m.root and
		//remove_one_child() persisted till m.root (but not including)
		nterm = nchild
	}

	log.Printf("Remove: last call to persistAll:\n"+
		"oterm=%s\nnterm=%s\npath=%s", oterm, nterm, path)
	nm.persistAll(oterm, nterm, path)
	nm.numEnts--

	return nm, retVal, true
}

//func (m *Map) replaceNode(on, child *node, path *nodeStack) (
//	/* nn */ *node /* nparent */, *node, *nodeStack,
//) {
//	var oparent = path.pop()
//	if oparent != nil {
//		var nparent = oparent.copy()
//
//		if on.isLeftChildOf(oparent) {
//			nparent.ln = child
//		} else {
//			nparent.rn = child
//		}
//	} else {
//		m.root = child
//	}
//
//	return child, nparent, path
//}

//func (m *Map) replaceNode(on, child, path) (
//	/* nparent */ *node, /* path */ *nodeStack,
//) {
func (m *Map) replaceNode(on, child, oparent *node) *node {
	var nparent *node
	if oparent != nil {
		nparent = oparent.copy()
		if on.isLeftChildOf(oparent) {
			nparent.ln = child
		} else {
			nparent.rn = child
		}
	} /* else {
		m.root = child
	} */
	//return nparent, path
	return nparent
}

func (m *Map) remove_one_child0(on, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	var child *node
	if on.ln != nil {
		child = on.ln
	} else {
		child = on.rn
	}
	//note: child could be nil

	if on.isRed() {
		assert(on.ln == nil, "remove_one_child: on.isRed && on.ln != nil")
		assert(on.rn == nil, "remove_one_child: on.isReg && on.rn != nil")
		// if on.isRed is must have a parent.
		var oparent = path.pop()
		var nparent = oparent.copy()

		if on.isLeftChildOf(oparent) {
			nparent.ln = nil
		} else {
			nparent.rn = nil
		}

		if oparent == term {
			return oparent, nparent, path
		}
		log.Printf("remove_one_child0: if on.on.Red(); calling persistTill:\n"+
			"oparent=%s\nnparent=%s\nterm=%s\npath=%s\n",
			oparent, nparent, term, path)
		m.persistTill(oparent, nparent, term, path)
	}
	//on.isBlack()

	if child != nil {
		//A non-nil child MUST be red.
		//So to put it in on's position we must change it to black.
		var nn = child.copy()
		nn.setBlack()
		return on, nn, path //return directly back to term handling & persistAll
	}
	//else nn==nil && nn.isBlack

	//child==nil & child's parent is on's parent cuz of the replaceNode

	//FIXME ???

	return m.delete_case1(child, term, path)
	//my rbtree.js had a magic Nil type node, not sure if I need one.
}

//remove_one_child() deletes a node that has only on child.
//
//remove_one_child() should be remove_node_with_zero_or_one_child() but that
//is to long :).
func (m *Map) remove_one_child(on, term *node, path *nodeStack) (
	/*ochild*/ *node /*nchild*/, *node /*path*/, *nodeStack,
) {
	//find the child of the node to be deleted.
	var ochild *node
	if on.ln != nil {
		ochild = on.ln
	} else {
		ochild = on.rn
	}
	//note: ochild could be nil

	var nn *node
	var oparent = path.peek() //could be nil
	var nparent = m.replaceNode(on, ochild, oparent)
	//if oparent is nil, then nparent is nil
	log.Printf("remove_one_child: called replaceNode on:\n"+
		"on=%s\nochild=%s\noparent=%s\n==> nparent=%s",
		on, ochild, oparent, nparent)
	if on.isBlack() {
		if ochild.isRed() {
			nn = ochild.copy()
			nn.setBlack()
		} else {
			//child.isBlack and on.isBlack
			//Fact: this only happens when child == nil
			//Reason: this child's sibling is nil (hence black), if this child
			//is a non-nil black child it would violate RBT property #4.
			//This function should be called delete_node_with_zero_or_one_child.
			ochild, nn, path = m.delete_case1(on, nparent, path)
		}
	}

	if oparent != nil {
		if on.isLeftChildOf(oparent) {
			nparent.ln = nn
		} else {
			nparent.rn = nn
		}
		path.pop() //take oparent off stack
		log.Println("remove_one_child: oparent !=nil; calling persistTill on: "+
			"oparent=%s\nnparent=%s\nterm=%s\npath=%s",
			oparent, nparent, term, path)
		return m.persistTill(oparent, nparent, term, path)
	}

	log.Println("remove_one_child: returning directly...")
	return on, nn, path
}

func (m *Map) delete_case1(on, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	if path.len() > 0 {
		log.Println("delete_case1: path.len() > 0; calling delete_case2:\n"+
			"on=%s\nterm=%s\npath=%s\n", on, term, path)
		return m.delete_case2(on, term, path)
	}
	log.Println("delete_case1: calling persistTill on:"+
		"on=%s\nnn=nil\nterm=%s\npath=%s", on, term, path)
	return m.persistTill(on, nil, term, path)
}

// delete_case2() ...
//
// when sibling is red we rotate away from it. My fuzzy understanding is that
// the sibling side is longer and we are trying to shorten the target side,
// hence we need to rotate to the short side.
func (m *Map) delete_case2(on, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("delete_case2: called with: on=%s\nterm=%s\npath=%s",
		on, term, path)
	//FACT: path.len() > 0; so a parent exists
	var nn = on.copy()
	var oparent = path.peek()
	var ogp = path.peekN(1)
	var osibling = on.sibling(oparent)

	var nparent *node
	var nsibling *node

	if osibling.isRed() {
		// I need to grok the scenario (if any) where
		var ngp *node
		if ogp != nil {
			ngp = ogp.copy()
		}
		nparent = oparent.copy()
		nsibling = osibling.copy()

		nparent.setRed()
		nsibling.setBlack()

		if on.isLeftChildOf(oparent) {
			m.rotateLeft(nparent, ngp)
		} else {
			m.rotateRight(nparent, ngp)
		}
	}
	log.Printf("delete_case2: calling persist_case3: "+
		"on=%s\nnn=%s\nterm=%s\npath=%s\n", on, nn, term, path)
	return m.delete_case3(on, nn, term, path) //?????????????????????????
}

func (m *Map) delete_case3(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	panic("not implemented")
	return nil, nil, nil
	//START HERE
}

func (m *Map) RangeLimit(start, end MapKey, fn func(MapKey, interface{}) bool) {
	//get iter
	var iter = m.iterRange(start, end)

	//walk iter
	for n := iter.next(); n != nil; n = iter.next() {
		if !fn(n.key, n.val) {
			return //STOP
		}
	}

	return
}

func (m *Map) Range(fn func(MapKey, interface{}) bool) {
	m.RangeLimit(ninf, pinf, fn)
	return
}

func (m *Map) Keys() []MapKey {
	var keys = make([]MapKey, m.NumEntries())
	var i int
	var fn = func(k MapKey, v interface{}) bool {
		keys[i] = k
		i++
		return true
	}
	m.Range(fn)
	return keys
}

func (m *Map) NumEntries() uint {
	return m.numEnts
}

//DELtype visitFn func(*node, *nodeStack) bool

func (m *Map) walkPreOrder(fn func(*node, *nodeStack) bool) bool {
	var path = newNodeStack()
	return m.root.visitPreOrder(fn, path)
}

func (m *Map) walkInOrder(fn func(*node, *nodeStack) bool) bool {
	var path = newNodeStack()
	return m.root.visitInOrder(fn, path)
}

func (m *Map) TreeString() string {
	var strs = make([]string, m.numEnts)
	var i int
	var fn = func(n *node, path *nodeStack) bool {
		var pk interface{}
		var parent = path.peek()
		if parent == nil {
			pk = nil
		} else {
			pk = parent.key
		}

		strs[i] = fmt.Sprintf("parent: %#v,%p, %s%p",
			pk, parent, n.String(), n)
		i++

		return true
	}
	m.walkPreOrder(fn)

	return strings.Join(strs, "\n")
}

func (m *Map) String() string {
	var strs = make([]string, m.numEnts)
	//log.Println("String: m.numEnts =", m.numEnts)

	//var i int
	//var fn = func(n *node, path *nodeStack) bool {
	//	strs[i] = fmt.Sprintf("%#v: %#v", n.key, n.val)
	//	i++
	//	return true
	//}
	//m.walkInOrder(fn)

	var iter = m.iterAll()
	var i int
	for n := iter.next(); n != nil; n = iter.next() {
		//log.Printf("String: i=%d; n=%s\n", i, n)
		strs[i] = fmt.Sprintf("%#v: %#v", n.key, n.val)
		//log.Printf("String: strs[%d] = %q", i, strs[i])
		i++
	}

	var s = "{" + strings.Join(strs, ", ") + "}"
	//log.Println("String: return ", s)
	return s
}

package sorted_map

import (
	"errors"
	"fmt"
	"strings"
)

type Map struct {
	numEnts uint
	root    *node
}

func New() *Map {
	var m = new(Map)
	return m
}

func (m *Map) valid() error {
	var _, err = m.root.valid()
	if err != nil {
		return err
	}
	var count = m.root.count()
	if uint(count) != m.numEnts {
		return errors.New("enumerated count of subnodes != m.numEnts")
	}
	return nil
}

func (m *Map) NumEntries() uint {
	return m.numEnts
}

func (m *Map) copy() *Map {
	var nm = new(Map)
	*nm = *m
	return nm
}

func (m *Map) iter() *nodeIter {
	return m.iterRange(ninf, pinf)
}

func (m *Map) iterRange(startKey, endKey MapKey) *nodeIter {
	//log.Printf("iterRange: called with: startKey=%s; endKey=%s;",
	//	startKey, endKey)
	var dir = less(startKey, endKey)
	var cur, path = m.root.findNodeIterPath(startKey, dir)
	//log.Printf("iterRange: findNodeIterPath returned:\ncur=%s\npath=%s",
	//	cur, path)
	//log.Printf("iterRange: cur=%s\npath=%s", cur, path)
	if cur == nil {
		cur = path.pop()
		if dir { //Forw
			//is cur to far?
			if less(cur.key, startKey) {
				cur = path.pop()
			}
		} else { //Back
			//log.Printf("iterRange: init Back: less(%s, %s) => %v",
			//	endKey, cur.key, less(startKey, cur.key))
			if less(startKey, cur.key) {
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
	//panic("not implemented")
	//return nil, nil, false

	var n, path = m.root.findNodeWithPath(k)
	if n != nil {
		return m, n.val, true
	}

	var npath = path.dup()
	var nm = m.copy()
	nm.insert(k, v, npath)

	return nm, nil, false

	//var val, found = m.Load(k)
	//if !found {
	//	return m, val, true
	//}
	//
	//var nm, _ = m.Store(k, v) //don't care if it was added or replaced
	//return nm, nil, false

}

func (m *Map) Put(k MapKey, v interface{}) *Map {
	var nm, _ = m.Store(k, v)
	return nm
}

// Store() inserts a new key:val pair and returns a new Map and a boolean
// indicatiing if the key:val was added(true) or merely replaced(false).
func (m *Map) Store(k MapKey, v interface{}) (*Map, bool) {
	var on, path = m.root.findNodeDupPath(k)
	//path is duped and stiched, but not anchored to m.root

	var nm = m.copy()

	if on != nil {
		nm.replace(k, v, on, path)
		return nm, false
	}

	nm.insert(k, v, path)

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

	m.persist(on, nn, path)

	return
}

// insert() inserts a new node, create from a key-value pair,  into the tip of
// the path, then balances and persists the *Map.
//
// path MUST be non-zero length.
//
// insert() MUST be called on a new *Map.
func (m *Map) insert(k MapKey, v interface{}, path *nodeStack) {
	var on *node           // = nil
	var nn = newNode(k, v) //nn.isRed ALWAYS!

	m.insertRepair(on, nn, path)
	m.numEnts++

	return
}

//persist() takes a duped path and sets the first element of the path to the
//map's root and stitches the new node in to the last element of the path. In
//the case where the path is empty, it simply sets the map's root to the new
//node.
func (m *Map) persist(on, nn *node, path *nodeStack) {
	//log.Printf("persist: called:\non=%s\nnn=%s\npath=%s", on, nn, path)

	if path.len() == 0 {
		//log.Printf("persist: path.len() == 0; SETTING ROOT = %s", nn)
		m.root = nn
		return
	}
	//log.Printf("persist: path.len() == %d; SETTING ROOT = %s",
	//	path.len(), path.head())

	m.root = path.head()

	var parent = path.peek()
	if on == nil {
		if less(nn.key, parent.key) {
			parent.ln = nn
		} else {
			parent.rn = nn
		}
	} else {
		if on.isLeftChildOf(parent) {
			parent.ln = nn
		} else {
			parent.rn = nn
		}
	}

	return
}

// rotateLeft() takes the target node(n) and its parent(p). We are rotating on
// target node(n) left. We assume that all arguments are mutable. We return
// the original n and it's new parent.
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
	var r = n.rn //assume n.rn is already a copy.

	//log.Printf("rotateLeft:\nn = %s\np = %s\nr = %s\n", n, p, r)

	if p != nil {
		if n.isLeftChildOf(p) {
			p.ln = r
		} else {
			p.rn = r
		}
	} /* else {
		m.root = r
	} */

	n.rn = r.ln //handle anticipated orphaned node
	r.ln = n    //now orphan it

	return n, r
}

// rotateRight() takes the target node(n) and its parent(p). We are rotating on
// target node(n) right. We assume that all arguments are mutable. We return
// the original n and it's new parent.
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
	var l = n.ln //assume n.ln is already a copy.

	//log.Printf("rotateRight:\nn = %s\np = %s\nl = %s\n", n, p, l)

	if p != nil {
		if n.isLeftChildOf(p) {
			p.ln = l
		} else {
			p.rn = l
		}
	} /* else {
		m.root = l
	} */

	n.ln = l.rn //handle anticipated orphaned node
	l.rn = n    //now orphan it

	return n, l
}

// insertRepair() MUST be called on a new *Map.
func (m *Map) insertRepair(on, nn *node, path *nodeStack) {
	assert(nn != nil, "nn == nil")
	//log.Printf("insertRepair:\non=%s\nnn=%s\npath=%s\n", on, nn, path)

	var parent, gp, uncle *node

	parent = path.peek()

	gp = path.peekN(1) // peek() == peekN(0); peekN is index from top

	if gp != nil {
		uncle = parent.sibling(gp)
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
		// grandparent exists because root is never Red
		// grandparent is black because parent is Red
		m.insertCase3(on, nn, path)
	} else {
		//we know:
		//  grandparent is black because parent is Red
		//  parent.isRed
		//  uncle.isBlack
		//  nn.isRed and
		m.insertCase4(on, nn, path)
	}
}

// insertCase1() MUST be called on a new *Map.
func (m *Map) insertCase1(on, nn *node, path *nodeStack) {
	//log.Printf("insertCase1:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")

	nn.setBlack()
	m.persist(on, nn, path)
	//return
}

// insertCase2() MUST be called on a new *Map.
func (m *Map) insertCase2(on, nn *node, path *nodeStack) {
	//log.Printf("insertCase2:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	m.persist(on, nn, path)
	return
}

// insertCase3() MUST be called on a new *Map.
func (m *Map) insertCase3(on, nn *node, path *nodeStack) {
	//log.Printf("insertCase3:\non=%s\nnn=%s\npath=%s\n", on, nn, path)

	var oparent = path.pop()
	var ogp = path.pop() //gp means grandparent

	var ouncle *node
	//if parent.isLeftChildOf(gp) {
	if less(oparent.key, ogp.key) {
		ouncle = ogp.rn
	} else {
		ouncle = ogp.ln
	}

	var nparent = oparent.copy() //new parent, cuz I am mutating it.
	nparent.setBlack()

	if less(nn.key, oparent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	var nuncle = ouncle.copy() //new uncle, cuz I am mutating it.
	nuncle.setBlack()

	var ngp = ogp.copy() //new grandparent, cuz I am mutating it.
	ngp.setRed()

	//if oparent.isLeftChildOf(ogp) {
	if less(oparent.key, ogp.key) {
		ngp.ln = nparent
		ngp.rn = nuncle
	} else {
		ngp.ln = nuncle
		ngp.rn = nparent
	}

	m.insertRepair(ogp, ngp, path)
	return
}

// insertCase4() MUST be called on a new *Map.
func (m *Map) insertCase4(on, nn *node, path *nodeStack) {
	//log.Printf("insertCase4:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	var parent = path.peek()
	var gp = path.peekN(1) //ogp means grandparent

	// insertCase4.1: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	//if nn.isLeftChildOf(nparent) && oparent.isRightChildOf(ogp) {
	//if less(nn.key, oparent.key) && less(ogp.key, oparent.key) {
	if less(nn.key, parent.key) && parent.isRightChildOf(gp) {
		//log.Println("insertCase4: nn is inside left grandchild")
		parent.ln = nn

		//nn, nparent = m.rotateRight(nparent, ngp)
		parent, nn = m.rotateRight(parent, gp)
		path.pop() //take parent off path
		path.push(nn)
		//path.push(parent)

		nn = nn.rn //nn.rn == parent

		//} else if nn.isRightChildOf(nparent) && oparent.isLeftChildOf(ogp) {
		//} else if less(oparent.key, nn.key) && less(oparent.key, ogp.key) {
	} else if less(parent.key, nn.key) && parent.isLeftChildOf(gp) {
		//log.Println("insertCase4: nn is inside right grandchild")
		parent.rn = nn

		//nn, nparent = m.rotateLeft(nparent, ngp)
		parent, nn = m.rotateLeft(parent, gp)
		path.pop() //take parent off path
		path.push(nn)
		//path.push(parent)

		nn = nn.ln //nn.ln == parent
	}

	m.insertCase4pt2(on, nn, path)
	return
}

func (m *Map) insertCase4pt2(on, nn *node, path *nodeStack) {
	//log.Printf("insertCase4pt2:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	var parent = path.pop()
	var gp = path.pop()

	if less(nn.key, parent.key) {
		parent.ln = nn
	} else {
		parent.rn = nn
	}

	parent.setBlack()
	gp.setRed()

	var ggp = path.peek()

	//log.Printf("ggp Tree =\n%s", ggp.treeString())

	//if less(nn.key, parent.key) {
	if nn.isLeftChildOf(parent) {
		//I am not sure that gp.ln == parent. Unless there is some deeper
		//logic, gp.ln could be parent's sibling (aka nn's uncle). That
		//'deeper logic' could be that if the uncle existed it would have been
		//rotated away? in insertCase4.1
		if gp.ln != parent {
			if gp.ln != nil {
				gp.ln = gp.ln.copy()
			}
		}

		var t *node
		gp, t = m.rotateRight(gp, ggp)
		gp = t
	} else {
		if gp.rn != parent {
			if gp.rn != nil {
				gp.rn = gp.rn.copy()
			}
		}

		var t *node
		gp, t = m.rotateLeft(gp, ggp)
		gp = t
	}

	m.persist(gp, gp, path)
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
	var on, path = m.root.findNodeDupPath(k)
	//log.Printf("Remove: k=%s;\non=%s\npath=%s", k, on, path)

	if on == nil {
		return m, nil, false
	}
	//found node associated with k
	var retVal = on.val

	var nm = m.copy()

	if on.ln == nil || on.rn == nil {
		nm.removeNodeWithZeroOrOneChild(on, path)
		nm.numEnts--
		return nm, retVal, true
	}
	//else has two children

	// if node has two children swap values with previous in-order node, then
	// delete that child, which will have at most one child of its' own.
	//log.Printf("Remove: on has two children;\non=%s", on)
	var otcn = on          //otcn == origninal-two-child-node
	var ntcn = otcn.copy() //ntcn == new-two-child-node

	var parent = path.peek()
	if parent != nil {
		if otcn.isLeftChildOf(parent) {
			parent.ln = ntcn
		} else {
			parent.rn = ntcn
		}
	}

	//find victim, building path
	path.push(ntcn)
	parent = ntcn
	on = on.ln
	for on.rn != nil {
		var nn = on.copy()
		if on.isLeftChildOf(parent) {
			parent.ln = nn
		} else {
			parent.rn = nn
		}
		path.push(nn)
		parent = nn
		on = on.rn
	}
	//on now points to previous node

	//ntcn's position (color, ln, rn) is otcn's
	//ntcn's content (key, val) is the previous node's (aka 'on').
	ntcn.key = on.key
	ntcn.val = on.val

	//log.Printf("Remove: has set otcn & ntcn;\notcn=%s\nntcn=%s\non=%s",
	//	otcn, ntcn, on)

	//log.Printf("Remove: ntcn != nil; path.peek() Tree =\n%s",
	//	path.peek().treeString())

	nm.removeNodeWithZeroOrOneChild(on, path)
	nm.numEnts--
	return nm, retVal, true
}

//removeNodeWithZeroOrOneChild() deletes a node that has only on child.
//Basically, we reparent the child to the parent of the deleted node, then
//balance the tree. The deleteCase?() methods are the balancing methods, but
//the deletion occurs here in removeNodeWithZeroOrOneChild().
//
//Was removeOneChild() but that was confusingly wrong name, just shorter.
func (m *Map) removeNodeWithZeroOrOneChild(on *node, path *nodeStack) {
	//log.Printf("removeNodeWithZeroOrOneChild: called:\non=%s\npath=%s",
	//	on, path)

	//find the child of the node to be deleted.
	var ochild *node
	if on.ln != nil {
		ochild = on.ln
	} else {
		ochild = on.rn
	}
	//note: ochild could be nil

	var nn *node

	if on.isBlack() {
		if ochild.isRed() {
			//only way 'on' can have a non-nil child
			nn = ochild.copy()
			nn.setBlack()

			m.persist(on, nn, path)
		} else {
			//child.isBlack and on.isBlack
			//Fact: this only happens when child == nil
			//Reason: this child's sibling is nil (hence black), if this child
			//is a non-nil black child it would violate RBT property #4.

			//log.Printf("removeNodeWithZeroOrOneChild: calling deleteCase1")

			m.deleteCase1(on, nn, path) //nn == nil
		}
		return
	} /* else {
		//on.isRed
		//on has no children. cuz we know it has only zero or one child (in this
		//case zero) cuz of RBT#4 (the count of black nodes on both sides).
		//nn == nil
	} */

	//on.isRed so just delete it
	m.persist(on, nn, path) //nn == nil
	return
}

func (m *Map) deleteCase1(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase1: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	//Fact: on.isBlack()

	var oparent = path.peek()

	if oparent == nil {
		//log.Printf("deleteCase1: path.len() == 0; returning directly...\n")
		m.persist(on, nn, path)
		return
	}

	m.deleteCase2(on, nn, path)
	return
}

// deleteCase2() ...
//
// when sibling is Red we rotate away from it. My fuzzy understanding is that
// the sibling side is longer and we are trying to shorten the target side,
// hence we need to rotate to the short side.
func (m *Map) deleteCase2(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase2: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	//Fact: on.isBlack()
	//Fact: path.len() > 0

	var parent = path.pop()
	var osibling = on.sibling(parent)
	//log.Printf("deleteCase2:\nosibling=%s\n", osibling)

	var gp = path.peek() //could be nil

	var nsibling *node

	if osibling.isRed() {
		nsibling = osibling.copy()
		//if on.isLeftChildOf(oparent) {
		//if on.key.Less(oparent.key) {
		if nsibling.key.Less(parent.key) {
			//parent.rn = nn
			parent.ln = nsibling
		} else {
			parent.rn = nsibling
			//parent.ln = nn
		}

		parent.setRed()
		nsibling.setBlack()

		//if on.key.Less(oparent.key) {
		if on.isLeftChildOf(parent) {
			//log.Println("deleteCase2: on.isLeftChildOf(parent) -> rotateLeft")
			parent, nsibling = m.rotateLeft(parent, gp)
			//parent childOf nsibling childOf gp
		} else {
			//log.Println("deleteCase2: on.isRightChildOf(parent) -> rotateLeft")
			parent, nsibling = m.rotateRight(parent, gp)
			//nparent childOf nsibling childOf ngp
		}

		path.push(nsibling) //new grandparent of nn
		path.push(parent)   //new parent or nn
		//log.Printf("deleteCase2: osibling.isRed condition: nsibling Tree =\n%s",
		//	nsibling.treeString())
	} else {
		//log.Println("deleteCase2: passing thru to deleteCase3")
		path.push(parent) //put oparent back, cuz we didn't use it.
	}

	m.deleteCase3(on, nn, path)
	return
}

func (m *Map) deleteCase3(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase3: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	//Fact: path.len() > 0
	//Face: on is black

	var parent = path.peek()
	var osibling = on.sibling(parent)

	//log.Printf("deleteCase3: parent Tree =\n%s", parent.treeString())

	if parent.isBlack() &&
		osibling.isBlack() &&
		osibling.ln.isBlack() &&
		osibling.rn.isBlack() {

		//log.Printf("deleteCase3: parent.isBlack && osibling.isBlack:\n"+
		//	"on=%s\nosibling=%s\nparent=%s\n", on, osibling, parent)

		var nsibling = osibling.copy()
		if osibling.isLeftChildOf(parent) {
			parent.ln = nsibling
			parent.rn = nn
		} else {
			parent.ln = nn
			parent.rn = nsibling
		}

		nsibling.setRed()

		//log.Println("deleteCase3: going to call deleteCase1 on parent")

		path.pop() //remove parent
		m.deleteCase1(parent, parent, path)
		return
	}
	//log.Println("deleteCase3: passing thru to deleteCase4")

	m.deleteCase4(on, nn, path)
	return
}

func (m *Map) deleteCase4(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase4: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	var parent = path.peek()
	var osibling = on.sibling(parent)

	//log.Printf("deleteCase4: parent Tree =\n%s", parent.treeString())

	if parent.isRed() &&
		osibling.isBlack() &&
		osibling.ln.isBlack() &&
		osibling.rn.isBlack() {

		//log.Println("deleteCase4: is completing the deleteCase line")

		var nsibling = osibling.copy()
		//if on.key.Less(parent.key) {
		if on.isLeftChildOf(parent) {
			//log.Println("deleteCase4: parent.ln = nn")
			//parent.ln = nn
			parent.rn = nsibling
		} else {
			//log.Println("deleteCase4: parent.rn = nn")
			parent.ln = nsibling
			//parent.rn = nn
		}

		nsibling.setRed()
		parent.setBlack()

		//log.Println("deleteCase4: returing m.persist(parent, parent, path)")
		//path.pop() //remove parent from path, becasue we're returning parent
		//m.persist(parent, parent, path)
		m.persist(on, nn, path)
		return
	} /* else {
		log.Println("deleteCase4: passing thru to deleteCase5")
	} */

	m.deleteCase5(on, nn, path)
	return
}

func (m *Map) deleteCase5(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase5: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	//Fact: path.len() > 0
	var parent = path.peek()
	var osibling = on.sibling(parent)

	//This is a potential pre-rotate phase for deleteCase6
	if osibling.isBlack() {
		if on.isLeftChildOf(parent) &&
			osibling.rn.isBlack() &&
			osibling.ln.isRed() {

			//log.Println("deleteCase5: pre-rotating tree to the Right")

			var nsibling = osibling.copy()
			nsibling.ln = osibling.ln.copy()
			nsibling.setRed()
			nsibling.ln.setBlack()

			//if on.key.Less(parent.key) {
			if on.isLeftChildOf(parent) {
				parent.rn = nsibling
			} else {
				parent.ln = nsibling
			}

			//log.Printf("before rotateRight: parent Tree =\n%s",
			//	parent.treeString())

			_, _ = m.rotateRight(nsibling, parent)

			//log.Printf("after rotateRight: parent Tree =\n%s",
			//	parent.treeString())
		} else if on.isRightChildOf(parent) &&
			osibling.ln.isBlack() &&
			osibling.rn.isRed() {

			//log.Println("deleteCase5: pre-rotating tree to the Left")

			var nsibling = osibling.copy()
			nsibling.rn = osibling.rn.copy()
			nsibling.setRed()
			nsibling.rn.setBlack()

			//if on.key.Less(parent.key) {
			if on.isLeftChildOf(parent) {
				parent.rn = nsibling
			} else {
				parent.ln = nsibling
			}

			_, _ = m.rotateLeft(nsibling, parent)

			//log.Printf("parent Tree =\n%s", parent.treeString())
		} /* else {
			log.Println("deleteCase5: secondary conditions failed: " +
				"passing thru to deleteCase6")
		} */
	} /* else {
		log.Println("deleteCase5: osibling.isRed: passing thru to deleteCase6")
	} */

	m.deleteCase6(on, nn, path)
	return
}

//deleteCase6()
//We know:
//  path.len() > 0 aka oparent != nil && oparent.isRed
//  osibling != nil
//  if on.isLeftChild
//    osibling.rn != nil and isRed and ln == nil
//  else
//    osibling.ln != nil and isRed and rn == nil
func (m *Map) deleteCase6(on, nn *node, path *nodeStack) {
	//log.Printf("deleteCase6: called with:\non=%s\nnn=%s\npath=%s",
	//	on, nn, path)

	//Fact: path.len() > 0
	//Fact: sibling.isRed()

	var parent = path.pop()

	//log.Printf("deleteCase6: parent Tree:\n%s", parent.treeString())

	var osibling = on.sibling(parent)

	var nsibling = osibling.copy()

	nsibling.color = parent.color
	parent.setBlack()

	var gp = path.peek()

	//if on.key.Less(parent.key) {
	if on.isLeftChildOf(parent) {
		//log.Println("on.isLeftChildOf(parent)")
		if osibling.rn != nil {
			nsibling.rn = osibling.rn.copy()
			nsibling.rn.setBlack()
		}

		//parent.ln = nn
		parent.rn = nsibling

		//log.Printf("before rotateLeft: parent Tree:\n%s", parent.treeString())

		parent, nsibling = m.rotateLeft(parent, gp)
		//position-wise sibling replaces parent and parent replaces on

		//log.Printf("after rotateLeft: nsibling Tree:\n%s",
		//	nsibling.treeString())

		path.push(nsibling)
		path.push(parent)
	} else {
		//log.Println("!on.isLeftChildOf(parent)")
		if osibling.ln != nil {
			nsibling.ln = osibling.ln.copy()
			nsibling.ln.setBlack()
		}

		parent.ln = nsibling
		//parent.rn = nn

		//log.Printf("before rotateRight: parent Tree:\n%s",
		//	parent.treeString())

		parent, nsibling = m.rotateRight(parent, gp)
		//position-wise sibling replaces parent and parent replaces on

		//log.Printf("after rotateRight: nsibling Tree:\n%s",
		//	nsibling.treeString())

		path.push(nsibling)
		path.push(parent)
	}

	//path.pop() //take nsibling off path
	//path.pop() //take parent off path
	//m.persist(parent, nsibling, path)
	m.persist(on, nn, path)
	return
}

//RangeLimit() executes the given function starting with the start key (if
//it exists), or the first key after the start key. It then stops at the end key
//(if it exists), or the last key before the end key. The traversal will stop
//immediately if the function returns false.
//
//If the start key is greater than the end key, then the traversal will be in
//reverse order.
//
//If you want to indicate a "key greater than any key" or a "key less than any
//other key", you can use the infinitely positive or negetive key, by calling
//sorted_map.InfKey(sign int). A call to sorted_map.InfKey(1) returns a key
//greater than any other key. A call to sorted_map.InfKey(-1) returns a key less
//than any other key.
func (m *Map) RangeLimit(start, end MapKey, fn func(MapKey, interface{}) bool) {
	var iter = m.iterRange(start, end)

	//walk iter
	for n := iter.Next(); n != nil; n = iter.Next() {
		if !fn(n.key, n.val) {
			return //STOP
		}
	}

	return
}

//Range() executes the given function on every key, value pair in order. If the
//function returns false the traversal of key, value pairs will stop.
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

//DELtype visitFn func(*node, *nodeStack) bool

func (m *Map) walkPreOrder(fn func(*node, *nodeStack) bool) bool {
	if m.root != nil {
		var path = newNodeStack(0)
		return m.root.visitPreOrder(fn, path)
	}
	return true
}

func (m *Map) walkInOrder(fn func(*node, *nodeStack) bool) bool {
	if m.root != nil {
		var path = newNodeStack(0)
		return m.root.visitInOrder(fn, path)
	}
	return true
}

//dup() is for testing only. It is a recusive copy.
func (m *Map) dup() *Map {
	var nm = &Map{
		numEnts: m.numEnts,
		root:    m.root.dup(),
	}
	nm.numEnts = m.numEnts
	nm.root = m.root.dup()
	return nm
}

//equiv() is for testing only. It is a equal-by-value method.
func (m *Map) equiv(m0 *Map) bool {
	return m.numEnts == m0.numEnts && m.root.equiv(m0.root)
}

func (m *Map) treeString() string {
	//var strs = make([]string, m.numEnts)
	//var i int
	//var fn = func(n *node, path *nodeStack) bool {
	//	var pk interface{}
	//	var parent = path.peek()
	//	if parent == nil {
	//		pk = nil
	//	} else {
	//		pk = parent.key
	//	}
	//
	//	var indent = strings.Repeat("  ", path.len())
	//	strs[i] = fmt.Sprintf("%sparent: %#v,%p, %s",
	//		indent, pk, parent, n.String())
	//	i++
	//
	//	return true
	//}
	//m.walkPreOrder(fn)
	//
	//return strings.Join(strs, "\n")
	return m.root.treeString()
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

	var iter = m.iter()
	var i int
	for n := iter.Next(); n != nil; n = iter.Next() {
		//log.Printf("String: i=%d; n=%s\n", i, n)
		strs[i] = fmt.Sprintf("%#v: %#v", n.key, n.val)
		//log.Printf("String: strs[%d] = %q", i, strs[i])
		i++
	}

	var s = "{" + strings.Join(strs, ", ") + "}"
	//log.Println("String: return ", s)
	return s
}

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
	numEnts uint
	root    *node
}

func New() *Map {
	var m = new(Map)
	return m
}

//MakeNode() is public for testing only.
func MakeMap(num uint, r *node) *Map {
	return &Map{num, r}
}

//Root() is public for testing only.
func (m *Map) Root() *node {
	return m.root
}

func (m *Map) NumEntries() uint {
	return m.numEnts
}

func (m *Map) copy() *Map {
	var nm = new(Map)
	*nm = *m
	return nm
}

func (m *Map) Iter() *nodeIter {
	return m.IterRange(ninf, pinf)
}

func (m *Map) IterRange(startKey, endKey MapKey) *nodeIter {
	var cur, path = m.root.findNodeWithPath(startKey)
	//log.Printf("iterRange: cur=%s\npath=%s", cur, path)
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
	panic("not implemented")
	return nil, nil, false
	//MAYBE:var n, path = m.root.findNodeWithPath(k)
	//MAYBE:
	//MAYBE:var nm *Map
	//MAYBE:if n == nil {
	//MAYBE:	nm = m.store(k, v, n, path)
	//MAYBE:	return nm, nil, false
	//MAYBE:}
	//MAYBE:
	//MAYBE:return m, n.val, true
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
	var on *node           // = nil
	var nn = newNode(k, v) //nn.isRed ALWAYS!

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
		log.Printf("OLD TREE(on):\n%s", on.ToString(-1))
		log.Printf("OLD TREE(m.root):\n%s", m.root.ToString(-1))
		log.Printf("NEW TREE:\n%s", nn.ToString(-1))
		assert(nn.IsBlack(), "new root 'nn' is not black")
		if !nn.IsBlack() {
			nn.setBlack()
		}
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

	//log.Printf("persistTill: recursing:\noparent=%s\nnparent=%s\npath=%s\n",
	//	oparent, nparent, path)
	return m.persistTill(oparent, nparent, term, path)
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

	log.Printf("rotateLeft:\nn = %s\np = %s\nr = %s\n", n, p, r)

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

	log.Printf("rotateRight:\nn = %s\np = %s\nl = %s\n", n, p, l)

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
	} else if parent.IsBlack() {
		// we know:
		// parent exists and is Black
		m.insertCase2(on, nn, path)
	} else if uncle.IsRed() {
		// we know:
		// parent.IsRed becuase of the previous condition
		// grandparent exists because root is never Red
		// grandparent is Black because parent is Red
		m.insertCase3(on, nn, path)
	} else {
		//we know:
		//  grandparent is Black because parent is Red
		//  parent.IsRed
		//  uncle.IsBlack
		//  nn.IsRed and
		m.insertCase4(on, nn, path)
	}
}

// insertCase1() MUST be called on a new *Map.
func (m *Map) insertCase1(on, nn *node, path *nodeStack) {
	log.Printf("insertCase1:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")
	assert(m.root == on, "path.peek()==nil BUT m.root != on")

	nn.setBlack()
	m.persistAll(on, nn, path)
	return
}

// insertCase2() MUST be called on a new *Map.
func (m *Map) insertCase2(on, nn *node, path *nodeStack) {
	log.Printf("insertCase2:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	m.persistAll(on, nn, path)
	return
}

// insertCase3() MUST be called on a new *Map.
func (m *Map) insertCase3(on, nn *node, path *nodeStack) {
	log.Printf("insertCase3:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
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
	log.Printf("insertCase4:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	var oparent = path.peek()
	var ogp = path.peekN(1) //ogp means original grandparent

	// insertCase4.1: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	//if nn.isRightChildOf(nparent) && oparent.isLeftChildOf(ogp) {
	if less(oparent.key, nn.key) && oparent.isLeftChildOf(ogp) {
		var nparent = oparent.copy()
		//if less(nn.key, oparent.key) {
		//	nparent.ln = nn
		//} else {
		nparent.rn = nn
		//}

		var ngp = ogp.copy()
		//if oparent.isLeftChildOf(ogp) {
		ngp.ln = nparent
		//} else {
		//	ngp.rn = nparent
		//}

		//nn, nparent = m.rotateLeft(nparent, ngp)
		nparent, nn = m.rotateLeft(nparent, ngp)

		path.pop() //take oparent off path
		path.pop() //take ogp off path
		path.push(ngp)
		path.push(nn)

		nn = nn.ln //nn.ln == nparent; see orig (commented out) rotateLeft call
		//} else if nn.isLeftChildOf(nparent) && oparent.isRightChildOf(ogp) {
	} else if less(nn.key, oparent.key) && oparent.isRightChildOf(ogp) {
		var nparent = oparent.copy()
		//if less(nn.key, oparent.key) {
		nparent.ln = nn
		//} else {
		//	nparent.rn = nn
		//}

		var ngp = ogp.copy()
		//if oparent.isLeftChildOf(ogp) {
		//	ngp.ln = nparent
		//} else {
		ngp.rn = nparent
		//}

		//nn, nparent = m.rotateRight(nparent, ngp)
		nparent, nn = m.rotateRight(nparent, ngp)

		path.pop() //take oparent off path
		path.pop() //take ogp off path
		path.push(ngp)
		path.push(nn)

		nn = nn.rn //nn.rn == nparent; see orig (commented out) rotateRight call
	}

	m.insertCase4pt2(on, nn, path)
	return
}

func (m *Map) insertCase4pt2(on, nn *node, path *nodeStack) {
	log.Printf("insertCase4pt2:\non=%s\nnn=%s\npath=%s\n", on, nn, path)
	var oparent = path.pop()
	var ogp = path.pop()

	var nparent = oparent.copy()
	if less(nn.key, nparent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	var ngp = ogp.copy()
	if oparent.isLeftChildOf(ogp) {
		ngp.ln = nparent
	} else {
		ngp.rn = nparent
	}

	if less(nn.key, nparent.key) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	nparent.setBlack()
	ngp.setRed()

	//if on.isLeftChildOf(oparent)
	if less(nn.key, oparent.key) {
		var oggp = path.peek() //old great grand parent
		var nggp *node         //new great grandparent
		if oggp != nil {
			nggp = oggp.copy()
			path.pop()      //take oggp off path
			path.push(nggp) //replace oggp with nggp on path
		}

		//I am not sure that ngp.ln == nparent. Unless there is some deeper
		//logic, ngp.ln could be nparent's sibling (aka nn's uncle). That
		//'deeper logic' could be that if the uncle existed it would have been
		//rotated away? in insertCase4.1
		if ngp.ln != nil {
			ngp.ln = ngp.ln.copy()
		}

		var t *node
		ngp, t = m.rotateRight(ngp, nggp)
		ngp = t
	} else {
		var oggp = path.peek()
		var nggp *node
		if oggp != nil {
			nggp = oggp.copy()
			path.pop()      //take oggp off path
			path.push(nggp) //replace oggp with nggp on path
		}

		if ngp.rn != nil {
			ngp.rn = ngp.rn.copy()
		}

		var t *node
		ngp, t = m.rotateLeft(ngp, nggp)
		ngp = t
	}

	m.persistAll(ogp, ngp, path)
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
	log.Printf("Remove: k=%s;\non=%s\npaht=%s", k, on, path)

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
		oterm = m.root
		log.Printf("Remove: setting oterm to m.root; oterm=%s", oterm)
	}

	var ochild, nchild *node
	log.Printf("Remove: calling removeNodeWithZeroOrOneChild() with:\n"+
		"on=%s\noterm=%s\npath=%s", on, oterm, path)
	ochild, nchild, path = nm.removeNodeWithZeroOrOneChild(on, oterm, path)
	//return values:
	//  ochild points to the non-persisted child of oterm.
	//  nchild points to a modified (ie persisted) oterm.
	//  path points to the path from m.root to ochild (not including ochild).
	log.Printf("Remove: removeNodeWithZeroOrOneChild returned:\n"+
		"ochild=%s\nnchild=%s\npath=%s", ochild, nchild, path)

	//if on != oterm {
	if nterm != nil {
		//nterm has been allocated. This means oterm has been set to "on" and
		//"on" has be set to the node previous the origninal on. Hence we have
		//to set one of the children of the new on to ...
		// FUCK IT this explination isn't working...
		log.Println("Remove: on != oterm")
		log.Printf("on=%s\n", on)
		log.Printf("oterm=%s\n", oterm)
		log.Printf("ochild=%s\n", ochild)
		//if oterm.ln == ochild {
		//	nterm.ln = nchild
		//} else if oterm.rn == ochild {
		//	nterm.rn = nchild
		//} else {
		//	panic("ochild != oterm child")
		//}

		nm.persistAll(oterm, nterm, path)
	} else {
		nm.persistAll(on, nchild, path)
	}

	nm.numEnts--

	return nm, retVal, true
}

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

//removeNodeWithZeroOrOneChild() deletes a node that has only on child.
//Basically, we reparent the child to the parent of the deleted node, then
//balance the tree. The deleteCase?() methods are the balancing methods, but
//the deletion occurs here in removeNodeWithZeroOrOneChild().
//
//Was removeOneChild() but that was confusingly wrong name, just shorter.
func (m *Map) removeNodeWithZeroOrOneChild(on, term *node, path *nodeStack) (
	/*ochild*/ *node,
	/*nchild*/ *node,
	/*path*/ *nodeStack,
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

	if on.IsBlack() {
		if ochild.IsRed() {
			//only way on can have a non-nil child
			nn = ochild.copy()
			nn.setBlack()
		} else {
			//child.IsBlack and on.IsBlack
			//Fact: this only happens when child == nil
			//Reason: this child's sibling is nil (hence Black), if this child
			//is a non-nil Black child it would violate RBT property #4.
			//This function should be called delete_node_with_zero_or_one_child.

			//This 'if' stmt is cuz I don't want the terminator to be nil.
			//That would polute the logic of the other deleteCase methods.
			if oparent == nil {
				log.Printf("removeNodeWithZeroOrOneChild:"+
					" calling deleteCase1:\non=%s\nnn=%s\nterm=%s\npath=%s\n",
					on, nn, term, path)
				on, nn, path = m.deleteCase1(on, nn, term, path)
			} else {
				log.Printf("removeNodeWithZeroOrOneChild:"+
					" calling deleteCase1:\non=%s\nnn=%s\nterm=%s\npath=%s\n",
					on, nn, m.root, path)
				on, nn, path = m.deleteCase1(on, nn, oparent, path)
			}
			log.Printf("removeNodeWithZeroOrOneChild:"+
				" returned from deleteCase1:\non=%s\nnn=%s\npath=%s\n",
				on, nn, path)
		}
	} /* else {
		//on.IsRed
		//on has no children. cuz we know it has only zero or one child (in this
		//case zero) cuz of RBT#4 (the count of Black nodes on both sides).
		//nn == nil
	} */

	//recheck oparent = path.peek() cuz path may have changed.
	oparent = path.peek()

	//if oparent == nil {
	if on == term { //more generalized??
		//we will let the last persistAll call in remove set m.root
		log.Println("removeNodeWithZeroOrOneChild:" +
			" on==term; returning directly...")
		return on, nn, path
	}

	var nparent = oparent.copy()
	if on.isLeftChildOf(oparent) {
		nparent.ln = nn
	} else {
		nparent.rn = nn
	}

	path.pop() //take oparent off stack

	if oparent == term {
		//no need to persistTill term
		log.Printf("removeNodeWithZeroOrOneChild:" +
			" oparent == term; returning directly...")
		return oparent, nparent, path
	}

	log.Printf("removeNodeWithZeroOrOneChild:"+
		" oparent!=nil; calling persistTill"+
		" on: oparent=%s\nnparent=%s\nterm=%s\npath=%s",
		oparent, nparent, term, path)
	return m.persistTill(oparent, nparent, term, path)
}

func (m *Map) deleteCase1(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase1: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	//Fact: on.IsBlack()
	//Fact: on != term; actually term = parent(on)

	var oparent = path.peek()

	if oparent == nil {
		log.Printf("deleteCase1: path.len() == 0; returning directly...\n")
		return on, nn, path
	}

	return m.deleteCase2(on, nn, term, path)
}

// deleteCase2() ...
//
// when sibling is Red we rotate away from it. My fuzzy understanding is that
// the sibling side is longer and we are trying to shorten the target side,
// hence we need to rotate to the short side.
func (m *Map) deleteCase2(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase2: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	//Fact: on.IsBlack()
	//Fact: on != term; actually term = parent(on)
	//Fact: path.len() > 0

	var oparent = path.pop()
	var osibling = on.sibling(oparent)
	log.Printf("deleteCase2:\nosibling=%s\n", osibling)

	var ogp = path.peek() //could be nil

	var nparent *node
	var nsibling *node

	if osibling.IsRed() {
		var ngp *node
		if ogp != nil {
			ngp = ogp.copy()
			path.pop()
			path.push(ngp) //replace ogp with ngp
		}

		nparent = oparent.copy()
		nsibling = osibling.copy()
		if on.isLeftChildOf(oparent) {
			//nparent.ln = on //unnecessary
			nparent.rn = nsibling
		} else {
			nparent.ln = nsibling
			//nparent.rn = on //unnecessary
		}

		nparent.setRed()
		nsibling.setBlack()

		if on.isLeftChildOf(oparent) {
			nparent, nsibling = m.rotateLeft(nparent, ngp)
			//nparent childOf nsibling childOf ngp
		} else {
			nparent, nsibling = m.rotateRight(nparent, ngp)
			//nparent childOf nsibling childOf ngp
		}

		path.push(nsibling) //new grandparent of nn
		path.push(nparent)  //new parent or nn
	} else {
		path.push(oparent) //put oparent back, cuz we didn't use it.
	}

	return m.deleteCase3(on, nn, term, path)
}

func (m *Map) deleteCase3(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase3: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	//Fact: path.len() > 0
	//Face: on is Black

	var oparent = path.peek()
	var osibling = on.sibling(oparent)

	if oparent.IsBlack() &&
		osibling.IsBlack() &&
		osibling.ln.IsBlack() &&
		osibling.rn.IsBlack() {

		log.Println("deleteCase4: going to call deleteCase1 on oparent")

		log.Printf("deleteCase3: oparent.isBlack && osibling.isBlack:\n"+
			"on=%s\nosibling=%s\noparent=%s\n", on, osibling, oparent)

		var nsibling = osibling.copy()
		var nparent = oparent.copy()
		if osibling.isLeftChildOf(oparent) {
			nparent.ln = nsibling
			nparent.rn = nn
		} else {
			nparent.ln = nn
			nparent.rn = nsibling
		}

		nsibling.setRed()

		path.pop()
		return m.deleteCase1(oparent, nparent, term, path)
	}

	return m.deleteCase4(on, nn, term, path)
}

func (m *Map) deleteCase4(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase4: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	var oparent = path.peek()
	var osibling = on.sibling(oparent)

	if oparent.IsRed() &&
		osibling.IsBlack() &&
		osibling.ln.IsBlack() &&
		osibling.rn.IsBlack() {

		log.Println("deleteCase4: is completing the deleteCase line")

		var nsibling = osibling.copy()
		var nparent = oparent.copy()
		if on.isLeftChildOf(oparent) {
			nparent.ln = nn
			nparent.rn = nsibling
		} else {
			nparent.ln = nsibling
			nparent.rn = nn
		}

		nsibling.setRed()
		nparent.setBlack()

		path.pop() //remove parent from path, becasue we're returning parent

		if oparent == term {
			return oparent, nparent, path
		} else {
			return m.persistTill(oparent, nparent, term, path)
		}
	}

	return m.deleteCase5(on, nn, term, path)
}

func (m *Map) deleteCase5(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase5: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	//Fact: path.len() > 0
	var oparent = path.peek()
	var osibling = on.sibling(oparent)

	//This is a potential pre-rotate phase for deleteCase6
	if osibling.IsBlack() {
		if on.isLeftChildOf(oparent) &&
			osibling.rn.IsBlack() &&
			osibling.ln.IsRed() {

			log.Println("deleteCase5: pre-rotating tree to the Right")

			var nsibling = osibling.copy()
			nsibling.ln = osibling.ln.copy()
			nsibling.setRed()
			nsibling.ln.setBlack()

			var nparent = oparent.copy()
			if on.isLeftChildOf(oparent) {
				nparent.rn = nsibling
			} else {
				nparent.ln = nsibling
			}

			_, _ = m.rotateRight(nsibling, nparent)

			path.pop()         //pop off oparent
			path.push(nparent) //replace oparent with nparent
		} else if on.isRightChildOf(oparent) &&
			osibling.ln.IsBlack() &&
			osibling.rn.IsRed() {

			log.Println("deleteCase5: pre-rotating tree to the Left")

			var nsibling = osibling.copy()
			nsibling.rn = osibling.ln.copy()
			nsibling.setRed()
			nsibling.rn.setBlack()

			var nparent = oparent.copy()
			if on.isLeftChildOf(oparent) {
				nparent.rn = nsibling
			} else {
				nparent.ln = nsibling
			}

			_, _ = m.rotateLeft(nsibling, nparent)

			path.pop()         //pop off oparent
			path.push(nparent) //replace oparent with nparent
		}
	}

	return m.deleteCase6(on, nn, term, path)
}

//deleteCase6()
//We know:
//  path.len() > 0 aka oparent != nil && oparent.isRed
//  osibling != nil
//  if on.isLeftChild
//    osibling.rn != nil and isRed and ln == nil
//  else
//    osibling.ln != nil and isRed and rn == nil
func (m *Map) deleteCase6(on, nn, term *node, path *nodeStack) (
	*node, *node, *nodeStack,
) {
	log.Printf("deleteCase6: called with:\non=%s\nnn=%s\nterm=%s\npath=%s",
		on, nn, term, path)

	//Fact: path.len() > 0
	//Fact: sibling.IsRed()

	var oparent = path.pop()
	var osibling = on.sibling(oparent)

	var nsibling = osibling.copy()
	var nparent = oparent.copy()
	if on.isLeftChildOf(oparent) {
		nparent.rn = nsibling
	} else {
		nparent.ln = nsibling
	}

	nsibling.color = oparent.color
	nparent.setBlack()

	if on.isLeftChildOf(oparent) {
		nsibling.rn = osibling.rn.copy()
		nsibling.rn.setBlack()

		nparent.ln = nn

		var ogp = path.pop()
		var ngp *node
		if ogp != nil {
			ngp = ogp.copy()
			path.push(ngp) //replace ogp with ngp
		}

		nparent, nsibling = m.rotateLeft(nparent, ngp)
		//nn, nparent = m.rotateLeft(nparent, ngp)
		//position-wise sibling replaces parent and parent replaces on

		path.push(nsibling)
		//path.push(nparent)
	} else {
		nsibling.ln = osibling.ln.copy()
		nsibling.ln.setBlack()

		//nparent.rn = nn

		var ogp = path.pop()
		var ngp *node
		if ogp != nil {
			ngp = ogp.copy()
			path.push(ngp) //replace ogp with ngp
		}

		nparent, nsibling = m.rotateLeft(nparent, ngp)
		//nn, nparent = m.rotateLeft(nparent, ngp)
		//position-wise sibling replaces parent and parent replaces on

		path.push(nsibling)
		//path.push(nparent)
	}

	//return on, nn, path
	//return on, nparent, path

	path.pop()
	return oparent, nsibling, path
}

func (m *Map) RangeLimit(start, end MapKey, fn func(MapKey, interface{}) bool) {
	//get iter
	var iter = m.IterRange(start, end)

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

//DELtype visitFn func(*node, *nodeStack) bool

func (m *Map) walkPreOrder(fn func(*node, *nodeStack) bool) bool {
	if m.root != nil {
		var path = newNodeStack()
		return m.root.visitPreOrder(fn, path)
	}
	return true
}

func (m *Map) walkInOrder(fn func(*node, *nodeStack) bool) bool {
	if m.root != nil {
		var path = newNodeStack()
		return m.root.visitInOrder(fn, path)
	}
	return true
}

//Dup() is for testing only. It is a recusive copy().
//
func (m *Map) Dup() *Map {
	var nm = &Map{
		numEnts: m.numEnts,
		root:    m.root.dup(),
	}
	nm.numEnts = m.numEnts
	nm.root = m.root.dup()
	return nm
}

//Equiv() is for testing only. It is a equal-by-value method.
func (m *Map) Equiv(m0 *Map) bool {
	return m.numEnts == m0.numEnts && m.root.equiv(m0.root)
}

func (m *Map) TreeString() string {
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
	return m.root.ToString(-1)
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

	var iter = m.Iter()
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

// Package sorted_map implements a functional Map data structure that preserves
// the ordering of the keys. The internal data structure of the sorted_map is
// a regular Red-Black Tree (as opposed to a Left-Leaning Red-Black Tree).
//
// Functional means that each data structure is immutable and persistent.
// The Map is immutable because you never modify a Map in place, but rather
// every modification (like a Store or Remove) creates a new Map with that
// modification. This is not as inefficient as it sounds like it would be. Each
// modification only changes the smallest  branch of the data structure it needs
// to in order to effect the new mapping. Otherwise, the new data structure
// shares the majority of the previous data structure. That is the Persistent
// property.
//
// Each method call that potentially modifies the Map, returns a new Map data
// structure in addition to the other pertinent return values.
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
		return errors.New("enumerated count of subnodes != m.NumEntries()")
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
	var dir = less(startKey, endKey)
	var cur, path = m.root.findNodeIterPath(startKey, dir)
	if cur == nil {
		cur = path.pop()
	}
	return newNodeIter(dir, cur, endKey, path)
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
	_ = assertOn && assert(cmp(on.key, k) == 0, "on.key != nn.key")

	var nn = on.copy()
	nn.val = v

	m.persist(on, nn, path)
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
}

//persist() takes a duped path and sets the first element of the path to the
//map's root and stitches the new node in to the last element of the path. In
//the case where the path is empty, it simply sets the map's root to the new
//node.
func (m *Map) persist(on, nn *node, path *nodeStack) {
	if path.len() == 0 {
		m.root = nn
		return
	}

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
	//_ = assertOn && assert(n == p.rn, "new node is not right child of new parent")
	//_ = assertOn && assert(p == nil || n == p.ln,
	//	"node is not the left child of parent")
	var r = n.rn //assume n.rn is already a copy.

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
	//_ = assertOn && assert(l == n.ln, "new node is not left child of new parent")
	//_ = assertOn && assert(p == nil || n == p.rn,
	//	"node is not the right child of parent")
	var l = n.ln //assume n.ln is already a copy.

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
	_ = assertOn && assert(nn != nil, "nn == nil")

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
	_ = assertOn && assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")

	nn.setBlack()
	m.persist(on, nn, path)
}

// insertCase2() MUST be called on a new *Map.
func (m *Map) insertCase2(on, nn *node, path *nodeStack) {
	m.persist(on, nn, path)
}

// insertCase3() MUST be called on a new *Map.
func (m *Map) insertCase3(on, nn *node, path *nodeStack) {
	var oparent = path.pop()
	var ogp = path.pop() //gp means grandparent

	var ouncle *node
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
}

// insertCase4() MUST be called on a new *Map.
func (m *Map) insertCase4(on, nn *node, path *nodeStack) {
	var parent = path.peek()
	var gp = path.peekN(1) //ogp means grandparent

	// insertCase4.1: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	//if nn.isLeftChildOf(nparent) && oparent.isRightChildOf(ogp) {
	//if less(nn.key, oparent.key) && less(ogp.key, oparent.key) {
	if less(nn.key, parent.key) && parent.isRightChildOf(gp) {
		parent.ln = nn

		parent, nn = m.rotateRight(parent, gp)
		path.pop() //take parent off path
		path.push(nn)

		nn = nn.rn //nn.rn == parent
	} else if less(parent.key, nn.key) && parent.isLeftChildOf(gp) {
		parent.rn = nn

		parent, nn = m.rotateLeft(parent, gp)
		path.pop() //take parent off path
		path.push(nn)

		nn = nn.ln //nn.ln == parent
	}

	m.insertCase4pt2(on, nn, path)
}

func (m *Map) insertCase4pt2(on, nn *node, path *nodeStack) {
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
}

func (m *Map) deleteCase1(on, nn *node, path *nodeStack) {
	//Fact: on.isBlack()
	var oparent = path.peek()

	if oparent == nil {
		m.persist(on, nn, path)
		return
	}

	m.deleteCase2(on, nn, path)
}

// deleteCase2() ...
//
// when sibling is Red we rotate away from it. My fuzzy understanding is that
// the sibling side is longer and we are trying to shorten the target side,
// hence we need to rotate to the short side.
func (m *Map) deleteCase2(on, nn *node, path *nodeStack) {
	//Fact: on.isBlack()
	//Fact: path.len() > 0

	var parent = path.pop()
	var osibling = on.sibling(parent)

	var gp = path.peek() //could be nil

	var nsibling *node

	if osibling.isRed() {
		nsibling = osibling.copy()
		if nsibling.key.Less(parent.key) {
			//parent.rn = nn
			parent.ln = nsibling
		} else {
			parent.rn = nsibling
			//parent.ln = nn
		}

		parent.setRed()
		nsibling.setBlack()

		if on.isLeftChildOf(parent) {
			parent, nsibling = m.rotateLeft(parent, gp)
			//parent childOf nsibling childOf gp
		} else {
			parent, nsibling = m.rotateRight(parent, gp)
			//nparent childOf nsibling childOf ngp
		}

		path.push(nsibling) //new grandparent of nn
		path.push(parent)   //new parent or nn
	} else {
		path.push(parent) //put oparent back, cuz we didn't use it.
	}

	m.deleteCase3(on, nn, path)
}

func (m *Map) deleteCase3(on, nn *node, path *nodeStack) {
	//Fact: path.len() > 0
	//Face: on is black

	var parent = path.peek()
	var osibling = on.sibling(parent)

	if parent.isBlack() &&
		osibling.isBlack() &&
		osibling.ln.isBlack() &&
		osibling.rn.isBlack() {

		var nsibling = osibling.copy()
		if osibling.isLeftChildOf(parent) {
			parent.ln = nsibling
			parent.rn = nn
		} else {
			parent.ln = nn
			parent.rn = nsibling
		}

		nsibling.setRed()

		path.pop() //remove parent
		m.deleteCase1(parent, parent, path)
		return
	}

	m.deleteCase4(on, nn, path)
}

func (m *Map) deleteCase4(on, nn *node, path *nodeStack) {
	var parent = path.peek()
	var osibling = on.sibling(parent)

	if parent.isRed() &&
		osibling.isBlack() &&
		osibling.ln.isBlack() &&
		osibling.rn.isBlack() {

		var nsibling = osibling.copy()
		//if on.key.Less(parent.key) {
		if on.isLeftChildOf(parent) {
			parent.rn = nsibling
		} else {
			parent.ln = nsibling
			//parent.rn = nn
		}

		nsibling.setRed()
		parent.setBlack()

		m.persist(on, nn, path)
		return
	}

	m.deleteCase5(on, nn, path)
}

func (m *Map) deleteCase5(on, nn *node, path *nodeStack) {
	//Fact: path.len() > 0
	var parent = path.peek()
	var osibling = on.sibling(parent)

	//This is a potential pre-rotate phase for deleteCase6
	if osibling.isBlack() {
		if on.isLeftChildOf(parent) &&
			osibling.rn.isBlack() &&
			osibling.ln.isRed() {

			var nsibling = osibling.copy()
			nsibling.ln = osibling.ln.copy()
			nsibling.setRed()
			nsibling.ln.setBlack()

			parent.rn = nsibling

			_, _ = m.rotateRight(nsibling, parent)
		} else if on.isRightChildOf(parent) &&
			osibling.ln.isBlack() &&
			osibling.rn.isRed() {

			var nsibling = osibling.copy()
			nsibling.rn = osibling.rn.copy()
			nsibling.setRed()
			nsibling.rn.setBlack()

			parent.ln = nsibling

			_, _ = m.rotateLeft(nsibling, parent)
		}
	}

	m.deleteCase6(on, nn, path)
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
	//Fact: path.len() > 0
	//Fact: sibling.isRed()

	var parent = path.pop()

	var osibling = on.sibling(parent)

	var nsibling = osibling.copy()

	nsibling.color = parent.color
	parent.setBlack()

	var gp = path.peek()

	if on.isLeftChildOf(parent) {
		if osibling.rn != nil {
			nsibling.rn = osibling.rn.copy()
			nsibling.rn.setBlack()
		}

		//parent.ln = nn
		parent.rn = nsibling

		parent, nsibling = m.rotateLeft(parent, gp)
		//position-wise sibling replaces parent and parent replaces on

		path.push(nsibling)
		path.push(parent)
	} else {
		if osibling.ln != nil {
			nsibling.ln = osibling.ln.copy()
			nsibling.ln.setBlack()
		}

		parent.ln = nsibling
		//parent.rn = nn

		parent, nsibling = m.rotateRight(parent, gp)
		//position-wise sibling replaces parent and parent replaces on

		path.push(nsibling)
		path.push(parent)
	}

	m.persist(on, nn, path)
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
}

//Range() executes the given function on every key, value pair in order. If the
//function returns false the traversal of key, value pairs will stop.
func (m *Map) Range(fn func(MapKey, interface{}) bool) {
	m.RangeLimit(ninf, pinf, fn)
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

//func (m *Map) walkPreOrder(fn func(*node, *nodeStack) bool) bool {
//	if m.root != nil {
//		var path = newNodeStack(0)
//		return m.root.visitPreOrder(fn, path)
//	}
//	return true
//}

//func (m *Map) walkInOrder(fn func(*node, *nodeStack) bool) bool {
//	if m.root != nil {
//		var path = newNodeStack(0)
//		return m.root.visitInOrder(fn, path)
//	}
//	return true
//}

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
		strs[i] = fmt.Sprintf("%#v: %#v", n.key, n.val)
		i++
	}

	var s = "{" + strings.Join(strs, ", ") + "}"
	return s
}

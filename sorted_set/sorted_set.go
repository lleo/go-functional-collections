package sorted_set

import (
	"errors"
	"fmt"
	"strings"
)

type Set struct {
	numEnts uint
	root    *node
}

func New() *Set {
	var s = new(Set)
	return s
}

func (s *Set) valid() error {
	var _, err = s.root.valid()
	if err != nil {
		return err
	}
	var count = s.root.count()
	if uint(count) != s.numEnts {
		return errors.New("enumerated count of subnodes != s.NumEntries()")
	}
	return nil
}

func (s *Set) NumEntries() uint {
	return s.numEnts
}

func (s *Set) copy() *Set {
	var nm = new(Set)
	*nm = *s
	return nm
}

func (s *Set) iter() *nodeIter {
	return s.iterRange(ninf, pinf)
}

func (s *Set) iterRange(startKey, endKey SetKey) *nodeIter {
	var dir = less(startKey, endKey)
	var cur, path = s.root.findNodeIterPath(startKey, dir)
	if cur == nil {
		cur = path.pop()
	}
	return newNodeIter(dir, cur, endKey, path)
}

func (s *Set) IsSet(k SetKey) bool {
	var n = s.root.findNode(k)
	return n != nil
}

func (s *Set) Set(k SetKey) *Set {
	var nm, _ = s.Add(k)
	return nm
}

// Add() inserts a new key and returns a new Set and a boolean
// indicatiing if the key was added(true) or merely replaced(false).
func (s *Set) Add(k SetKey) (*Set, bool) {
	var on, path = s.root.findNodeDupPath(k)
	//path is duped and stiched, but not anchored to s.root

	if on != nil {
		return s, false
	}

	var nm = s.copy()
	nm.insert(k, path)

	return nm, true
}

//insert() inserts a new node into the tip of the path, then balances and
//persists the *Set.
//
//path MUST be non-zero length.
//
//insert() MUST be called on a new *Set.
func (s *Set) insert(k SetKey, path *nodeStack) {
	var on *node        // = nil
	var nn = newNode(k) //nn.isRed ALWAYS!

	s.insertRepair(on, nn, path)
	s.numEnts++
}

//persist() takes a duped path and sets the first element of the path to the
//set's root and stitches the new node in to the last element of the path. In
//the case where the path is empty, it simply sets the set's root to the new
//node.
func (s *Set) persist(on, nn *node, path *nodeStack) {
	if path.len() == 0 {
		s.root = nn
		return
	}

	s.root = path.head()

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
func (s *Set) rotateLeft(n, p *node) (*node, *node) {
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
		s.root = r
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
func (s *Set) rotateRight(n, p *node) (*node, *node) {
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
		s.root = l
	} */

	n.ln = l.rn //handle anticipated orphaned node
	l.rn = n    //now orphan it

	return n, l
}

// insertRepair() MUST be called on a new *Set.
func (s *Set) insertRepair(on, nn *node, path *nodeStack) {
	_ = assertOn && assert(nn != nil, "nn == nil")

	var parent, gp, uncle *node

	parent = path.peek()

	gp = path.peekN(1) // peek() == peekN(0); peekN is index from top

	if gp != nil {
		uncle = parent.sibling(gp)
	}

	if parent == nil {
		s.insertCase1(on, nn, path)
	} else if parent.isBlack() {
		// we know:
		// parent exists and is black
		s.insertCase2(on, nn, path)
	} else if uncle.isRed() {
		// we know:
		// parent.isRed becuase of the previous condition
		// grandparent exists because root is never Red
		// grandparent is black because parent is Red
		s.insertCase3(on, nn, path)
	} else {
		//we know:
		//  grandparent is black because parent is Red
		//  parent.isRed
		//  uncle.isBlack
		//  nn.isRed and
		s.insertCase4(on, nn, path)
	}
}

// insertCase1() MUST be called on a new *Set.
func (s *Set) insertCase1(on, nn *node, path *nodeStack) {
	_ = assertOn && assert(path.len() == 0, "path.peek()==nil BUT path.len() != 0")

	nn.setBlack()
	s.persist(on, nn, path)
}

// insertCase2() MUST be called on a new *Set.
func (s *Set) insertCase2(on, nn *node, path *nodeStack) {
	s.persist(on, nn, path)
}

// insertCase3() MUST be called on a new *Set.
func (s *Set) insertCase3(on, nn *node, path *nodeStack) {
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

	s.insertRepair(ogp, ngp, path)
}

// insertCase4() MUST be called on a new *Set.
func (s *Set) insertCase4(on, nn *node, path *nodeStack) {
	var parent = path.peek()
	var gp = path.peekN(1) //ogp means grandparent

	// insertCase4.1: conditional prep-rotate
	// We pre-rotate when nn is the inner child of the grandparent.
	//if nn.isLeftChildOf(nparent) && oparent.isRightChildOf(ogp) {
	//if less(nn.key, oparent.key) && less(ogp.key, oparent.key) {
	if less(nn.key, parent.key) && parent.isRightChildOf(gp) {
		parent.ln = nn

		parent, nn = s.rotateRight(parent, gp)
		path.pop() //take parent off path
		path.push(nn)

		nn = nn.rn //nn.rn == parent
	} else if less(parent.key, nn.key) && parent.isLeftChildOf(gp) {
		parent.rn = nn

		parent, nn = s.rotateLeft(parent, gp)
		path.pop() //take parent off path
		path.push(nn)

		nn = nn.ln //nn.ln == parent
	}

	s.insertCase4pt2(on, nn, path)
}

func (s *Set) insertCase4pt2(on, nn *node, path *nodeStack) {
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
		gp, t = s.rotateRight(gp, ggp)
		gp = t
	} else {
		if gp.rn != parent {
			if gp.rn != nil {
				gp.rn = gp.rn.copy()
			}
		}

		var t *node
		gp, t = s.rotateLeft(gp, ggp)
		gp = t
	}

	s.persist(gp, gp, path)
}

// Del() calls Remove() but only returns the modified *Set.
//
// I wonder if this is inlined as Delete() may have.
func (s *Set) Unset(k SetKey) *Set {
	var nm, _ = s.Remove(k)
	return nm
}

// Remove() eliminates the node pointed to by the SetKey argument (and
// rebalances) a persistent version of the given *Set.
func (s *Set) Remove(k SetKey) (*Set, bool) {
	var on, path = s.root.findNodeDupPath(k)

	if on == nil {
		return s, false
	}
	//found node associated with k

	var nm = s.copy()

	if on.ln == nil || on.rn == nil {
		nm.removeNodeWithZeroOrOneChild(on, path)
		nm.numEnts--
		return nm, true
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
	//ntcn's content (key) is the previous node's (aka 'on').
	ntcn.key = on.key

	nm.removeNodeWithZeroOrOneChild(on, path)
	nm.numEnts--
	return nm, true
}

//removeNodeWithZeroOrOneChild() deletes a node that has only on child.
//Basically, we reparent the child to the parent of the deleted node, then
//balance the tree. The deleteCase?() methods are the balancing methods, but
//the deletion occurs here in removeNodeWithZeroOrOneChild().
//
//Was removeOneChild() but that was confusingly wrong name, just shorter.
func (s *Set) removeNodeWithZeroOrOneChild(on *node, path *nodeStack) {
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

			s.persist(on, nn, path)
		} else {
			//child.isBlack and on.isBlack
			//Fact: this only happens when child == nil
			//Reason: this child's sibling is nil (hence black), if this child
			//is a non-nil black child it would violate RBT property #4.

			s.deleteCase1(on, nn, path) //nn == nil
		}
		return
	} /* else {
		//on.isRed
		//on has no children. cuz we know it has only zero or one child (in this
		//case zero) cuz of RBT#4 (the count of black nodes on both sides).
		//nn == nil
	} */

	//on.isRed so just delete it
	s.persist(on, nn, path) //nn == nil
}

func (s *Set) deleteCase1(on, nn *node, path *nodeStack) {
	//Fact: on.isBlack()
	var oparent = path.peek()

	if oparent == nil {
		s.persist(on, nn, path)
		return
	}

	s.deleteCase2(on, nn, path)
}

// deleteCase2() ...
//
// when sibling is Red we rotate away from it. My fuzzy understanding is that
// the sibling side is longer and we are trying to shorten the target side,
// hence we need to rotate to the short side.
func (s *Set) deleteCase2(on, nn *node, path *nodeStack) {
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
			parent, nsibling = s.rotateLeft(parent, gp)
			//parent childOf nsibling childOf gp
		} else {
			parent, nsibling = s.rotateRight(parent, gp)
			//nparent childOf nsibling childOf ngp
		}

		path.push(nsibling) //new grandparent of nn
		path.push(parent)   //new parent or nn
	} else {
		path.push(parent) //put oparent back, cuz we didn't use it.
	}

	s.deleteCase3(on, nn, path)
}

func (s *Set) deleteCase3(on, nn *node, path *nodeStack) {
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
		s.deleteCase1(parent, parent, path)
		return
	}

	s.deleteCase4(on, nn, path)
}

func (s *Set) deleteCase4(on, nn *node, path *nodeStack) {
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

		s.persist(on, nn, path)
		return
	}

	s.deleteCase5(on, nn, path)
}

func (s *Set) deleteCase5(on, nn *node, path *nodeStack) {
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

			_, _ = s.rotateRight(nsibling, parent)
		} else if on.isRightChildOf(parent) &&
			osibling.ln.isBlack() &&
			osibling.rn.isRed() {

			var nsibling = osibling.copy()
			nsibling.rn = osibling.rn.copy()
			nsibling.setRed()
			nsibling.rn.setBlack()

			parent.ln = nsibling

			_, _ = s.rotateLeft(nsibling, parent)
		}
	}

	s.deleteCase6(on, nn, path)
}

//deleteCase6()
//We know:
//  path.len() > 0 aka oparent != nil && oparent.isRed
//  osibling != nil
//  if on.isLeftChild
//    osibling.rn != nil and isRed and ln == nil
//  else
//    osibling.ln != nil and isRed and rn == nil
func (s *Set) deleteCase6(on, nn *node, path *nodeStack) {
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

		parent, nsibling = s.rotateLeft(parent, gp)
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

		parent, nsibling = s.rotateRight(parent, gp)
		//position-wise sibling replaces parent and parent replaces on

		path.push(nsibling)
		path.push(parent)
	}

	s.persist(on, nn, path)
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
//sorted_set.InfKey(sign int). A call to sorted_set.InfKey(1) returns a key
//greater than any other key. A call to sorted_set.InfKey(-1) returns a key less
//than any other key.
func (s *Set) RangeLimit(start, end SetKey, fn func(SetKey) bool) {
	var iter = s.iterRange(start, end)

	//walk iter
	for n := iter.Next(); n != nil; n = iter.Next() {
		if !fn(n.key) {
			return //STOP
		}
	}
}

//Range() executes the given function on every key, value pair in order. If the
//function returns false the traversal of key, value pairs will stop.
func (s *Set) Range(fn func(SetKey) bool) {
	s.RangeLimit(ninf, pinf, fn)
}

func (s *Set) Keys() []SetKey {
	var keys = make([]SetKey, s.NumEntries())
	var i int
	var fn = func(k SetKey) bool {
		keys[i] = k
		i++
		return true
	}
	s.Range(fn)
	return keys
}

//func (s *Set) walkPreOrder(fn func(*node, *nodeStack) bool) bool {
//	if s.root != nil {
//		var path = newNodeStack(0)
//		return s.root.visitPreOrder(fn, path)
//	}
//	return true
//}

//func (s *Set) walkInOrder(fn func(*node, *nodeStack) bool) bool {
//	if s.root != nil {
//		var path = newNodeStack(0)
//		return s.root.visitInOrder(fn, path)
//	}
//	return true
//}

//dup() is for testing only. It is a recusive copy.
func (s *Set) dup() *Set {
	var nm = &Set{
		numEnts: s.numEnts,
		root:    s.root.dup(),
	}
	nm.numEnts = s.numEnts
	nm.root = s.root.dup()
	return nm
}

//equiv() is for testing only. It is a equal-by-value method.
func (s *Set) equiv(m0 *Set) bool {
	return s.numEnts == m0.numEnts && s.root.equiv(m0.root)
}

func (s *Set) treeString() string {
	//var strs = make([]string, s.numEnts)
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
	//s.walkPreOrder(fn)
	//
	//return strings.Join(strs, "\n")
	return s.root.treeString()
}

func (s *Set) String() string {
	var strs = make([]string, s.numEnts)

	//var i int
	//var fn = func(n *node, path *nodeStack) bool {
	//	strs[i] = fmt.Sprintf("%#v", n.key)
	//	i++
	//	return true
	//}
	//s.walkInOrder(fn)

	var iter = s.iter()
	var i int
	for n := iter.Next(); n != nil; n = iter.Next() {
		strs[i] = fmt.Sprintf("%#v", n.key)
		i++
	}

	return "{" + strings.Join(strs, ", ") + "}"
}

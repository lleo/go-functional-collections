package sorted_map

import (
	"errors"
	"fmt"
)

type colorType bool

func (c colorType) String() string {
	if !c {
		return "RED"
	}
	return "BLACK"
}

const (
	//Black is public for testing only.
	Black = colorType(true)
	//Red is public for testing only.
	Red = colorType(false)
)

// color() returns the color of a node, the reason for its existence is to
// treat nil *node values as Black.
func color(n *node) colorType {
	if n == nil {
		return Black
	}
	return n.color
}

type node struct {
	key   MapKey
	val   interface{}
	color colorType //default node is RED aka false
	ln    *node
	rn    *node
}

func newNode(k MapKey, v interface{}) *node {
	var n = new(node)
	n.key = k
	n.val = v
	//n.color = Red   //default
	//n.ln = nil      //default
	//n.rn = nil      //default
	return n
}

//func NewNode(k MapKey, v interface{}) *node {
//	var n = new(node)
//	n.key = k
//	n.val = v
//	//n.color = Red   //default
//	//n.ln = nil      //default
//	//n.rn = nil      //default
//	return n
//}

//MakeIntNode() exists for testing only.
func MakeIntNode(i int, c colorType, ln, rn *node) *node {
	return &node{IntKey(i), i, c, ln, rn}
}

//Key() exists for testing only.
func (n *node) Key() MapKey {
	return n.key
}

//Val() exists for testing only.
func (n *node) Val() interface{} {
	return n.val
}

//Color() exists for testing only.
func (n *node) Color() colorType {
	return color(n)
}

//Ln() exists for testing only.
func (n *node) Ln() *node {
	return n.ln
}

//Rn() exists for testing only.
func (n *node) Rn() *node {
	return n.rn
}

func (n *node) copy() *node {
	var nn = new(node)
	*nn = *n
	return nn
}

//count() sums up the number of sub-nodes plus this node.
func (n *node) count() int {
	if n == nil {
		return 0
	}
	return n.ln.count() + n.rn.count() + 1
}

//Ln() is public for testing only. It returns a boolean indicating if the
//sub-tree represented by this node is valid w/r RED-BLACK-TREE-PROPERTIES.md,
//and the count of the black nodes in the left sub-tree path plus one for this
//node (if it is black). If the node is not valid, then the black node count
//will be -1.
func (n *node) Valid() (int, error) {
	//RBT#2
	if n == nil {
		return 1, nil
	}
	var lcount, lerr = n.ln.Valid()
	var rcount, rerr = n.rn.Valid()

	if lerr != nil {
		return -1, lerr
	}
	if rerr != nil {
		return -1, rerr
	}

	//RBT#4
	if lcount != rcount || lcount < 0 {
		var errStr = fmt.Sprintf("left count,%d != right count,%d",
			lcount, rcount)
		return -1, errors.New(errStr)
	}

	//RBT#3
	if n.IsRed() {
		if n.ln.IsRed() || n.rn.IsRed() {
			return -1, errors.New("Red-Red violation.")
		}
	} else {
		lcount++
	}

	return lcount, nil
}

//dup() is for testing only. It is a recursive copy().
func (n *node) dup() *node {
	if n == nil {
		return nil
	}
	var nn = &node{
		key:   n.key,
		val:   n.val,
		color: n.color,
		ln:    n.ln.dup(),
		rn:    n.rn.dup(),
	}
	return nn
}

//equiv() is for testing only. It is a equal-by-value method.
func (n *node) equiv(n0 *node) bool {
	if n == nil {
		return n0 == nil
	} else if n0 == nil {
		return false
	}
	//n != nil && n0 != nil

	if cmp(n.key, n0.key) != 0 {
		return false
	}
	if n.val != n0.val {
		return false
	}
	if n.color != n0.color {
		return false
	}

	if !n.ln.equiv(n0.ln) {
		return false
	}
	if !n.rn.equiv(n0.rn) {
		return false
	}
	return true
}

// IsRed() is public for testing only.
func (n *node) IsRed() bool {
	return bool(!color(n)) //given that Red is encoded with a false value
}

// IsBlack() is public for testing only.
func (n *node) IsBlack() bool {
	return bool(color(n)) //given that Black is encoded as true
}

func (n *node) setBlack() *node {
	n.color = Black
	return n
}

func (n *node) setRed() *node {
	n.color = Red
	return n
}

func (n *node) isLeftChildOf(parent *node) bool {
	if parent.ln == n {
		return true
	}
	return false
}

func (n *node) isRightChildOf(parent *node) bool {
	if parent.rn == n {
		return true
	}
	return false
}

func (n *node) sibling(parent *node) *node {
	if parent.ln == n {
		return parent.rn
	}
	return parent.ln
}

func (n *node) findNode(k MapKey) *node {
	if n == nil {
		return nil
	}

	var cur = n
	for cur != nil {
		switch {
		case less(k, cur.key):
			cur = cur.ln
		case less(cur.key, k):
			cur = cur.rn
		default:
			return cur
		}
	}
	return nil
}

func (n *node) findNodeWithPath(k MapKey) (*node, *nodeStack) {
	var path = newNodeStack(0)
	var cur = n
	for cur != nil {
		var ocur = cur
		switch {
		case less(k, cur.key):
			cur = cur.ln
		case less(cur.key, k):
			cur = cur.rn
		default:
			return cur, path
		}
		path.push(ocur)
	}
	return nil, path
}

func (n *node) findNodeDupPath(k MapKey) (*node, *nodeStack) {
	var path = newNodeStack(0)
	if n == nil {
		return nil, path
	}

	var cur = n
	var ocur = cur
	switch {
	case less(k, cur.key):
		cur = cur.ln
	case less(cur.key, k):
		cur = cur.rn
	default: //cur.key == k
		return cur, path
	}
	var parent = ocur.copy()
	path.push(parent)

	for cur != nil {
		ocur = cur
		switch {
		case less(k, cur.key):
			cur = cur.ln
		case less(cur.key, k):
			cur = cur.rn
		default: //cur.key == k
			return cur, path
		}

		var ncur = ocur.copy()
		if ocur.isLeftChildOf(parent) {
			parent.ln = ncur
		} else {
			parent.rn = ncur
		}
		parent = ncur
		path.push(parent)
	}

	return nil, path
}

//	var parent *node
//	for cur != nil {
//		//var ocur = cur
//		ocur = cur
//		switch {
//		case less(k, cur.key):
//			cur = cur.ln
//		case less(cur.key, k):
//			cur = cur.rn
//		default: //cur.key == k
//			return cur, path
//		}
//
//		//var parent = path.peek()
//		var ncur = ocur.copy()
//		if parent != nil {
//			//if less(ncur.key, parent.key) {
//			if ocur.isLeftChildOf(parent) {
//				parent.ln = ncur
//			} else {
//				parent.rn = ncur
//			}
//		}
//
//		path.push(ncur)
//		parent = ncur
//	}
//	return nil, path
//}

func (n *node) findNodeIterPath(k MapKey, dir bool) (*node, *nodeStack) {
	var path = newNodeStack(0)
	var cur = n
	for cur != nil {
		switch {
		case less(k, cur.key):
			if dir { //if dir==forw(true) then path.push(cur)
				path.push(cur)
			}
			cur = cur.ln
		case less(cur.key, k):
			if !dir { //if dir=back(false) then path.push(cur)
				path.push(cur)
			}
			cur = cur.rn
		default: //cur.key == k
			return cur, path
		}
	}
	return nil, path
}

// visitPreOrder() calls the visit function on the current node, then
// conditionally calls visitPreOrder on its children. The condition is
// if the node exists.
//
// should never be called when n == nil.
func (n *node) visitPreOrder(
	fn func(*node, *nodeStack) bool,
	path *nodeStack,
) bool {
	assert(n != nil, "visitPreOrder() called when n == nil")

	if !fn(n, path) {
		return false
	}

	if n.ln != nil {
		path.push(n)
		if !n.ln.visitPreOrder(fn, path) {
			return false
		}
		path.pop()
	}

	if n.rn != nil {
		path.push(n)
		if !n.rn.visitPreOrder(fn, path) {
			return false
		}
		path.pop()
	}

	return true
}

// visitInOrder() conditionally calls visitInOrder() on its left child, then
// calls the visit function on the current node, and finnaly conditionally
// calls visitInOrder() on its' right child. The condition is if the node
// exists.
//
// should never be called when n == nil.
func (n *node) visitInOrder(
	fn func(*node, *nodeStack) bool,
	path *nodeStack,
) bool {
	assert(n != nil, "visitInOrder() called when n == nil")

	if n.ln != nil {
		path.push(n)
		if !n.ln.visitInOrder(fn, path) {
			return false
		}
		path.pop()
	}

	if !fn(n, path) {
		return false
	}

	if n.rn != nil {
		path.push(n)
		if !n.rn.visitInOrder(fn, path) {
			return false
		}
		path.pop()
	}

	return true
}

const toStrFmt0 = "%p,node{key:%s, val:%#v, color:%s\n"
const toStrFmt1 = "%s  ln: %s,\n"
const toStrFmt2 = "%s  rn: %s,\n"
const toStrFmt3 = "%s}\n"

// TreeString() prints the node and all children to a depth of d. For example,
// if d==0 it only prints the given node; if d==1 then it prints the node and
// it's left and write children. Finnaly, if d < 0 it will print the entire
// tree starting at the given node.
//func (n *node) TreeString(d int) string {
func (n *node) TreeString() string {
	return n.treeString(-1, "")
}

func (n *node) treeString(d int, indent string) string {
	if n == nil {
		return "<nil>"
	}
	//if d < 0 {
	//	d = 0
	//}
	if d == 0 {
		return n.String()
	}
	if n.ln == nil && n.rn == nil {
		return n.String()
	}
	return fmt.Sprintf(toStrFmt0, n, n.key, n.val, n.color) +
		fmt.Sprintf(toStrFmt1, indent, n.ln.treeString(d-1, indent+"  ")) +
		fmt.Sprintf(toStrFmt2, indent, n.rn.treeString(d-1, indent+"  ")) +
		indent + "}"
}

func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%p,node{key:%s, val:%#v, color:%s ln:%p, rn:%p}",
		n, n.key, n.val, n.color, n.ln, n.rn)
}

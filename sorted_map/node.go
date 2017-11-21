package sorted_map

import (
	"fmt"
)

type Color bool

func (c Color) String() string {
	if !c {
		return "RED"
	}
	return "BLACK"
}

const (
	//Black is public for testing only.
	Black = Color(true)
	//Red is public for testing only.
	Red = Color(false)
)

// color() returns the color of a Node, the reason for its existence is to
// treat nil *Node values as Black.
func color(n *Node) Color {
	if n == nil {
		return Black
	}
	return n.color
}

//Node struct is public for testing only.
type Node struct {
	key   MapKey
	val   interface{}
	color Color //default Node is RED aka false
	ln    *Node
	rn    *Node
}

func newNode(k MapKey, v interface{}) *Node {
	var n = new(Node)
	n.key = k
	n.val = v
	//n.color = Red   //default
	//n.ln = nil      //default
	//n.rn = nil      //default
	return n
}

//MakeNode() is here for testing only.
func MakeNode(k MapKey, v interface{}, c Color, ln, rn *Node) *Node {
	return &Node{k, v, c, ln, rn}
}

//Key() is here for testing only.
func (n *Node) Key() MapKey {
	return n.key
}

//Val() is here for testing only.
func (n *Node) Val() interface{} {
	return n.val
}

//Color() is here for testing only.
func (n *Node) Color() Color {
	return n.color
}

//Ln() is here for testing only.
func (n *Node) Ln() *Node {
	return n.ln
}

//Rn() is here for testing only.
func (n *Node) Rn() *Node {
	return n.rn
}

func (n *Node) copy() *Node {
	var nn = new(Node)
	*nn = *n
	return nn
}

// IsRed() returns true if the color is Red.
func (n *Node) IsRed() bool {
	return bool(!color(n)) //given that Red is encoded with a false value
}

func (n *Node) IsBlack() bool {
	return bool(color(n)) //given that Black is encoded as true
}

func (n *Node) setBlack() *Node {
	n.color = Black
	return n
}

func (n *Node) setRed() *Node {
	n.color = Red
	return n
}

func (n *Node) isLeftChildOf(parent *Node) bool {
	if parent.ln == n {
		return true
	}
	return false
}

func (n *Node) isRightChildOf(parent *Node) bool {
	if parent.rn == n {
		return true
	}
	return false
}

func (n *Node) sibling(parent *Node) *Node {
	if parent.ln == n {
		return parent.rn
	}
	return parent.ln
}

func (n *Node) findNode(k MapKey) *Node {
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

func (n *Node) findNodeWithPath(k MapKey) (*Node, *nodeStack) {
	var path = newNodeStack()
	var cur = n
	//log.Printf("findNodeWithPath: cur=%s\n", cur)
	for cur != nil {
		switch {
		case less(k, cur.key):
			//log.Printf("findNodeWithPath: k,%s < cur.key,%s\n", k, cur.key)
			path.push(cur)
			cur = cur.ln
		case less(cur.key, k):
			//log.Printf("findNodeWithPath: cur.key,%s < k,%s\n", cur.key, k)
			path.push(cur)
			cur = cur.rn
		default:
			//log.Printf("findNodeWithPath: returning cur=%s\n%s", cur, path)
			return cur, path
		}
	}
	//log.Printf("findNodeWithPath: returning cur=nil\n%s", path)
	return nil, path
}

// visitPreOrder() calls the visit function on the current Node, then
// conditionally calls visitPreOrder on its children. The condition is
// if the Node exists.
//
// should never be called when n == nil.
func (n *Node) visitPreOrder(
	fn func(*Node, *nodeStack) bool,
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
// calls the visit function on the current Node, and finnaly conditionally
// calls visitInOrder() on its' right child. The condition is if the Node
// exists.
//
// should never be called when n == nil.
func (n *Node) visitInOrder(
	fn func(*Node, *nodeStack) bool,
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

func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}

	//var lnStr, rnStr string
	//if n.ln == nil {
	//	lnStr = "nil"
	//} else {
	//	lnStr = "!nil"
	//}
	//if n.rn == nil {
	//	rnStr = "nil"
	//} else {
	//	rnStr = "!nil"
	//}

	return fmt.Sprintf("Node{key:%s, val:%#v, color:%s ln:%p, rn:%p}",
		n.key, n.val, n.color, n.ln, n.rn)
}

package sorted_map

import (
	"fmt"
)

type node struct {
	key   MapKey
	val   interface{}
	color colorType //default node is RED aka false
	ln    *node
	rn    *node
}

type colorType bool

func (c colorType) String() string {
	if !c {
		return "RED"
	}
	return "BLACK"
}

const (
	black = colorType(true)
	red   = colorType(false)
)

// color() returns the color of a node, the reason for its existence is to
// treat nil *node values as black.
func color(n *node) colorType {
	if n == nil {
		return black
	}
	return n.color
}

func newNode(k MapKey, v interface{}) *node {
	var n = new(node)
	n.key = k
	n.val = v
	//n.color = red   //default
	//n.ln = nil      //default
	//n.rn = nil      //default
	return n
}

func (n *node) copy() *node {
	var nn = new(node)
	*nn = *n
	return nn
}

// isRed() returns true if the color is red.
func (n *node) isRed() bool {
	return bool(!color(n)) //given that red is encoded with a false value
}

func (n *node) isBlack() bool {
	return bool(color(n)) //given that black is encoded as true
}

func (n *node) setBlack() *node {
	n.color = black
	return n
}

func (n *node) setRed() *node {
	n.color = red
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

func (n *node) visitPreOrder(
	fn func(*node, *nodeStack) bool,
	path *nodeStack,
) bool {
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

func (n *node) visitInOrder(
	fn func(*node, *nodeStack) bool,
	path *nodeStack,
) bool {
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

func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}

	var lnStr, rnStr string
	if n.ln == nil {
		lnStr = "nil"
	} else {
		lnStr = "!nil"
	}
	if n.rn == nil {
		rnStr = "nil"
	} else {
		rnStr = "!nil"
	}

	return fmt.Sprintf("node{key:%s, val:%#v, color:%s ln:%s, rn:%s}",
		n.key, n.val, n.color, lnStr, rnStr)
}

package sorted_map

import "strings"

type nodeStack []*node

func newNodeStack() *nodeStack {
	var ns nodeStack = make([]*node, 0)
	return &ns
}

func (ns *nodeStack) push(n *node) *nodeStack {
	(*ns) = append(*ns, n)
	return ns
}

func (ns *nodeStack) pop() *node {
	if len(*ns) == 0 {
		return nil
	}
	var n = (*ns)[len(*ns)-1]
	*ns = (*ns)[:len(*ns)-1]
	return n
}

func (ns *nodeStack) peek() *node {
	if len(*ns) == 0 {
		return nil
	}
	return (*ns)[len(*ns)-1]
}

// peekN() is index from top. ie peekN(0) = (*ns)[len(*ns)-1]
func (ns *nodeStack) peekN(n int) *node {
	if len(*ns) < 1+n {
		return nil
	}
	return (*ns)[len(*ns)-1-n]
}

func (ns *nodeStack) len() int {
	return len(*ns)
}

func (ns *nodeStack) String() string {
	var strs = make([]string, ns.len())
	for i, n := range *ns {
		strs[i] = n.String()
	}
	return "[ " + strings.Join(strs, ",\n  ") + " ]"
}

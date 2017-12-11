package sorted_map

type nodeIter struct {
	dir    bool // true == lower-to-higher; false == higher-to-lower
	endKey MapKey
	cur    *node
	path   *nodeStack
}

func newNodeIter(start *node, endKey MapKey, path *nodeStack) *nodeIter {
	var iter = new(nodeIter)
	iter.dir = less(start.key, endKey)
	iter.endKey = endKey
	iter.cur = start
	iter.path = path
	return iter
}

func (it *nodeIter) Next() *node {
	if it.dir {
		return it.Forw()
	} else {
		return it.Back()
	}
}

func (it *nodeIter) Forw() *node {
	if it.cur == nil {
		//the iterator is over
		return nil
	}
	if it.toFar() {
		return nil
	}

	var ret = it.cur

	// set it.cur to next node
	if it.cur.rn != nil {
		//go right one then left-most
		it.cur = it.cur.rn
		for it.cur.ln != nil {
			//only push when going left
			it.path.push(it.cur)
			it.cur = it.cur.ln
		}
	} else {
		it.cur = it.path.pop()
	}

	return ret
}

func (it *nodeIter) Back() *node {
	if it.cur == nil {
		//the iterator is over
		return nil
	}
	if it.toFar() {
		return nil
	}
	var ret = it.cur

	// set it.cur to previous node
	if it.cur.ln != nil {
		//go left one then right-most
		it.cur = it.cur.ln
		for it.cur.rn != nil {
			//only push when going right
			it.path.push(it.cur)
			it.cur = it.cur.rn
		}
	} else {
		it.cur = it.path.pop()
	}

	return ret
}

func (it *nodeIter) toFar() bool {
	if it.dir {
		// lower to higher
		return less(it.endKey, it.cur.key) // cur <= end
	} else {
		// higher to lower
		return less(it.cur.key, it.endKey) // end <= cur
	}
}

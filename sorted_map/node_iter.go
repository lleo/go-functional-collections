package sorted_map

type nodeIter struct {
	dir    bool // true == lower-to-higher; false == higher-to-lower
	endKey MapKey
	cur    *node
	path   *nodeStack
}

func newNodeIter(start *node, endKey MapKey, path *nodeStack) *nodeIter {
	//log.Printf("newNodeIter: start = %s\n", start)
	//log.Printf("newNodeIter: endKey = %s\n", endKey)
	//log.Printf("newNodeIter: path = \n%s\n", path)
	var iter = new(nodeIter)
	iter.dir = less(start.key, endKey)
	//log.Printf("newNodeIter: dir = %v\n", iter.dir)
	iter.endKey = endKey
	iter.cur = start
	iter.path = path
	return iter
}

func (it *nodeIter) next() *node {
	if it.dir {
		return it.forw()
	} else {
		return it.back()
	}
}

func (it *nodeIter) forw() *node {
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
		it.path.push(it.cur)
		it.cur = it.cur.rn
		for it.cur.ln != nil {
			it.path.push(it.cur)
			it.cur = it.cur.ln
		}
	} else if it.path.len() != 0 {
		if it.cur.isLeftChildOf(it.path.peek()) {
			it.cur = it.path.pop()
		} else { // it.cur is the right child
			it.cur = it.path.pop() //temporary
			if it.path.len() != 0 {
				if it.cur.isLeftChildOf(it.path.peek()) {
					it.cur = it.path.pop()
				} else {
					it.cur = nil
				}
			} else {
				it.cur = nil
			}
		}
	} else {
		it.cur = nil
	}

	return ret
}

func (it *nodeIter) back() *node {
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
		it.path.push(it.cur)
		it.cur = it.cur.ln
		for it.cur.rn != nil {
			it.path.push(it.cur)
			it.cur = it.cur.rn
		}
	} else if it.path.len() != 0 {
		if it.cur.isRightChildOf(it.path.peek()) {
			it.cur = it.path.pop()
		} else { // it.cur is the left child
			it.path.pop() //temporary
			if it.path.len() != 0 {
				if it.cur.isRightChildOf(it.path.peek()) {
					it.cur = it.path.pop()
				} else {
					it.cur = nil
				}
			} else {
				it.cur = nil
			}
		}
	} else {
		it.cur = nil
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

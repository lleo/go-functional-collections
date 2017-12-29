package sorted_map

type Iter struct {
	dir    bool // true == lower-to-higher; false == higher-to-lower
	endKey MapKey
	cur    *node
	path   *nodeStack
}

func newNodeIter(
	dir bool,
	start *node,
	endKey MapKey,
	path *nodeStack,
) *Iter {
	var iter = new(Iter)
	iter.dir = dir
	iter.endKey = endKey
	iter.cur = start
	iter.path = path
	return iter
}

func (it *Iter) Next() (MapKey, interface{}) {
	if it.dir {
		return it.forw()
	} else {
		return it.back()
	}
}

func (it *Iter) forw() (MapKey, interface{}) {
	if it.cur == nil {
		//the iterator is over
		return nil, nil
	}
	if it.toFar() {
		return nil, nil
	}

	var retKey, retVal = it.cur.key, it.cur.val

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

	return retKey, retVal
}

func (it *Iter) back() (MapKey, interface{}) {
	if it.cur == nil {
		//the iterator is over
		return nil, nil
	}
	if it.toFar() {
		return nil, nil
	}

	var retKey, retVal = it.cur.key, it.cur.val

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

	return retKey, retVal
}

func (it *Iter) toFar() bool {
	if it.dir {
		// lower to higher
		return less(it.endKey, it.cur.key) // cur <= end
	} else {
		// higher to lower
		return less(it.cur.key, it.endKey) // end <= cur
	}
}

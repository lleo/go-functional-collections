package sortedMap

import "github.com/lleo/go-functional-collections/key"

// Iter struct mmaintins the current state for walking the *Set data structure.
type Iter struct {
	dir    bool // true == lower-to-higher; false == higher-to-lower
	endKey key.Sort
	cur    *node
	path   *nodeStack
}

func newNodeIter(
	dir bool,
	start *node,
	endKey key.Sort,
	path *nodeStack,
) *Iter {
	var iter = new(Iter)
	iter.dir = dir
	iter.endKey = endKey
	iter.cur = start
	iter.path = path
	return iter
}

// Next returns each sucessive key in the *Set. When all entries have been
// returned it will return a nil key.Sort
func (it *Iter) Next() (key.Sort, interface{}) {
	if it.dir {
		return it.forw()
	}
	return it.back()
}

func (it *Iter) forw() (key.Sort, interface{}) {
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

func (it *Iter) back() (key.Sort, interface{}) {
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
		return key.Less(it.endKey, it.cur.key) // cur <= end
	}
	// higher to lower
	return key.Less(it.cur.key, it.endKey) // end <= cur
}

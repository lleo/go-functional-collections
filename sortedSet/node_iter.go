package sorted_set

import (
	"github.com/lleo/go-functional-collections/sorted"
)

// Iter struct mmaintins the current state for walking the *Map data structure.
type Iter struct {
	dir    bool // true == lower-to-higher; false == higher-to-lower
	endKey sorted.Key
	cur    *node
	path   *nodeStack
}

func newNodeIter(
	dir bool,
	start *node,
	endKey sorted.Key,
	path *nodeStack,
) *Iter {
	var iter = new(Iter)
	iter.dir = dir
	iter.endKey = endKey
	iter.cur = start
	iter.path = path
	return iter
}

// Next returns each sucessive key/value mapping in the *Map. When all entries
// have been returned it will return a nil sorted.Key.
func (it *Iter) Next() sorted.Key {
	if it.dir {
		return it.forw()
	} else {
		return it.back()
	}
}

func (it *Iter) forw() sorted.Key {
	if it.cur == nil {
		//the iterator is over
		return nil
	}
	if it.toFar() {
		return nil
	}

	var ret = it.cur.key

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

func (it *Iter) back() sorted.Key {
	if it.cur == nil {
		//the iterator is over
		return nil
	}
	if it.toFar() {
		return nil
	}
	var ret = it.cur.key

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

func (it *Iter) toFar() bool {
	if it.dir {
		// lower to higher
		return sorted.Less(it.endKey, it.cur.key) // cur <= end
	} else {
		// higher to lower
		return sorted.Less(it.cur.key, it.endKey) // end <= cur
	}
}

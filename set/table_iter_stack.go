package set

import (
	"github.com/lleo/go-functional-collections/set/hash"
)

type tableIterStack []tableIterFunc

func newTableIterStack() tableIterStack {
	var tis tableIterStack = make([]tableIterFunc, 0, hash.DepthLimit)
	return tis
}

func (tis *tableIterStack) push(f tableIterFunc) {
	(*tis) = append(*tis, f)
}

func (tis *tableIterStack) pop() tableIterFunc {
	if len(*tis) == 0 {
		return nil
	}

	var last = len(*tis) - 1
	var f = (*tis)[last]
	*tis = (*tis)[:last]

	return f
}

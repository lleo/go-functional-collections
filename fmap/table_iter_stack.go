package fmap

import (
	"github.com/lleo/go-functional-collections/key/hash"
)

type tableIterStack []tableIterFunc

func newTableIterStack() *tableIterStack {
	var tis tableIterStack = make([]tableIterFunc, 0, hash.DepthLimit)
	return &tis
}

func (tis *tableIterStack) push(f tableIterFunc) {
	(*tis) = append(*tis, f)
}

func (tis *tableIterStack) pop() tableIterFunc {
	if tis == nil {
		panic("WTF!!! iterator is nil")
	}
	if len(*tis) == 0 {
		return nil
	}

	var last = len(*tis) - 1
	var f = (*tis)[last]
	*tis = (*tis)[:last]

	return f
}

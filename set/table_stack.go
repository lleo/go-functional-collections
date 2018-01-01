package set

import (
	"strings"

	"github.com/lleo/go-functional-collections/hash"
)

type tableStack []tableI

func newTableStack() *tableStack {
	var ts = make(tableStack, 0, hash.MaxDepth)
	return &ts
}

// ts.peek() returns the last entry without removing it from the tableStack
func (ts *tableStack) peek() tableI {
	if len(*ts) == 0 {
		return nil
	}
	return (*ts)[len(*ts)-1]
}

// Put a new tableI in the ts object.
// You should never push nil, but we are not checking to prevent this.
func (ts *tableStack) push(tab tableI) *tableStack {
	*ts = append(*ts, tab)
	return ts
}

// ts.pop() returns & remmoves the last entry inserted with ts.push(...).
func (ts *tableStack) pop() tableI {
	if len(*ts) == 0 {
		return nil
	}

	var parent = (*ts)[len(*ts)-1]
	*ts = (*ts)[:len(*ts)-1]
	return parent
}

//func (ts *tableStack) shift() tableI {
//	var t tableI
//	t, *ts = (*ts)[0], (*ts)[1:]
//	return t
//}
//
//func (ts *tableStack) unshift(t tableI) *tableStack {
//	*ts = append([]tableI{t}, (*ts)...)
//	return ts
//}

//// ts.isEmpty() returns true if there are no entries in the ts object,
//// otherwise it returns false.
//func (ts *tableStack) isEmpty() bool {
//	return len(*ts) == 0
//}

func (ts *tableStack) len() int {
	return len(*ts)
}

// Convert ts to a string representation. This is only good for debug messages.
// It is not a string format to convert back from.
func (ts *tableStack) String() string {
	var tsStrs = make([]string, len(*ts))

	for i, tab := range *ts {
		tsStrs[i] = tab.String()
	}

	var vals = strings.Join(tsStrs, ", ")

	return "[ " + vals + " ]"
}

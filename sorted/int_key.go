package sorted

import (
	"strconv"
)

// IntKey is a type wrapper of int values that implements the sorted.Key
// interface.
type IntKey int

// Less returns true if passed another IntKey that it based on an int value that
// is less than the integer value of the receiver. Iteger comparison is used.
//
// If the sorted.Key value that is passed to the Less method that is not a
// StringKey type, Less will panic.
func (ik IntKey) Less(o Key) bool {
	var oik, ok = o.(IntKey)
	if !ok {
		panic("o is not a IntKey")
	}
	return ik < oik
}

func (ik IntKey) String() string {
	return strconv.Itoa(int(ik))
}

// Package sorted provides the sorted.Key interface used by sorted_map and
// sorted_set, as well as implementations of the sorted.Key interface for basic
// types like int an string. It also provides supporting functions for the
// sorted_map and sorted_set data structures.
package sorted

type Key interface {
	Less(Key) bool
	String() string
}

// nInf is a Key for negative infinity
type nInf struct{}

func (nInf) Less(Key) bool {
	return true
}

func (nInf) String() string {
	return "nInf"
}

// pInf is a Key for positive infinity
type pInf struct{}

func (pInf) Less(Key) bool {
	return false
}

func (pInf) String() string {
	return "pInf"
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

// InfKey() is the constructor for the generic infinite sorted.Key values. If
// InfKey() is passed a non-negitive integer it will return a sorted.Key that
// is greater than any other sorted.Key. If it is passed a negative integer it
// will return a sorted.Key value that is less than any other sorted.Key value.
func InfKey(sign int) Key {
	if sign < 0 {
		return ninf
	}
	return pinf
}

// Less returns the result of x.Less(y), but first checks if either x or y is
// an infinite key value.
//
// If x is a positive infinity key , or y is a negitive inifinity key, Less()
// returns false.
//
// Then if x is a negative infinity key, or y is a positive infinity key, Less()
// returns true.
//
// Finally, the result of x.Less(y) is returned.
func Less(x, y Key) bool {
	if x == pinf || y == ninf {
		return false
	}
	if x == ninf || y == pinf {
		return true
	}
	return x.Less(y)
}

// Cmp does a standard unix-style comparison of two sorted.Key inferface
// implementing structs. Internally is uses sorted.Less(). It returns a
// positive integer if x > y, a negative integer if x < y, and zero if x == y.
// Because it uses sorted.Less() internally, it honors positive and negative
// infinity keys; sorted.Inf(1) and sorted.InfKey(-1)
func Cmp(x, y Key) int {
	if Less(x, y) {
		return -1
	} else if Less(y, x) {
		return 1
	}
	return 0
}

package key

import (
	"github.com/lleo/go-functional-collections/key/hash"
)

type Hash interface {
	Hash() hash.Val
	Equals(Hash) bool
	String() string
}

type Sort interface {
	Less(Sort) bool
	String() string
}

// nInf is a Sort for negative infinity
type nInf struct{}

func (nInf) Less(Sort) bool {
	return true
}

func (nInf) String() string {
	return "nInf"
}

// pInf is a Sort for positive infinity
type pInf struct{}

func (pInf) Less(Sort) bool {
	return false
}

func (pInf) String() string {
	return "pInf"
}

var (
	ninf = nInf{}
	pinf = pInf{}
)

// Inf() is the constructor for the generic infinite key.Sort values. If
// Inf() is passed a non-negitive integer it will return a key.Sort that
// is greater than any other key.Sort. If it is passed a negative integer it
// will return a key.Sort value that is less than any other key.Sort value.
func Inf(sign int) Sort {
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
func Less(x, y Sort) bool {
	if x == pinf || y == ninf {
		return false
	}
	if x == ninf || y == pinf {
		return true
	}
	return x.Less(y)
}

// Cmp does a standard unix-style comparison of two key.Sort inferface
// implementing structs. Internally is uses sorted.Less(). It returns a
// positive integer if x > y, a negative integer if x < y, and zero if x == y.
// Because it uses sorted.Less() internally, it honors positive and negative
// infinity keys; sorted.Inf(1) and sorted.Inf(-1)
func Cmp(x, y Sort) int {
	if Less(x, y) {
		return -1
	} else if Less(y, x) {
		return 1
	}
	return 0
}

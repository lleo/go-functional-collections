package key

import (
	"github.com/lleo/go-functional-collections/hash"
)

// Str is a type wrapper of string values that implements the hash.Key and
// sorted.Key interfaces.
type Str string

func (sk Str) Less(okey Sort) bool {
	var osk, ok = okey.(Str)
	if !ok {
		panic("okey is not a key.Str")
		//return false
	}
	return sk < osk
}

// Hash calculates the hash.Val if the key.Str receiver every time it is
// called.
func (sk Str) Hash() hash.Val {
	return hash.Calculate([]byte(sk))
}

// Equals determines if the given Key is equivalent, by value, to the receiver.
func (sk Str) Equals(okey Hash) bool {
	var osk, ok = okey.(Str)
	if !ok {
		return false
	}
	return sk == osk
}

// String returns a string representation of the receiver.
func (sk Str) String() string {
	return string(sk)
}

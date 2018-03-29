package key

import (
	"github.com/lleo/go-functional-collections/key/hash"
)

type Sort interface {
	Less(Sort) bool
	String() string
}

//func (s0 Sort) Equals(s1 Sort) bool {
//	return !(s0.Less(s1) || s1.Less(s0))
//}

type Hash interface {
	Hash() hash.Val
	Equals(Hash) bool
	String() string
}

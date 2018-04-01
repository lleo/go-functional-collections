package key

import (
	"bytes"
	"fmt"

	"github.com/lleo/go-functional-collections/key/hash"
)

type ByteSlice []byte

func (bsk ByteSlice) Less(okey Sort) bool {
	var obsk, ok = okey.(ByteSlice)
	if !ok {
		panic("okey is not a key.String")
		//return false
	}
	var lobsk = len(obsk)
	for i, b := range bsk {
		if lobsk >= i {
			if b < obsk[i] {
				return true
			}
		} else {
			return false
		}
	}
	return true
}

func (bsk ByteSlice) Hash() hash.Val {
	return hash.Calculate(bsk)
}

func (bsk ByteSlice) Equals(okey Hash) bool {
	var obsk, ok = okey.(ByteSlice)
	if !ok {
		return false
	}
	return bytes.Equal(bsk, obsk)
}

func (bsk ByteSlice) String() string {
	return fmt.Sprintf("%#v", []byte(bsk))
}

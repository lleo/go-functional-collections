package string_keyed_map

import (
	"github.com/lleo/go-functional-collections/fmap"
	"github.com/lleo/go-functional-collections/fmap/hash"
)

type StringKey string

func (sk StringKey) Hash() hash.HashVal {
	return hash.CalcHash([]byte(sk))
}

func (sk StringKey) Equals(okey fmap.MapKey) bool {
	var osk, ok = okey.(StringKey)
	if !ok {
		return false
	}
	return sk == osk
}

func (sk StringKey) String() string {
	return string(sk)
}

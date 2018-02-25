package hash

import (
	"log"
	"strconv"
)

// IntKey defines a type of 'int' that has the methods that satisfy the hash.Key
// interface.
type IntKey int

const uintSize = 4 << (^uint(0) >> 32 & 1) // 32 or 64

// Int2ByteSlice is a function that gets initialized at runtime based on the
// size of a 'int'.
//
// If the size of an 'int' is neither 4bytes nor 8bytes the function will panic
// if called.
var Int2ByteSlice func(int) []byte

func init() {
	switch uintSize {
	case 4:
		Int2ByteSlice = func(i int) []byte {
			return []byte{
				byte(i),
				byte((uint(i) >> 8)),
				byte((uint(i) >> 16)),
				byte((uint(i) >> 24)),
			}
		}
	case 8:
		Int2ByteSlice = func(i int) []byte {
			return []byte{
				byte(i),
				byte((uint(i) >> 8)),
				byte((uint(i) >> 16)),
				byte((uint(i) >> 24)),
				byte((uint(i) >> 32)),
				byte((uint(i) >> 40)),
				byte((uint(i) >> 48)),
				byte((uint(i) >> 56)),
			}
		}
	default:
		Int2ByteSlice = func(i int) []byte {
			log.Panicf("Int2ByteSlice not defined for uintSize=%d\n", uintSize)
			return nil
		}
	}
}

// Hash calculates the hash.Val of the IntKey receiver every time it is called.
func (ik IntKey) Hash() Val {
	var ib = Int2ByteSlice(int(ik))
	return CalcHash(ib)
}

// Equals determines if the given Key is equivalent, by value, to the receiver.
func (ik IntKey) Equals(okey Key) bool {
	var oik, ok = okey.(IntKey)
	if !ok {
		return false
	}
	return ik == oik
}

// String returns a string representation of the receiver.
func (ik IntKey) String() string {
	return strconv.FormatInt(int64(ik), 10)
}
